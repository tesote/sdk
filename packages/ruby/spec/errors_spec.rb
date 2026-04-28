require 'spec_helper'

RSpec.describe TesoteSdk::ApiError do
  let(:summary) do
    { method: 'GET', path: '/api/v3/accounts', query: { 'q' => '1' }, body_shape: nil }
  end

  def fake_response(status:, headers: {}, message: 'Error')
    response = Net::HTTPResponse.send(:response_class, status.to_s).new('1.1', status.to_s, message)
    headers.each { |k, v| response[k] = v }
    response
  end

  describe '.from_response' do
    {
      'UNAUTHORIZED' => [401, TesoteSdk::UnauthorizedError],
      'API_KEY_REVOKED' => [401, TesoteSdk::ApiKeyRevokedError],
      'WORKSPACE_SUSPENDED' => [403, TesoteSdk::WorkspaceSuspendedError],
      'ACCOUNT_DISABLED' => [403, TesoteSdk::AccountDisabledError],
      'HISTORY_SYNC_FORBIDDEN' => [403, TesoteSdk::HistorySyncForbiddenError],
      'MUTATION_CONFLICT' => [409, TesoteSdk::MutationDuringPaginationError],
      'UNPROCESSABLE_CONTENT' => [422, TesoteSdk::UnprocessableContentError],
      'INVALID_DATE_RANGE' => [422, TesoteSdk::InvalidDateRangeError],
      'RATE_LIMIT_EXCEEDED' => [429, TesoteSdk::RateLimitExceededError]
    }.each do |code, (status, klass)|
      it "maps error_code #{code} → #{klass}" do
        response = fake_response(status: status, headers: { 'X-Request-Id' => 'req_1' })
        body = { error: 'msg', error_code: code, error_id: 'eid_1' }.to_json
        err = described_class.from_response(response, body, summary, attempts: 2)
        expect(err).to be_a(klass)
        expect(err.error_code).to eq(code)
        expect(err.message).to eq('msg')
        expect(err.http_status).to eq(status)
        expect(err.request_id).to eq('req_1')
        expect(err.error_id).to eq('eid_1')
        expect(err.attempts).to eq(2)
        expect(err.response_body).to eq(body)
        expect(err.request_summary).to eq(summary)
      end
    end

    it 'maps 503 with no error_code to ServiceUnavailableError' do
      response = fake_response(status: 503)
      err = described_class.from_response(response, '{}', summary)
      expect(err).to be_a(TesoteSdk::ServiceUnavailableError)
    end

    it 'falls back to ApiError for unknown error_code' do
      response = fake_response(status: 418)
      err = described_class.from_response(response, '{"error_code":"TEAPOT"}', summary)
      expect(err.class).to eq(TesoteSdk::ApiError)
      expect(err.error_code).to eq('TEAPOT')
    end

    it 'parses retry_after from header (integer string)' do
      response = fake_response(status: 429, headers: { 'Retry-After' => '17' })
      body = '{"error_code":"RATE_LIMIT_EXCEEDED"}'
      err = described_class.from_response(response, body, summary)
      expect(err.retry_after).to eq(17)
    end

    it 'falls back to envelope retry_after when header missing' do
      response = fake_response(status: 429)
      body = '{"error_code":"RATE_LIMIT_EXCEEDED","retry_after":42}'
      err = described_class.from_response(response, body, summary)
      expect(err.retry_after).to eq(42)
    end

    it 'tolerates non-JSON body' do
      response = fake_response(status: 500, message: 'Internal Server Error')
      err = described_class.from_response(response, '<html>boom</html>', summary)
      expect(err.message).to include('500')
      expect(err.response_body).to eq('<html>boom</html>')
    end
  end
end

RSpec.describe TesoteSdk::Transport do
  it 'redacts the bearer token in the logger payload' do
    captured = []
    transport = described_class.new(
      api_key: 'sk_live_1234567890abcd',
      version_segment: 'v3',
      base_url: 'https://equipo.tesote.com/api',
      base_delay: 0.0, max_delay: 0.0,
      sleeper: ->(_s) {},
      randomizer: ->(_) { 0.0 },
      logger: ->(event, payload) { captured << [event, payload] }
    )
    stub_request(:get, 'https://equipo.tesote.com/api/v3/accounts')
      .to_return(status: 200, body: '[]')

    transport.request('GET', 'accounts')

    request_event = captured.find { |(event, _p)| event == :request }
    expect(request_event).not_to be_nil
    auth_value = request_event[1][:headers]['Authorization']
    expect(auth_value).to eq('Bearer ****abcd')
    expect(auth_value).not_to include('sk_live_1234567890abcd')
  end
end
