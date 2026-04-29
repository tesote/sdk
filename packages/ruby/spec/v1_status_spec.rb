require 'spec_helper'

RSpec.describe TesoteSdk::V1::Status do
  let(:api_key) { 'test_api_key' }
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) { TesoteSdk::V1::Client.new(api_key: api_key, base_url: base_url, base_delay: 0.0, max_delay: 0.0) }

  describe '#status' do
    it 'GETs /status (unversioned)' do
      stub_request(:get, "#{base_url}/status")
        .to_return(status: 200, body: { status: 'ok', authenticated: false }.to_json)

      result = client.status.status
      expect(result).to be_a(TesoteSdk::Models::StatusResult)
      expect(result.status).to eq('ok')
      expect(result.authenticated).to eq(false)
    end
  end

  describe '#whoami' do
    it 'GETs /whoami (unversioned) and returns Whoami' do
      stub_request(:get, "#{base_url}/whoami")
        .to_return(status: 200, body: { client: { id: 'c_1', name: 'Acme', type: 'workspace' } }.to_json)

      result = client.status.whoami
      expect(result).to be_a(TesoteSdk::Models::Whoami)
      expect(result.client['id']).to eq('c_1')
    end

    it 'maps 401 UNAUTHORIZED to UnauthorizedError' do
      stub_request(:get, "#{base_url}/whoami")
        .to_return(status: 401, body: { error: 'no auth', error_code: 'UNAUTHORIZED' }.to_json)

      expect { client.status.whoami }.to raise_error(TesoteSdk::UnauthorizedError)
    end
  end
end
