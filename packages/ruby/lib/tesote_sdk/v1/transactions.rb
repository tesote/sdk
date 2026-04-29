module TesoteSdk
  module V1
    class Transactions
      SHOW_CACHE_TTL = 300

      def initialize(transport)
        @transport = transport
      end

      # GET /v1/accounts/{id}/transactions — convenience pass-through to
      # V1::Accounts#list_transactions for callers that prefer to start from
      # the transactions client.
      def list_for_account(account_id, query = {}, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?

        body = @transport.request('GET', "accounts/#{account_id}/transactions", query: query, opts: opts)
        Models::TransactionList.from_hash(body)
      end

      # GET /v1/transactions/{id}
      def get(id, opts: {})
        raise ArgumentError, 'id is required' if id.nil? || id.to_s.empty?

        merged = opts.key?(:cache) ? opts : opts.merge(cache: { ttl: SHOW_CACHE_TTL })
        body = @transport.request('GET', "transactions/#{id}", opts: merged)
        Models::Transaction.from_hash(body)
      end
    end
  end
end
