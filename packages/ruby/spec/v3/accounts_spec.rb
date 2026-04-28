require 'spec_helper'

RSpec.describe TesoteSdk::V3::Accounts do
  let(:api_key) { 'sk_test_abcd1234' }
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) do
    TesoteSdk::V3::Client.new(
      api_key: api_key,
      base_url: base_url,
      base_delay: 0.0,
      max_delay: 0.0,
      sleeper: ->(_s) {},
      randomizer: ->(_) { 0.0 }
    )
  end

  describe '#list' do
    it 'GETs /v3/accounts and parses the JSON' do
      stub = stub_request(:get, "#{base_url}/v3/accounts")
             .with(headers: { 'Authorization' => "Bearer #{api_key}" })
             .to_return(status: 200,
                        body: '{"data":[{"id":"acct_1"}],"meta":{"cursor":null}}',
                        headers: { 'Content-Type' => 'application/json',
                                   'X-Request-Id' => 'req_a',
                                   'X-RateLimit-Remaining' => '199' })

      result = client.accounts.list
      expect(stub).to have_been_requested
      expect(result['data']).to eq([{ 'id' => 'acct_1' }])
      expect(client.last_request_id).to eq('req_a')
      expect(client.last_rate_limit.remaining).to eq(199)
    end

    it 'forwards query params' do
      stub = stub_request(:get, "#{base_url}/v3/accounts")
             .with(query: { 'limit' => '50', 'cursor' => 'abc' })
             .to_return(status: 200, body: '{"data":[]}')
      client.accounts.list({ limit: 50, cursor: 'abc' })
      expect(stub).to have_been_requested
    end
  end

  describe '#get' do
    it 'GETs /v3/accounts/:id' do
      stub = stub_request(:get, "#{base_url}/v3/accounts/acct_42")
             .to_return(status: 200, body: '{"id":"acct_42","balance":1234}')
      result = client.accounts.get('acct_42')
      expect(stub).to have_been_requested
      expect(result['id']).to eq('acct_42')
    end

    it 'raises ArgumentError on blank id' do
      expect { client.accounts.get('') }.to raise_error(ArgumentError)
    end

    it 'maps a 401 response to UnauthorizedError with all required fields' do
      stub_request(:get, "#{base_url}/v3/accounts/acct_42")
        .to_return(status: 401,
                   body: '{"error":"bad key","error_code":"UNAUTHORIZED","error_id":"e1"}',
                   headers: { 'X-Request-Id' => 'req_z' })
      expect { client.accounts.get('acct_42') }.to raise_error(TesoteSdk::UnauthorizedError) do |err|
        expect(err.error_code).to eq('UNAUTHORIZED')
        expect(err.http_status).to eq(401)
        expect(err.request_id).to eq('req_z')
        expect(err.error_id).to eq('e1')
        expect(err.request_summary[:method]).to eq('GET')
        expect(err.request_summary[:path]).to eq('/api/v3/accounts/acct_42')
      end
    end
  end

  describe 'unwired methods' do
    it 'raises NotImplementedError for stubs' do
      expect { client.accounts.sync('acct_1') }.to raise_error(NotImplementedError)
      expect { client.transactions.list_for_account('acct_1') }.to raise_error(NotImplementedError)
      expect { client.webhooks.list }.to raise_error(NotImplementedError)
    end
  end

  describe 'webhook signature helper' do
    it 'verifies a valid signature' do
      body = '{"event":"x"}'
      secret = 'whsec_test'
      sig = OpenSSL::HMAC.hexdigest('SHA256', secret, body)
      expect(TesoteSdk::V3.verify_webhook_signature(body: body, signature_header: sig, secret: secret)).to be(true)
    end

    it 'raises on mismatch' do
      expect do
        TesoteSdk::V3.verify_webhook_signature(body: 'b', signature_header: 'deadbeef' * 8, secret: 'k')
      end.to raise_error(TesoteSdk::Error)
    end
  end
end
