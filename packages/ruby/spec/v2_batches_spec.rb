require 'spec_helper'

RSpec.describe TesoteSdk::V2::Batches do
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) { TesoteSdk::V2::Client.new(api_key: 'k', base_url: base_url, base_delay: 0.0, max_delay: 0.0) }

  describe '#create' do
    it 'POSTs orders and returns BatchCreateResult with idempotency-key' do
      header_seen = nil
      stub_request(:post, "#{base_url}/v2/accounts/a_1/batches")
        .with do |req|
          header_seen = req.headers['Idempotency-Key']
          req.headers['Content-Type'] == 'application/json'
        end
        .to_return(status: 201,
                   body: { batch_id: 'b_1',
                           orders: [{ id: 'to_1', status: 'draft', batch_id: 'b_1' }],
                           errors: [] }.to_json)
      orders = [{ amount: '1.00', currency: 'VES', beneficiary: { name: 'X' } }]
      result = client.batches.create('a_1', orders: orders)
      expect(result).to be_a(TesoteSdk::Models::BatchCreateResult)
      expect(result.batch_id).to eq('b_1')
      expect(result.orders.first).to be_a(TesoteSdk::Models::TransactionOrder)
      expect(header_seen).to match(/\A[0-9a-f-]{36}\z/i)
    end

    it 'maps BATCH_VALIDATION_ERROR' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/batches")
        .to_return(status: 400, body: { error_code: 'BATCH_VALIDATION_ERROR' }.to_json)
      expect { client.batches.create('a_1', orders: [{ amount: '1' }]) }
        .to raise_error(TesoteSdk::BatchValidationError)
    end

    it 'raises ArgumentError on empty orders' do
      expect { client.batches.create('a_1', orders: []) }.to raise_error(ArgumentError)
    end
  end

  describe '#get' do
    it 'returns BatchSummary' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/batches/b_1")
        .to_return(status: 200, body: { batch_id: 'b_1', total_orders: 2,
                                        total_amount_cents: 200,
                                        amount_currency: 'VES',
                                        statuses: { 'draft' => 2 },
                                        batch_status: 'draft',
                                        created_at: 'now',
                                        orders: [{ id: 'to_1', status: 'draft' }] }.to_json)
      summary = client.batches.get('a_1', 'b_1')
      expect(summary).to be_a(TesoteSdk::Models::BatchSummary)
      expect(summary.orders.first).to be_a(TesoteSdk::Models::TransactionOrder)
    end

    it 'maps BATCH_NOT_FOUND' do
      stub_request(:get, "#{base_url}/v2/accounts/a_1/batches/missing")
        .to_return(status: 404, body: { error_code: 'BATCH_NOT_FOUND' }.to_json)
      expect { client.batches.get('a_1', 'missing') }.to raise_error(TesoteSdk::BatchNotFoundError)
    end
  end

  describe '#approve / #submit / #cancel' do
    it 'approve POSTs and returns BatchApproveResult' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/batches/b_1/approve")
        .to_return(status: 200, body: { approved: 5, failed: 0 }.to_json)
      result = client.batches.approve('a_1', 'b_1')
      expect(result).to be_a(TesoteSdk::Models::BatchApproveResult)
      expect(result.approved).to eq(5)
    end

    it 'submit POSTs token and returns BatchSubmitResult' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/batches/b_1/submit")
        .with(body: { token: 'mfa' }.to_json)
        .to_return(status: 200, body: { enqueued: 5, failed: 0 }.to_json)
      result = client.batches.submit('a_1', 'b_1', token: 'mfa')
      expect(result.enqueued).to eq(5)
    end

    it 'cancel POSTs and returns BatchCancelResult' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/batches/b_1/cancel")
        .to_return(status: 200, body: { cancelled: 3, skipped: 2, errors: [] }.to_json)
      result = client.batches.cancel('a_1', 'b_1')
      expect(result.cancelled).to eq(3)
    end

    it 'maps INVALID_ORDER_STATE on approve' do
      stub_request(:post, "#{base_url}/v2/accounts/a_1/batches/b_1/approve")
        .to_return(status: 409, body: { error_code: 'INVALID_ORDER_STATE' }.to_json)
      expect { client.batches.approve('a_1', 'b_1') }.to raise_error(TesoteSdk::InvalidOrderStateError)
    end
  end
end
