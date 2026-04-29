require 'spec_helper'

RSpec.describe TesoteSdk::V2::TransactionOrders do
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) { TesoteSdk::V2::Client.new(api_key: 'k', base_url: base_url, base_delay: 0.0, max_delay: 0.0) }

  describe '#list' do
    it 'returns OffsetPage of TransactionOrder' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/transaction_orders")
        .to_return(status: 200, body: { items: [{ id: 'to_1', status: 'draft' }],
                                        has_more: false, limit: 50, offset: 0 }.to_json)
      page = client.transaction_orders.list('a_1')
      expect(page).to be_a(TesoteSdk::Models::OffsetPage)
      expect(page.items.first).to be_a(TesoteSdk::Models::TransactionOrder)
    end
  end

  describe '#get' do
    it 'returns TransactionOrder' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/transaction_orders/to_1")
        .to_return(status: 200, body: { id: 'to_1', status: 'draft',
                                        latest_attempt: { id: 'la_1', status: 'pending', attempt_number: 1 } }.to_json)
      result = client.transaction_orders.get('a_1', 'to_1')
      expect(result.latest_attempt).to be_a(TesoteSdk::Models::LatestAttempt)
    end

    it 'maps TRANSACTION_ORDER_NOT_FOUND' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/transaction_orders/missing")
        .to_return(status: 404, body: { error_code: 'TRANSACTION_ORDER_NOT_FOUND' }.to_json)
      expect { client.transaction_orders.get('a_1', 'missing') }
        .to raise_error(TesoteSdk::TransactionOrderNotFoundError)
    end
  end

  describe '#create' do
    it 'POSTs and returns TransactionOrder; sends Idempotency-Key header' do
      header_seen = nil
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transaction_orders")
        .with do |req|
          header_seen = req.headers['Idempotency-Key']
          true
        end
        .to_return(status: 201, body: { id: 'to_1', status: 'draft' }.to_json)
      order = { amount: '10.00', currency: 'VES', description: 'test', beneficiary: { name: 'X' } }
      result = client.transaction_orders.create('a_1', order: order)
      expect(result).to be_a(TesoteSdk::Models::TransactionOrder)
      expect(header_seen).to match(/\A[0-9a-f-]{36}\z/i)
    end

    it 'forwards a caller-supplied idempotency_key opt to the header' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transaction_orders")
        .with(headers: { 'Idempotency-Key' => 'caller-key' })
        .to_return(status: 201, body: { id: 'to_1', status: 'draft' }.to_json)
      client.transaction_orders.create('a_1',
                                       order: { amount: '1.00', currency: 'VES' },
                                       opts: { idempotency_key: 'caller-key' })
    end

    it 'maps VALIDATION_ERROR' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transaction_orders")
        .to_return(status: 400, body: { error_code: 'VALIDATION_ERROR', error: 'bad' }.to_json)
      expect { client.transaction_orders.create('a_1', order: {}) }
        .to raise_error(TesoteSdk::ValidationError)
    end

    it 'maps 415 to plain ApiError when Content-Type missing on server' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transaction_orders")
        .to_return(status: 415, body: { error: 'unsupported' }.to_json)
      expect { client.transaction_orders.create('a_1', order: {}) }
        .to raise_error(TesoteSdk::ApiError) { |err| expect(err.http_status).to eq(415) }
    end
  end

  describe '#submit' do
    it 'POSTs to submit endpoint with token' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transaction_orders/to_1/submit")
        .with(body: { token: 'mfa-1' }.to_json)
        .to_return(status: 202, body: { id: 'to_1', status: 'processing' }.to_json)
      result = client.transaction_orders.submit('a_1', 'to_1', token: 'mfa-1')
      expect(result.status).to eq('processing')
    end

    it 'maps INVALID_ORDER_STATE' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transaction_orders/to_1/submit")
        .to_return(status: 409, body: { error_code: 'INVALID_ORDER_STATE' }.to_json)
      expect { client.transaction_orders.submit('a_1', 'to_1') }
        .to raise_error(TesoteSdk::InvalidOrderStateError)
    end
  end

  describe '#cancel' do
    it 'POSTs to cancel endpoint' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/transaction_orders/to_1/cancel")
        .to_return(status: 200, body: { id: 'to_1', status: 'cancelled' }.to_json)
      result = client.transaction_orders.cancel('a_1', 'to_1')
      expect(result.status).to eq('cancelled')
    end
  end
end
