require 'spec_helper'

RSpec.describe TesoteSdk::V2::PaymentMethods do
  let(:base_url) { 'https://equipo.tesote.com/api' }
  let(:client) { TesoteSdk::V2::Client.new(api_key: 'k', base_url: base_url, base_delay: 0.0, max_delay: 0.0) }

  describe '#list' do
    it 'returns OffsetPage of PaymentMethod' do
      stub_request(:get, "#{base_url}/v2/payment_methods")
        .to_return(status: 200, body: { items: [{ id: 'pm_1', method_type: 'bank_account', currency: 'VES' }],
                                        has_more: false, limit: 50, offset: 0 }.to_json)
      page = client.payment_methods.list
      expect(page).to be_a(TesoteSdk::Models::OffsetPage)
      expect(page.items.first).to be_a(TesoteSdk::Models::PaymentMethod)
    end
  end

  describe '#get' do
    it 'returns PaymentMethod' do
      stub_request(:get, "#{base_url}/v2/payment_methods/pm_1")
        .to_return(status: 200, body: { id: 'pm_1', method_type: 'bank_account', currency: 'VES',
                                        counterparty: { id: 'cp_1', name: 'X' } }.to_json)
      pm = client.payment_methods.get('pm_1')
      expect(pm).to be_a(TesoteSdk::Models::PaymentMethod)
      expect(pm.counterparty).to be_a(TesoteSdk::Models::Counterparty)
    end

    it 'maps PAYMENT_METHOD_NOT_FOUND' do
      stub_request(:get, "#{base_url}/v2/payment_methods/missing")
        .to_return(status: 404, body: { error_code: 'PAYMENT_METHOD_NOT_FOUND' }.to_json)
      expect { client.payment_methods.get('missing') }.to raise_error(TesoteSdk::PaymentMethodNotFoundError)
    end
  end

  describe '#create' do
    it 'POSTs and returns PaymentMethod with Content-Type and Idempotency-Key' do
      stub_request(:post, "#{base_url}/v2/payment_methods")
        .with(headers: { 'Content-Type' => 'application/json' })
        .to_return(status: 201, body: { id: 'pm_1', method_type: 'bank_account', currency: 'VES' }.to_json)
      pm = client.payment_methods.create(payment_method: { method_type: 'bank_account', currency: 'VES' })
      expect(pm).to be_a(TesoteSdk::Models::PaymentMethod)
    end

    it 'maps VALIDATION_ERROR' do
      stub_request(:post, "#{base_url}/v2/payment_methods")
        .to_return(status: 400, body: { error_code: 'VALIDATION_ERROR' }.to_json)
      expect { client.payment_methods.create(payment_method: {}) }
        .to raise_error(TesoteSdk::ValidationError)
    end
  end

  describe '#update' do
    it 'PATCHes and returns PaymentMethod' do
      stub_request(:patch, "#{base_url}/v2/payment_methods/pm_1")
        .with(headers: { 'Content-Type' => 'application/json' },
              body: { payment_method: { label: 'New' } }.to_json)
        .to_return(status: 200, body: { id: 'pm_1', method_type: 'bank_account', currency: 'VES',
                                        label: 'New' }.to_json)
      result = client.payment_methods.update('pm_1', payment_method: { label: 'New' })
      expect(result.label).to eq('New')
    end
  end

  describe '#delete' do
    it 'DELETEs and returns nil on 204' do
      stub_request(:delete, "#{base_url}/v2/payment_methods/pm_1")
        .to_return(status: 204, body: '')
      expect(client.payment_methods.delete('pm_1')).to be_nil
    end

    it 'maps 409 VALIDATION_ERROR (in-use)' do
      stub_request(:delete, "#{base_url}/v2/payment_methods/pm_1")
        .to_return(status: 409, body: { error_code: 'VALIDATION_ERROR' }.to_json)
      expect { client.payment_methods.delete('pm_1') }.to raise_error(TesoteSdk::ValidationError)
    end
  end
end
