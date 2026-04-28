require 'spec_helper'

RSpec.describe TesoteSdk::Transport do
  let(:api_key) { 'test_api_key_abcd1234' }
  let(:base_url) { 'https://equipo.tesote.com/api' }

  def build_transport(**overrides)
    described_class.new(
      api_key: api_key,
      version_segment: 'v3',
      base_url: base_url,
      max_attempts: overrides.delete(:max_attempts) || 3,
      base_delay: 0.0,
      max_delay: 0.0,
      sleeper: ->(_s) {},
      randomizer: ->(_) { 0.0 },
      **overrides
    )
  end

  describe 'auth + headers' do
    it 'injects bearer token, accept, user-agent' do
      stub = stub_request(:get, "#{base_url}/v3/accounts")
             .with(headers: {
                     'Authorization' => "Bearer #{api_key}",
                     'Accept' => 'application/json',
                     'User-Agent' => %r{\Atesote-sdk-rb/\d+\.\d+\.\d+ \(ruby/.+\)\z}
                   })
             .to_return(status: 200, body: '{"data":[]}', headers: { 'Content-Type' => 'application/json' })

      result = build_transport.request('GET', 'accounts')
      expect(result).to eq({ 'data' => [] })
      expect(stub).to have_been_requested
    end

    it 'raises ConfigError when api_key is empty' do
      expect { described_class.new(api_key: '', version_segment: 'v3') }
        .to raise_error(TesoteSdk::ConfigError)
    end

    it 'sets Content-Type for mutating requests with body' do
      stub = stub_request(:post, "#{base_url}/v3/accounts/acct_1/sync")
             .with(headers: { 'Content-Type' => 'application/json' },
                   body: { reason: 'manual' }.to_json)
             .to_return(status: 200, body: '{"ok":true}')
      build_transport.request('POST', 'accounts/acct_1/sync', body: { reason: 'manual' })
      expect(stub).to have_been_requested
    end
  end

  describe 'idempotency' do
    it 'auto-generates an Idempotency-Key for POST when caller does not pass one' do
      header_seen = nil
      capture = lambda do |req|
        header_seen = req.headers['Idempotency-Key']
        true
      end
      stub = stub_request(:post, "#{base_url}/v3/batches")
             .with(&capture)
             .to_return(status: 200, body: '{}')

      build_transport.request('POST', 'batches', body: {})
      expect(stub).to have_been_requested
      expect(header_seen).to match(/\A[0-9a-f-]{36}\z/i)
    end

    it 'forwards a caller-supplied idempotency key verbatim' do
      stub = stub_request(:post, "#{base_url}/v3/batches")
             .with(headers: { 'Idempotency-Key' => 'caller-supplied-key' })
             .to_return(status: 200, body: '{}')
      build_transport.request('POST', 'batches', body: {}, opts: { idempotency_key: 'caller-supplied-key' })
      expect(stub).to have_been_requested
    end

    it 'does not send an Idempotency-Key on GET' do
      stub = stub_request(:get, "#{base_url}/v3/accounts")
             .with { |req| !req.headers.key?('Idempotency-Key') }
             .to_return(status: 200, body: '{"data":[]}')
      build_transport.request('GET', 'accounts')
      expect(stub).to have_been_requested
    end
  end

  describe 'rate-limit headers' do
    it 'captures X-RateLimit-* headers into last_rate_limit' do
      stub_request(:get, "#{base_url}/v3/accounts")
        .to_return(status: 200, body: '[]', headers: {
                     'X-RateLimit-Limit' => '200',
                     'X-RateLimit-Remaining' => '197',
                     'X-RateLimit-Reset' => '1700000000'
                   })

      transport = build_transport
      transport.request('GET', 'accounts')
      info = transport.last_rate_limit
      expect(info.limit).to eq(200)
      expect(info.remaining).to eq(197)
      expect(info.reset).to eq(1_700_000_000)
    end

    it 'propagates X-Request-Id into last_request_id' do
      stub_request(:get, "#{base_url}/v3/accounts")
        .to_return(status: 200, body: '[]', headers: { 'X-Request-Id' => 'req_xyz' })
      transport = build_transport
      transport.request('GET', 'accounts')
      expect(transport.last_request_id).to eq('req_xyz')
    end
  end

  describe 'retries' do
    it 'retries on 503 then succeeds' do
      stub = stub_request(:get, "#{base_url}/v3/accounts")
             .to_return({ status: 503, body: '{}' },
                        { status: 200, body: '{"ok":true}' })
      result = build_transport(max_attempts: 3).request('GET', 'accounts')
      expect(result).to eq({ 'ok' => true })
      expect(stub).to have_been_requested.twice
    end

    it 'raises RateLimitExceededError after exhausting retries on 429' do
      stub_request(:get, "#{base_url}/v3/accounts")
        .to_return(status: 429,
                   body: '{"error":"rate","error_code":"RATE_LIMIT_EXCEEDED","retry_after":1}',
                   headers: { 'Retry-After' => '1', 'X-Request-Id' => 'req_rl' })

      expect do
        build_transport(max_attempts: 2).request('GET', 'accounts')
      end.to raise_error(TesoteSdk::RateLimitExceededError) do |err|
        expect(err.attempts).to eq(2)
        expect(err.retry_after).to eq(1)
        expect(err.request_id).to eq('req_rl')
      end
    end

    it 'does not retry on 400-class errors other than 429' do
      stub = stub_request(:get, "#{base_url}/v3/accounts/acct_x")
             .to_return(status: 422,
                        body: '{"error":"bad","error_code":"UNPROCESSABLE_CONTENT"}')
      expect { build_transport.request('GET', 'accounts/acct_x') }
        .to raise_error(TesoteSdk::UnprocessableContentError)
      expect(stub).to have_been_requested.once
    end

    it 'wraps connection errors in NetworkError after retries' do
      stub_request(:get, "#{base_url}/v3/accounts")
        .to_raise(Errno::ECONNRESET)
      expect { build_transport(max_attempts: 2).request('GET', 'accounts') }
        .to raise_error(TesoteSdk::NetworkError) do |err|
          expect(err.attempts).to eq(2)
          expect(err.cause).to be_a(Errno::ECONNRESET)
        end
    end

    it 'wraps timeouts on idempotent methods in TimeoutError after retries' do
      stub_request(:get, "#{base_url}/v3/accounts")
        .to_raise(Net::ReadTimeout)
      expect { build_transport(max_attempts: 2).request('GET', 'accounts') }
        .to raise_error(TesoteSdk::TimeoutError)
    end

    it 'does not retry timeouts on POST without an idempotency key path' do
      stub_request(:post, "#{base_url}/v3/batches")
        .to_raise(Net::ReadTimeout)
      expect { build_transport(max_attempts: 3).request('POST', 'batches', body: {}) }
        .to raise_error(TesoteSdk::TimeoutError) do |err|
          expect(err.attempts).to eq(1)
        end
    end
  end

  describe 'cache backend' do
    it 'returns cached body on a second GET when ttl is set' do
      backend = TesoteSdk::CacheBackend.new
      stub = stub_request(:get, "#{base_url}/v3/accounts")
             .to_return(status: 200, body: '{"data":[1]}')
      transport = build_transport(cache_backend: backend)
      a = transport.request('GET', 'accounts', opts: { cache: { ttl: 60 } })
      b = transport.request('GET', 'accounts', opts: { cache: { ttl: 60 } })
      expect(a).to eq(b)
      expect(stub).to have_been_requested.once
    end

    it 'bypasses cache when opts[:cache] is false' do
      backend = TesoteSdk::CacheBackend.new
      stub = stub_request(:get, "#{base_url}/v3/accounts")
             .to_return(status: 200, body: '{"data":[1]}')
      transport = build_transport(cache_backend: backend)
      transport.request('GET', 'accounts', opts: { cache: { ttl: 60 } })
      transport.request('GET', 'accounts', opts: { cache: false })
      expect(stub).to have_been_requested.twice
    end
  end

  describe 'request summary + redaction' do
    it 'redacts token-like query params in the request summary' do
      stub_request(:get, "#{base_url}/v3/accounts").with(query: { api_key: 'secret' })
                                                   .to_return(status: 422, body: '{"error_code":"UNPROCESSABLE_CONTENT"}')
      expect { build_transport.request('GET', 'accounts', query: { api_key: 'secret' }) }
        .to raise_error(TesoteSdk::ApiError) do |err|
          expect(err.request_summary[:query]).to eq({ 'api_key' => '[REDACTED]' })
        end
    end
  end
end
