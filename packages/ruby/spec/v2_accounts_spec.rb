require 'spec_helper'

RSpec.describe TesoteSdk::V2::Accounts do
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) { TesoteSdk::V2::Client.new(api_key: 'k', base_url: base_url, base_delay: 0.0, max_delay: 0.0) }

  describe '#list' do
    it 'returns AccountList' do
      stub_request(:get, "#{base_url}/v2/accounts")
        .to_return(status: 200, body: { total: 0, accounts: [], pagination: {} }.to_json)
      expect(client.accounts.list).to be_a(TesoteSdk::Models::AccountList)
    end
  end

  describe '#get' do
    it 'returns an Account' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1")
        .to_return(status: 200, body: { id: 'a_1', name: 'X', data: {}, bank: {} }.to_json)
      expect(client.accounts.get('a_1')).to be_a(TesoteSdk::Models::Account)
    end

    it 'maps 404 ACCOUNT_NOT_FOUND' do
      stub_request(:get, "#{base_url}/v2/accounts/missing")
        .to_return(status: 404, body: { error_code: 'ACCOUNT_NOT_FOUND' }.to_json)
      expect { client.accounts.get('missing') }.to raise_error(TesoteSdk::AccountNotFoundError)
    end
  end

  describe '#sync' do
    it 'POSTs and returns SyncStartResult' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/sync")
        .with(headers: { 'Content-Type' => 'application/json' })
        .to_return(status: 202, body: { message: 'started', sync_session_id: 'ss_1',
                                        status: 'pending', started_at: '2026-04-28' }.to_json)
      result = client.accounts.sync('a_1')
      expect(result).to be_a(TesoteSdk::Models::SyncStartResult)
      expect(result.sync_session_id).to eq('ss_1')
    end

    it 'auto-generates an Idempotency-Key for the POST' do
      header_seen = nil
      stub_request(:post, "#{base_url}/v2/accounts/a_1/sync")
        .with do |req|
          header_seen = req.headers['Idempotency-Key']
          true
        end
        .to_return(status: 202, body: { message: 's', sync_session_id: 'ss', status: 'pending',
                                        started_at: 'now' }.to_json)
      client.accounts.sync('a_1')
      expect(header_seen).to match(/\A[0-9a-f-]{36}\z/i)
    end

    it 'maps 409 SYNC_IN_PROGRESS' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/sync")
        .to_return(status: 409, body: { error_code: 'SYNC_IN_PROGRESS', error: 'busy' }.to_json)
      expect { client.accounts.sync('a_1') }.to raise_error(TesoteSdk::SyncInProgressError)
    end

    it 'maps 503 BANK_UNDER_MAINTENANCE' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/sync")
        .to_return(status: 503, body: { error_code: 'BANK_UNDER_MAINTENANCE' }.to_json,
                   headers: { 'Retry-After' => '120' })
      tight = TesoteSdk::V2::Client.new(api_key: 'k', base_url: base_url, max_attempts: 1,
                                        base_delay: 0.0, max_delay: 0.0, sleeper: ->(_) {})
      expect { tight.accounts.sync('a_1') }.to raise_error(TesoteSdk::BankUnderMaintenanceError)
    end

    it 'maps 429 SYNC_RATE_LIMIT_EXCEEDED' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/sync")
        .to_return(status: 429, body: { error_code: 'SYNC_RATE_LIMIT_EXCEEDED', retry_after: 60 }.to_json,
                   headers: { 'Retry-After' => '60' })
      tight = TesoteSdk::V2::Client.new(api_key: 'k', base_url: base_url, max_attempts: 1,
                                        base_delay: 0.0, max_delay: 0.0, sleeper: ->(_) {})
      expect { tight.accounts.sync('a_1') }.to raise_error(TesoteSdk::SyncRateLimitExceededError)
    end
  end
end
