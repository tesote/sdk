module TesoteSdk
  module V2
    class Batches
      def initialize(transport)
        @transport = transport
      end

      # POST /v2/accounts/{id}/batches
      # `orders` is an array of order hashes (per spec).
      def create(account_id, orders:, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?
        raise ArgumentError, 'orders is required' if orders.nil? || orders.empty?

        payload = { orders: orders }
        body = @transport.request('POST', "accounts/#{account_id}/batches", body: payload, opts: opts)
        Models::BatchCreateResult.from_hash(body)
      end

      # GET /v2/accounts/{id}/batches/{batch_id}
      def get(account_id, batch_id, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?
        raise ArgumentError, 'batch_id is required' if batch_id.nil? || batch_id.to_s.empty?

        body = @transport.request('GET', "accounts/#{account_id}/batches/#{batch_id}", opts: opts)
        Models::BatchSummary.from_hash(body)
      end

      # POST /v2/accounts/{id}/batches/{batch_id}/approve
      def approve(account_id, batch_id, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?
        raise ArgumentError, 'batch_id is required' if batch_id.nil? || batch_id.to_s.empty?

        body = @transport.request(
          'POST',
          "accounts/#{account_id}/batches/#{batch_id}/approve",
          body: {},
          opts: opts
        )
        Models::BatchApproveResult.from_hash(body)
      end

      # POST /v2/accounts/{id}/batches/{batch_id}/submit
      def submit(account_id, batch_id, token: nil, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?
        raise ArgumentError, 'batch_id is required' if batch_id.nil? || batch_id.to_s.empty?

        payload = token.nil? ? {} : { token: token }
        body = @transport.request(
          'POST',
          "accounts/#{account_id}/batches/#{batch_id}/submit",
          body: payload,
          opts: opts
        )
        Models::BatchSubmitResult.from_hash(body)
      end

      # POST /v2/accounts/{id}/batches/{batch_id}/cancel
      def cancel(account_id, batch_id, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?
        raise ArgumentError, 'batch_id is required' if batch_id.nil? || batch_id.to_s.empty?

        body = @transport.request(
          'POST',
          "accounts/#{account_id}/batches/#{batch_id}/cancel",
          body: {},
          opts: opts
        )
        Models::BatchCancelResult.from_hash(body)
      end
    end
  end
end
