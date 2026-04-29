module TesoteSdk
  module V2
    class TransactionOrders
      def initialize(transport)
        @transport = transport
      end

      # GET /v2/accounts/{id}/transaction_orders
      # Returns OffsetPage with TransactionOrder items.
      def list(account_id, query = {}, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?

        body = @transport.request('GET', "accounts/#{account_id}/transaction_orders", query: query, opts: opts)
        items_raw = (body && body['items']) || []
        Models::OffsetPage.new(
          items: Models::TransactionOrder.from_array(items_raw),
          limit: body && body['limit'],
          offset: body && body['offset'],
          has_more: body && body['has_more']
        )
      end

      # GET /v2/accounts/{id}/transaction_orders/{order_id}
      def get(account_id, order_id, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?
        raise ArgumentError, 'order_id is required' if order_id.nil? || order_id.to_s.empty?

        body = @transport.request('GET', "accounts/#{account_id}/transaction_orders/#{order_id}", opts: opts)
        Models::TransactionOrder.from_hash(body)
      end

      # POST /v2/accounts/{id}/transaction_orders
      # `order` is a hash matching the spec's `transaction_order` body field.
      def create(account_id, order:, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?
        raise ArgumentError, 'order is required' if order.nil?

        payload = { transaction_order: order }
        body = @transport.request('POST', "accounts/#{account_id}/transaction_orders", body: payload, opts: opts)
        Models::TransactionOrder.from_hash(body)
      end

      # POST /v2/accounts/{id}/transaction_orders/{order_id}/submit
      def submit(account_id, order_id, token: nil, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?
        raise ArgumentError, 'order_id is required' if order_id.nil? || order_id.to_s.empty?

        payload = token.nil? ? {} : { token: token }
        body = @transport.request(
          'POST',
          "accounts/#{account_id}/transaction_orders/#{order_id}/submit",
          body: payload,
          opts: opts
        )
        Models::TransactionOrder.from_hash(body)
      end

      # POST /v2/accounts/{id}/transaction_orders/{order_id}/cancel
      def cancel(account_id, order_id, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?
        raise ArgumentError, 'order_id is required' if order_id.nil? || order_id.to_s.empty?

        body = @transport.request(
          'POST',
          "accounts/#{account_id}/transaction_orders/#{order_id}/cancel",
          body: {},
          opts: opts
        )
        Models::TransactionOrder.from_hash(body)
      end
    end
  end
end
