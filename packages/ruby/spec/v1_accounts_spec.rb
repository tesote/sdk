require 'spec_helper'

RSpec.describe TesoteSdk::V1::Accounts do
  let(:api_key) { 'test_api_key' }
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) { TesoteSdk::V1::Client.new(api_key: api_key, base_url: base_url, base_delay: 0.0, max_delay: 0.0) }

  describe '#list' do
    it 'returns AccountList of typed Account models' do
      payload = {
        total: 1,
        accounts: [{
          id: 'a_1',
          name: 'Checking',
          data: { masked_account_number: '1234', currency: 'VES' },
          bank: { name: 'Banesco' },
          legal_entity: { id: 'le_1', legal_name: 'Acme' },
          tesote_created_at: '2026-04-01T00:00:00Z',
          tesote_updated_at: '2026-04-02T00:00:00Z'
        }],
        pagination: { current_page: 1, per_page: 50, total_pages: 1, total_count: 1 }
      }
      stub_request(:get, "#{base_url}/v1/accounts")
        .to_return(status: 200, body: payload.to_json)

      result = client.accounts.list
      expect(result).to be_a(TesoteSdk::Models::AccountList)
      expect(result.total).to eq(1)
      expect(result.accounts.size).to eq(1)
      acct = result.accounts.first
      expect(acct).to be_a(TesoteSdk::Models::Account)
      expect(acct.id).to eq('a_1')
      expect(acct.bank).to be_a(TesoteSdk::Models::Bank)
      expect(acct.bank.name).to eq('Banesco')
    end
  end

  describe '#get' do
    it 'returns Account model' do
      stub_request(:get, "#{base_url}/v1/accounts/a_1")
        .to_return(status: 200,
                   body: { id: 'a_1', name: 'Checking', data: { currency: 'VES' }, bank: { name: 'B' } }.to_json)
      result = client.accounts.get('a_1')
      expect(result).to be_a(TesoteSdk::Models::Account)
      expect(result.id).to eq('a_1')
      expect(result.data.currency).to eq('VES')
    end

    it 'maps 404 ACCOUNT_NOT_FOUND to AccountNotFoundError' do
      stub_request(:get, "#{base_url}/v1/accounts/missing")
        .to_return(status: 404, body: { error: 'nope', error_code: 'ACCOUNT_NOT_FOUND' }.to_json)
      expect { client.accounts.get('missing') }.to raise_error(TesoteSdk::AccountNotFoundError)
    end

    it 'raises ArgumentError on blank id' do
      expect { client.accounts.get('') }.to raise_error(ArgumentError)
    end
  end

  describe '#list_transactions' do
    it 'returns TransactionList with cursor pagination fields' do
      payload = {
        total: 2,
        transactions: [
          { id: 't_1', status: 'posted',
            data: { amount_cents: 100, currency: 'VES', description: 'a', transaction_date: '2026-04-01' },
            tesote_imported_at: '2026-04-01', tesote_updated_at: '2026-04-01',
            transaction_categories: [], counterparty: { name: 'Vendor' } }
        ],
        pagination: { has_more: true, per_page: 50, after_id: 't_1', before_id: 't_1' }
      }
      stub_request(:get, "#{base_url}/v1/accounts/a_1/transactions")
        .with(query: hash_including('start_date' => '2026-04-01'))
        .to_return(status: 200, body: payload.to_json)

      result = client.accounts.list_transactions('a_1', { start_date: '2026-04-01' })
      expect(result).to be_a(TesoteSdk::Models::TransactionList)
      expect(result.transactions.first.counterparty.name).to eq('Vendor')
      expect(result.pagination.has_more).to eq(true)
    end

    it 'maps INVALID_DATE_RANGE to InvalidDateRangeError' do
      stub_request(:get, "#{base_url}/v1/accounts/a_1/transactions")
        .to_return(status: 422, body: { error: 'bad', error_code: 'INVALID_DATE_RANGE' }.to_json)
      expect { client.accounts.list_transactions('a_1') }.to raise_error(TesoteSdk::InvalidDateRangeError)
    end
  end

  describe '#each_transaction_page' do
    it 'walks cursor pages until has_more is false' do
      page1 = {
        total: 2,
        transactions: [{ id: 't_1', status: 'posted', data: {}, tesote_imported_at: nil, tesote_updated_at: nil,
                         transaction_categories: [], counterparty: nil }],
        pagination: { has_more: true, per_page: 1, after_id: 't_1', before_id: 't_1' }
      }.to_json
      page2 = {
        total: 2,
        transactions: [{ id: 't_2', status: 'posted', data: {}, tesote_imported_at: nil, tesote_updated_at: nil,
                         transaction_categories: [], counterparty: nil }],
        pagination: { has_more: false, per_page: 1, after_id: 't_2', before_id: 't_2' }
      }.to_json

      stub_request(:get, %r{/v1/accounts/a_1/transactions})
        .to_return({ status: 200, body: page1 }, { status: 200, body: page2 })

      pages = client.accounts.each_transaction_page('a_1').to_a
      expect(pages.size).to eq(2)
      expect(pages.last.transactions.first.id).to eq('t_2')
    end
  end
end
