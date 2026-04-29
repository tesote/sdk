require 'spec_helper'

RSpec.describe TesoteSdk::V2::Transactions do
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) { TesoteSdk::V2::Client.new(api_key: 'k', base_url: base_url, base_delay: 0.0, max_delay: 0.0) }

  describe '#list_for_account' do
    it 'returns TransactionList' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/transactions")
        .to_return(status: 200, body: { total: 0, transactions: [], pagination: { has_more: false } }.to_json)
      expect(client.transactions.list_for_account('a_1')).to be_a(TesoteSdk::Models::TransactionList)
    end

    it 'maps INVALID_DATE_RANGE' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/transactions")
        .to_return(status: 422, body: { error_code: 'INVALID_DATE_RANGE' }.to_json)
      expect { client.transactions.list_for_account('a_1') }.to raise_error(TesoteSdk::InvalidDateRangeError)
    end
  end

  describe '#each_page_for_account (cursor)' do
    it 'walks pages until has_more false' do
      page1 = { total: 2, transactions: [{ id: 't_1' }],
                pagination: { has_more: true, after_id: 't_1' } }.to_json
      page2 = { total: 2, transactions: [{ id: 't_2' }],
                pagination: { has_more: false, after_id: 't_2' } }.to_json

      stub_request(:get, %r{/v2/accounts/a_1/transactions(?:\?|\z)})
        .to_return({ status: 200, body: page1 }, { status: 200, body: page2 })

      pages = client.transactions.each_page_for_account('a_1').to_a
      expect(pages.size).to eq(2)
    end
  end

  describe '#export' do
    it 'returns RawResponse with CSV body' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/transactions/export")
        .with(query: hash_including('format' => 'csv'))
        .to_return(status: 200, body: "id,date\n1,2026-04-01\n",
                   headers: { 'Content-Type' => 'text/csv',
                              'Content-Disposition' => 'attachment; filename=tx_a_1_now.csv',
                              'X-Request-Id' => 'req_export' })
      raw = client.transactions.export('a_1', { format: 'csv' })
      expect(raw).to be_a(TesoteSdk::Transport::RawResponse)
      expect(raw.body).to include('id,date')
      expect(raw.content_type).to eq('text/csv')
      expect(raw.content_disposition).to include('attachment')
    end

    it 'maps UNPROCESSABLE_CONTENT for export' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/transactions/export")
        .to_return(status: 422, body: { error_code: 'UNPROCESSABLE_CONTENT' }.to_json)
      expect { client.transactions.export('a_1') }.to raise_error(TesoteSdk::UnprocessableContentError)
    end
  end

  describe '#sync' do
    it 'POSTs to /accounts/{id}/transactions/sync and returns SyncResult' do
      payload = { added: [{ transaction_id: 't_1', account_id: 'a_1', amount: 1.0, date: '2026-04-01',
                            name: 'lunch', pending: false }],
                  modified: [], removed: [], next_cursor: 'c1', has_more: false }
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transactions/sync")
        .with(headers: { 'Content-Type' => 'application/json' })
        .to_return(status: 200, body: payload.to_json)
      result = client.transactions.sync('a_1', count: 100)
      expect(result).to be_a(TesoteSdk::Models::SyncResult)
      expect(result.added.first).to be_a(TesoteSdk::Models::SyncTransaction)
      expect(result.next_cursor).to eq('c1')
    end

    it 'maps INVALID_COUNT' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transactions/sync")
        .to_return(status: 422, body: { error_code: 'INVALID_COUNT' }.to_json)
      expect { client.transactions.sync('a_1') }.to raise_error(TesoteSdk::InvalidCountError)
    end

    it 'maps INVALID_CURSOR' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transactions/sync")
        .to_return(status: 422, body: { error_code: 'INVALID_CURSOR' }.to_json)
      expect { client.transactions.sync('a_1', cursor: 'x') }.to raise_error(TesoteSdk::InvalidCursorError)
    end

    it 'maps HISTORY_SYNC_FORBIDDEN' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transactions/sync")
        .to_return(status: 403, body: { error_code: 'HISTORY_SYNC_FORBIDDEN' }.to_json)
      expect { client.transactions.sync('a_1') }.to raise_error(TesoteSdk::HistorySyncForbiddenError)
    end
  end

  describe '#sync_legacy' do
    it 'POSTs to /transactions/sync (legacy non-nested)' do
      stub_request(:post, "#{base_url}/v2/transactions/sync")
        .to_return(status: 200, body: { added: [], modified: [], removed: [], next_cursor: nil,
                                        has_more: false }.to_json)
      expect(client.transactions.sync_legacy).to be_a(TesoteSdk::Models::SyncResult)
    end
  end

  describe '#get' do
    it 'returns Transaction' do
      stub_request(:get, "#{base_url}/v2/transactions/t_1")
        .to_return(status: 200, body: { id: 't_1', status: 'posted', data: {},
                                        transaction_categories: [] }.to_json)
      expect(client.transactions.get('t_1')).to be_a(TesoteSdk::Models::Transaction)
    end

    it 'maps TRANSACTION_NOT_FOUND' do
      stub_request(:get, "#{base_url}/v2/transactions/t_x")
        .to_return(status: 404, body: { error_code: 'TRANSACTION_NOT_FOUND' }.to_json)
      expect { client.transactions.get('t_x') }.to raise_error(TesoteSdk::TransactionNotFoundError)
    end
  end

  describe '#bulk' do
    it 'POSTs to /transactions/bulk and parses BulkResult' do
      payload = { bulk_results: [{ account_id: 'a_1', transactions: [],
                                   pagination: { has_more: false } }] }
      stub_request(:post, "#{base_url}/v2/transactions/bulk")
        .with(body: hash_including('account_ids' => ['a_1']))
        .to_return(status: 200, body: payload.to_json)
      result = client.transactions.bulk(account_ids: ['a_1'])
      expect(result).to be_a(TesoteSdk::Models::BulkResult)
      expect(result.bulk_results.first.account_id).to eq('a_1')
    end

    it 'raises ArgumentError on empty account_ids' do
      expect { client.transactions.bulk(account_ids: []) }.to raise_error(ArgumentError)
    end

    it 'maps UNPROCESSABLE_CONTENT' do
      stub_request(:post, "#{base_url}/v2/transactions/bulk")
        .to_return(status: 422, body: { error_code: 'UNPROCESSABLE_CONTENT' }.to_json)
      expect { client.transactions.bulk(account_ids: ['a_1']) }
        .to raise_error(TesoteSdk::UnprocessableContentError)
    end
  end

  describe '#search' do
    it 'returns SearchResult' do
      stub_request(:get, "#{base_url}/v2/transactions/search")
        .with(query: hash_including('q' => 'cafe'))
        .to_return(status: 200, body: { transactions: [], total: 0 }.to_json)
      result = client.transactions.search(q: 'cafe')
      expect(result).to be_a(TesoteSdk::Models::SearchResult)
    end

    it 'raises ArgumentError when q is missing' do
      expect { client.transactions.search }.to raise_error(ArgumentError)
    end
  end
end
