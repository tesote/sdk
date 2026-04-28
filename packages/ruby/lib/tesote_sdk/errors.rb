require 'json'

module TesoteSdk
  # Base class. Catch-all only as last resort.
  class Error < StandardError
    attr_reader :error_code, :http_status, :request_id, :error_id,
                :retry_after, :response_body, :request_summary, :attempts

    def initialize(message,
                   error_code: nil,
                   http_status: nil,
                   request_id: nil,
                   error_id: nil,
                   retry_after: nil,
                   response_body: nil,
                   request_summary: nil,
                   attempts: nil,
                   cause: nil)
      super(message)
      @error_code = error_code
      @http_status = http_status
      @request_id = request_id
      @error_id = error_id
      @retry_after = retry_after
      @response_body = response_body
      @request_summary = request_summary
      @attempts = attempts
      @cause_override = cause
    end

    def cause
      @cause_override || super
    end

    def to_h
      {
        error_code: error_code,
        message: message,
        http_status: http_status,
        request_id: request_id,
        error_id: error_id,
        retry_after: retry_after,
        response_body: response_body,
        request_summary: request_summary,
        attempts: attempts
      }
    end

    def inspect
      parts = [
        "code=#{error_code.inspect}",
        "status=#{http_status.inspect}",
        "request_id=#{request_id.inspect}",
        "attempts=#{attempts.inspect}"
      ].join(' ')
      "#<#{self.class.name}: #{message} (#{parts})>"
    end
  end

  # Bad SDK config; raised at construction.
  class ConfigError < Error; end

  # Calling a method whose upstream endpoint is gone in this version.
  class EndpointRemovedError < Error; end

  # Server returned a usable HTTP response with an error.
  class ApiError < Error
    # Map error_code → typed subclass. Unknown codes fall back to ApiError.
    CODE_REGISTRY = {} # rubocop:disable Style/MutableConstant -- registered via .register at load time

    def self.register(code, klass)
      CODE_REGISTRY[code] = klass
    end

    # why: single dispatcher used by Transport — keeps mapping logic out of caller paths.
    def self.from_response(response, body, request_summary, attempts: 1)
      parsed = parse_body(body)
      envelope = parsed.is_a?(Hash) ? parsed : {}
      error_code = envelope['error_code']
      message = envelope['error'] || synthesize_message(response)
      error_id = envelope['error_id']
      retry_after_value = parse_retry_after(response, envelope['retry_after'])
      request_id = response['x-request-id'] || response['X-Request-Id']
      http_status = response.code.to_i

      klass = pick_class(http_status, error_code)
      klass.new(
        message,
        error_code: error_code,
        http_status: http_status,
        request_id: request_id,
        error_id: error_id,
        retry_after: retry_after_value,
        response_body: body,
        request_summary: request_summary,
        attempts: attempts
      )
    end

    def self.pick_class(http_status, error_code)
      return CODE_REGISTRY[error_code] if error_code && CODE_REGISTRY.key?(error_code)
      return ServiceUnavailableError if http_status == 503

      ApiError
    end

    def self.parse_body(body)
      return nil if body.nil? || body.empty?

      JSON.parse(body)
    rescue JSON::ParserError
      nil
    end

    def self.synthesize_message(response)
      "HTTP #{response.code} #{response.message}"
    end

    def self.parse_retry_after(response, envelope_value)
      header = response['retry-after'] || response['Retry-After']
      return Integer(header) if header && header.match?(/\A\d+\z/)
      return Integer(envelope_value) if envelope_value.is_a?(Integer)
      return Integer(envelope_value) if envelope_value.is_a?(String) && envelope_value.match?(/\A\d+\z/)

      nil
    end
  end

  class UnauthorizedError < ApiError; end
  class ApiKeyRevokedError < ApiError; end
  class WorkspaceSuspendedError < ApiError; end
  class AccountDisabledError < ApiError; end
  class HistorySyncForbiddenError < ApiError; end
  class MutationDuringPaginationError < ApiError; end
  class UnprocessableContentError < ApiError; end
  class InvalidDateRangeError < ApiError; end
  class RateLimitExceededError < ApiError; end
  class ServiceUnavailableError < ApiError; end

  ApiError.register('UNAUTHORIZED', UnauthorizedError)
  ApiError.register('API_KEY_REVOKED', ApiKeyRevokedError)
  ApiError.register('WORKSPACE_SUSPENDED', WorkspaceSuspendedError)
  ApiError.register('ACCOUNT_DISABLED', AccountDisabledError)
  ApiError.register('HISTORY_SYNC_FORBIDDEN', HistorySyncForbiddenError)
  ApiError.register('MUTATION_CONFLICT', MutationDuringPaginationError)
  ApiError.register('UNPROCESSABLE_CONTENT', UnprocessableContentError)
  ApiError.register('INVALID_DATE_RANGE', InvalidDateRangeError)
  ApiError.register('RATE_LIMIT_EXCEEDED', RateLimitExceededError)

  # Transport-level failures: no usable HTTP response.
  class TransportError < Error; end
  class NetworkError < TransportError; end
  class TimeoutError < TransportError; end
  class TlsError < TransportError; end
end
