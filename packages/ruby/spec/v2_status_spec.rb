require 'spec_helper'

RSpec.describe TesoteSdk::V2::Status do
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) { TesoteSdk::V2::Client.new(api_key: 'k', base_url: base_url, base_delay: 0.0, max_delay: 0.0) }

  it 'GETs /v2/status' do
    stub_request(:get, "#{base_url}/v2/status")
      .to_return(status: 200, body: { status: 'ok', authenticated: false }.to_json)
    result = client.status.status
    expect(result).to be_a(TesoteSdk::Models::StatusResult)
    expect(result.status).to eq('ok')
  end

  it 'GETs /v2/whoami' do
    stub_request(:get, "#{base_url}/v2/whoami")
      .to_return(status: 200, body: { client: { id: 'c_1' } }.to_json)
    result = client.status.whoami
    expect(result).to be_a(TesoteSdk::Models::Whoami)
  end
end
