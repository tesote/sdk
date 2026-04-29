require 'json'
require 'net/http'
require 'openssl'
require 'securerandom'
require 'uri'

require_relative 'errors'
require_relative 'version'

module TesoteSdk
  RateLimitInfo = Struct.new(:limit, :remaining, :reset, keyword_init: true)

  # Minimal LRU + TTL cache. Duck-typed: any object responding to #read(key)
  # and #write(key, value, ttl:) can be passed via :cache_backend.
  class CacheBackend
    Entry = Struct.new(:value, :expires_at)

    def initialize(max_size: 256)
      @max_size = max_size
      @data = {}
      @mutex = Mutex.new
    end

    def read(key)
      @mutex.synchronize do
        entry = @data.delete(key)
        return nil if entry.nil?
        if entry.expires_at < monotonic_now
          return nil
        end

        @data[key] = entry
        entry.value
      end
    end

    def write(key, value, ttl:)
      @mutex.synchronize do
        @data.delete(key)
        @data[key] = Entry.new(value, monotonic_now + ttl)
        @data.shift while @data.size > @max_size
        value
      end
    end

    def clear
      @mutex.synchronize { @data.clear }
    end

    private

    def monotonic_now
      Process.clock_gettime(Process::CLOCK_MONOTONIC)
    end
  end

  class Transport
    DEFAULT_BASE_URL = 'https://equipo.tesote.com/api'.freeze
    DEFAULT_OPEN_TIMEOUT = 5
    DEFAULT_READ_TIMEOUT = 30
    DEFAULT_MAX_ATTEMPTS = 3
    DEFAULT_BASE_DELAY = 0.250
    DEFAULT_MAX_DELAY = 8.0
    MUTATING_METHODS = %w[POST PUT PATCH DELETE].freeze
    RETRYABLE_STATUSES = [429, 502, 503, 504].freeze

    attr_reader :api_key, :base_url, :version_segment, :user_agent,
                :open_timeout, :read_timeout, :max_attempts, :base_delay,
                :max_delay, :logger, :cache_backend

    attr_accessor :last_rate_limit, :last_request_id

    def initialize(api_key:,
                   version_segment:,
                   base_url: DEFAULT_BASE_URL,
                   user_agent: nil,
                   open_timeout: DEFAULT_OPEN_TIMEOUT,
                   read_timeout: DEFAULT_READ_TIMEOUT,
                   max_attempts: DEFAULT_MAX_ATTEMPTS,
                   base_delay: DEFAULT_BASE_DELAY,
                   max_delay: DEFAULT_MAX_DELAY,
                   logger: nil,
                   cache_backend: nil,
                   sleeper: nil,
                   randomizer: nil)
      raise ConfigError, 'api_key is required' if api_key.nil? || api_key.to_s.empty?
      raise ConfigError, 'version_segment is required' if version_segment.nil? || version_segment.to_s.empty?

      @api_key = api_key.to_s
      @version_segment = version_segment.to_s
      @base_url = base_url.to_s.sub(%r{/+\z}, '')
      @user_agent = user_agent || default_user_agent
      @open_timeout = open_timeout
      @read_timeout = read_timeout
      @max_attempts = max_attempts
      @base_delay = base_delay
      @max_delay = max_delay
      @logger = logger
      @cache_backend = cache_backend
      @sleeper = sleeper || ->(seconds) { sleep(seconds) }
      @randomizer = randomizer || ->(max) { Kernel.rand(max) }
      @last_rate_limit = nil
      @last_request_id = nil
    end

    # opts:
    #   :idempotency_key      → forwarded as Idempotency-Key (auto-gen for mutations)
    #   :cache                → false to bypass; { ttl: int } to enable TTL cache
    #   :extra_headers        → hash of additional headers
    def request(method, path, query: nil, body: nil, opts: {})
      method_upper = method.to_s.upcase
      uri = build_uri(path, query)
      execute_request(method_upper, uri, body, opts)
    end

    # why: GET /status and GET /whoami live at the API root, not under
    # /v1 or /v2 — bypass the version_segment but reuse all cross-cutting.
    def request_unversioned(method, path, query: nil, body: nil, opts: {})
      method_upper = method.to_s.upcase
      uri = build_unversioned_uri(path, query)
      execute_request(method_upper, uri, body, opts)
    end

    # Returns a RawResponse with body string + headers — used for file-download
    # endpoints (CSV/JSON export) where the SDK should not parse the body.
    RawResponse = Struct.new(:status, :body, :content_type, :content_disposition, :request_id, keyword_init: true)

    def request_raw(method, path, query: nil, body: nil, opts: {})
      method_upper = method.to_s.upcase
      uri = build_uri(path, query)
      request_summary = build_request_summary(method_upper, uri, body)
      response, body_str, attempts = perform_with_retries(method_upper, uri, body, opts, request_summary)

      record_rate_limit(response)
      @last_request_id = response['x-request-id'] || response['X-Request-Id']
      status = response.code.to_i

      if status >= 200 && status < 300
        return RawResponse.new(
          status: status,
          body: body_str,
          content_type: response['content-type'] || response['Content-Type'],
          content_disposition: response['content-disposition'] || response['Content-Disposition'],
          request_id: @last_request_id
        )
      end

      raise ApiError.from_response(response, body_str, request_summary, attempts: attempts)
    end

    private

    def execute_request(method_upper, uri, body, opts)
      cache_key = cache_key_for(method_upper, uri, opts)
      cached = cache_lookup(method_upper, cache_key, opts)
      return cached unless cached.nil?

      request_summary = build_request_summary(method_upper, uri, body)
      response, body_str, attempts = perform_with_retries(method_upper, uri, body, opts, request_summary)

      record_rate_limit(response)
      @last_request_id = response['x-request-id'] || response['X-Request-Id']

      status = response.code.to_i
      if status >= 200 && status < 300
        parsed = parse_json(body_str)
        cache_store(method_upper, cache_key, parsed, opts)
        bust_cache_for_mutation(method_upper, uri)
        return parsed
      end

      raise ApiError.from_response(response, body_str, request_summary, attempts: attempts)
    end

    def default_user_agent
      "tesote-sdk-rb/#{TesoteSdk::VERSION} (ruby/#{RUBY_VERSION})"
    end

    def build_uri(path, query)
      joined = "#{base_url}/#{version_segment}/#{path.to_s.sub(%r{\A/+}, '')}"
      uri = URI.parse(joined)
      if query && !query.empty?
        uri.query = URI.encode_www_form(stringify_query(query))
      end
      uri
    end

    def build_unversioned_uri(path, query)
      joined = "#{base_url}/#{path.to_s.sub(%r{\A/+}, '')}"
      uri = URI.parse(joined)
      if query && !query.empty?
        uri.query = URI.encode_www_form(stringify_query(query))
      end
      uri
    end

    def stringify_query(query)
      query.each_with_object([]) do |(key, value), acc|
        next if value.nil?

        if value.is_a?(Array)
          value.each { |v| acc << [key.to_s, v.to_s] }
        else
          acc << [key.to_s, value.to_s]
        end
      end
    end

    def perform_with_retries(method, uri, body, opts, request_summary)
      attempt = 0
      idempotency_key = opts[:idempotency_key]
      if MUTATING_METHODS.include?(method) && idempotency_key.nil?
        idempotency_key = SecureRandom.uuid
      end

      loop do
        attempt += 1
        begin
          response, body_str = send_once(method, uri, body, opts, idempotency_key)
        rescue Net::OpenTimeout, Net::ReadTimeout => e
          raise wrap_timeout(e, request_summary, attempt) if !retryable_network?(method) || attempt >= max_attempts

          sleep_for(backoff_delay(attempt))
          next
        rescue OpenSSL::SSL::SSLError => e
          raise TlsError.new("TLS error: #{e.message}", request_summary: request_summary, attempts: attempt, cause: e)
        rescue Errno::ECONNREFUSED, Errno::ECONNRESET, Errno::EHOSTUNREACH,
               Errno::ENETUNREACH, SocketError, EOFError => e
          raise wrap_network(e, request_summary, attempt) if attempt >= max_attempts

          sleep_for(backoff_delay(attempt))
          next
        end

        status = response.code.to_i
        if RETRYABLE_STATUSES.include?(status) && attempt < max_attempts
          sleep_for(retry_after_for(response, attempt))
          next
        end

        return [response, body_str, attempt]
      end
    end

    def send_once(method, uri, body, opts, idempotency_key)
      request_obj = build_request(method, uri, body, opts, idempotency_key)
      http = Net::HTTP.new(uri.host, uri.port)
      http.use_ssl = (uri.scheme == 'https')
      http.open_timeout = open_timeout
      http.read_timeout = read_timeout

      log(:request, method: method, path: uri.request_uri, headers: redact_headers(request_obj))
      response = http.request(request_obj)
      body_str = response.body.to_s
      log(:response, status: response.code, request_id: response['x-request-id'])
      [response, body_str]
    end

    def build_request(method, uri, body, opts, idempotency_key)
      klass = case method
              when 'GET' then Net::HTTP::Get
              when 'POST' then Net::HTTP::Post
              when 'PUT' then Net::HTTP::Put
              when 'PATCH' then Net::HTTP::Patch
              when 'DELETE' then Net::HTTP::Delete
              else raise ArgumentError, "unsupported HTTP method: #{method}"
              end
      req = klass.new(uri.request_uri)
      req['Authorization'] = "Bearer #{api_key}"
      req['Accept'] = 'application/json'
      req['User-Agent'] = user_agent
      req['Idempotency-Key'] = idempotency_key if idempotency_key
      apply_extra_headers(req, opts[:extra_headers])
      if MUTATING_METHODS.include?(method) && !body.nil?
        req['Content-Type'] = 'application/json'
        req.body = body.is_a?(String) ? body : JSON.generate(body)
      end
      req
    end

    def apply_extra_headers(req, headers)
      return if headers.nil?

      headers.each { |k, v| req[k.to_s] = v.to_s }
    end

    def retryable_network?(method)
      # why: only retry timeouts on idempotent methods unless an idempotency key was sent
      return true if %w[GET HEAD OPTIONS PUT DELETE].include?(method)

      false
    end

    def wrap_network(error, request_summary, attempts)
      NetworkError.new("network error: #{error.class}: #{error.message}",
                       request_summary: request_summary,
                       attempts: attempts,
                       cause: error)
    end

    def wrap_timeout(error, request_summary, attempts)
      TimeoutError.new("timeout: #{error.class}: #{error.message}",
                       request_summary: request_summary,
                       attempts: attempts,
                       cause: error)
    end

    def backoff_delay(attempt)
      capped = [base_delay * (2**(attempt - 1)), max_delay].min
      jitter = @randomizer.call(capped)
      [capped + jitter, max_delay].min
    end

    def retry_after_for(response, attempt)
      header = response['retry-after'] || response['Retry-After']
      return Integer(header) if header && header.match?(/\A\d+\z/)

      backoff_delay(attempt)
    end

    def sleep_for(seconds)
      return if seconds.nil? || seconds <= 0

      @sleeper.call(seconds)
    end

    def record_rate_limit(response)
      limit = response['x-ratelimit-limit'] || response['X-RateLimit-Limit']
      remaining = response['x-ratelimit-remaining'] || response['X-RateLimit-Remaining']
      reset = response['x-ratelimit-reset'] || response['X-RateLimit-Reset']
      return if limit.nil? && remaining.nil? && reset.nil?

      @last_rate_limit = RateLimitInfo.new(
        limit: parse_int(limit),
        remaining: parse_int(remaining),
        reset: parse_int(reset)
      )
    end

    def parse_int(value)
      return nil if value.nil?

      Integer(value)
    rescue ArgumentError, TypeError
      nil
    end

    def parse_json(body_str)
      return nil if body_str.nil? || body_str.empty?

      JSON.parse(body_str)
    rescue JSON::ParserError
      body_str
    end

    def build_request_summary(method, uri, body)
      {
        method: method,
        path: uri.path,
        query: redact_query(uri.query),
        body_shape: describe_body_shape(body)
      }
    end

    def redact_query(query_string)
      return nil if query_string.nil? || query_string.empty?

      URI.decode_www_form(query_string).to_h { |(k, v)| [k, redact_secret_value(k, v)] }
    end

    def redact_secret_value(key, value)
      return '[REDACTED]' if key.to_s.match?(/key|token|secret|password/i)

      value
    end

    def redact_headers(req)
      req.each_capitalized.to_h do |key, value|
        if key.casecmp('authorization').zero?
          [key, redact_bearer(value)]
        else
          [key, value]
        end
      end
    end

    def redact_bearer(value)
      return value unless value.is_a?(String) && value.start_with?('Bearer ')

      token = value.sub(/\ABearer\s+/, '')
      last4 = token[-4..] || ''
      "Bearer ****#{last4}"
    end

    def describe_body_shape(body)
      return nil if body.nil?
      return { type: 'string', bytes: body.bytesize } if body.is_a?(String)
      return { type: 'array', items: body.size } if body.is_a?(Array)
      return { type: 'hash', keys: body.keys.map(&:to_s) } if body.is_a?(Hash)

      { type: body.class.name }
    end

    def cache_key_for(method, uri, opts)
      return nil if cache_backend.nil?
      return nil unless method == 'GET'
      return nil if opts[:cache] == false
      return nil unless opts[:cache].is_a?(Hash)

      [method, uri.path, uri.query.to_s, api_key_fingerprint].join('|')
    end

    def cache_lookup(_method, cache_key, opts)
      return nil if cache_key.nil?
      return nil if opts[:cache] == false

      cache_backend.read(cache_key)
    end

    def cache_store(method, cache_key, parsed, opts)
      return if cache_key.nil?
      return unless method == 'GET'
      return unless opts[:cache].is_a?(Hash)

      ttl = opts[:cache][:ttl] || opts[:cache]['ttl']
      return if ttl.nil?

      cache_backend.write(cache_key, parsed, ttl: ttl)
    end

    def bust_cache_for_mutation(method, _uri)
      return if cache_backend.nil?
      return unless MUTATING_METHODS.include?(method)
      return unless cache_backend.respond_to?(:clear)

      cache_backend.clear
    end

    def api_key_fingerprint
      # why: avoid cross-tenant cache bleed; do not store the raw key
      @api_key_fingerprint ||= begin
        last4 = api_key[-4..] || api_key
        "key-#{last4}"
      end
    end

    def log(event, payload)
      return if logger.nil?

      # why: a misbehaving user logger must never break the request path
      begin
        logger.call(event, payload)
      rescue NoMethodError, ArgumentError, TypeError
        nil
      end
    end
  end
end
