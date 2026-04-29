require 'spec_helper'

RSpec.describe TesoteSdk::V1::Transactions do
  let(:api_key) { 'test_api_key' }
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) { TesoteSdk::V1::Client.new(api_key: api_key, base_url: base_url, base_delay: 0.0, max_delay: 0.0) }

  describe '#get' do
    it 'returns a typed Transaction' do
      payload = {
        id: 't_1',
        status: 'posted',
        data: { amount_cents: 1000, currency: 'VES', description: 'lunch', transaction_date: '2026-04-01' },
        tesote_imported_at: '2026-04-01T00:00:00Z',
        tesote_updated_at: '2026-04-01T00:00:00Z',
        transaction_categories: [{ name: 'food', external_category_code: 'FOOD',
                                   created_at: '2026-04-01', updated_at: '2026-04-01' }],
        counterparty: { name: 'Cafe' }
      }
      stub_request(:get, "#{base_url}/v1/transactions/t_1")
        .to_return(status: 200, body: payload.to_json)

      result = client.transactions.get('t_1')
      expect(result).to be_a(TesoteSdk::Models::Transaction)
      expect(result.data).to be_a(TesoteSdk::Models::TransactionData)
      expect(result.data.amount_cents).to eq(1000)
      expect(result.transaction_categories.first.name).to eq('food')
    end

    it 'maps 404 TRANSACTION_NOT_FOUND to TransactionNotFoundError' do
      stub_request(:get, "#{base_url}/v1/transactions/missing")
        .to_return(status: 404, body: { error_code: 'TRANSACTION_NOT_FOUND', error: 'nope' }.to_json)
      expect { client.transactions.get('missing') }.to raise_error(TesoteSdk::TransactionNotFoundError)
    end

    it 'raises ArgumentError on blank id' do
      expect { client.transactions.get(nil) }.to raise_error(ArgumentError)
    end
  end

  describe '#list_for_account' do
    it 'returns a TransactionList' do
      stub_request(:get, "#{base_url}/v1/accounts/a_1/transactions")
        .to_return(status: 200, body: { total: 0, transactions: [], pagination: { has_more: false } }.to_json)
      result = client.transactions.list_for_account('a_1')
      expect(result).to be_a(TesoteSdk::Models::TransactionList)
      expect(result.transactions).to eq([])
    end
  end
end
