require 'spec_helper'

RSpec.describe TesoteSdk::V2::SyncSessions do
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) { TesoteSdk::V2::Client.new(api_key: 'k', base_url: base_url, base_delay: 0.0, max_delay: 0.0) }

  describe '#list' do
    it 'returns OffsetPage of SyncSession' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/sync_sessions")
        .to_return(status: 200,
                   body: { sync_sessions: [{ id: 'ss_1', status: 'completed', started_at: '2026-04-01' }],
                           limit: 50, offset: 0, has_more: false }.to_json)
      page = client.sync_sessions.list('a_1')
      expect(page).to be_a(TesoteSdk::Models::OffsetPage)
      expect(page.items.first).to be_a(TesoteSdk::Models::SyncSession)
      expect(page.has_more).to eq(false)
    end

    it 'maps BANK_CONNECTION_NOT_FOUND' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/sync_sessions")
        .to_return(status: 404, body: { error_code: 'BANK_CONNECTION_NOT_FOUND' }.to_json)
      expect { client.sync_sessions.list('a_1') }.to raise_error(TesoteSdk::BankConnectionNotFoundError)
    end
  end

  describe '#get' do
    it 'returns a SyncSession' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/sync_sessions/ss_1")
        .to_return(status: 200, body: { id: 'ss_1', status: 'completed', started_at: 'now' }.to_json)
      expect(client.sync_sessions.get('a_1', 'ss_1')).to be_a(TesoteSdk::Models::SyncSession)
    end

    it 'maps SYNC_SESSION_NOT_FOUND' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/sync_sessions/missing")
        .to_return(status: 404, body: { error_code: 'SYNC_SESSION_NOT_FOUND' }.to_json)
      expect { client.sync_sessions.get('a_1', 'missing') }.to raise_error(TesoteSdk::SyncSessionNotFoundError)
    end
  end

  describe '#each_page (offset)' do
    it 'walks until has_more false' do
      page1 = { sync_sessions: [{ id: 'ss_1' }], limit: 1, offset: 0, has_more: true }.to_json
      page2 = { sync_sessions: [{ id: 'ss_2' }], limit: 1, offset: 1, has_more: false }.to_json
      stub_request(:get, %r{/v2/accounts/a_1/sync_sessions})
        .to_return({ status: 200, body: page1 }, { status: 200, body: page2 })

      pages = client.sync_sessions.each_page('a_1', {}, page_size: 1).to_a
      expect(pages.size).to eq(2)
    end
  end
end
