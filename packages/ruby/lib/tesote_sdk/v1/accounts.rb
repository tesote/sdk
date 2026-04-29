module TesoteSdk
  module V1
    class Accounts
      INDEX_CACHE_TTL = 60      # 1 minute ETag (spec)
      SHOW_CACHE_TTL = 300      # 5 minutes ETag (spec)

      def initialize(transport)
        @transport = transport
      end

      # GET /v1/accounts
      def list(query = {}, opts: {})
        body = @transport.request('GET', 'accounts', query: query, opts: with_cache(opts, INDEX_CACHE_TTL))
        Models::AccountList.from_hash(body)
      end

      # GET /v1/accounts/{id}
      def get(id, opts: {})
        raise ArgumentError, 'id is required' if id.nil? || id.to_s.empty?

        body = @transport.request('GET', "accounts/#{id}", opts: with_cache(opts, SHOW_CACHE_TTL))
        Models::Account.from_hash(body)
      end

      # GET /v1/accounts/{id}/transactions
      def list_transactions(account_id, query = {}, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?

        body = @transport.request('GET', "accounts/#{account_id}/transactions", query: query, opts: opts)
        Models::TransactionList.from_hash(body)
      end

      # Enumerate all transactions across pages (cursor-based).
      # Yields TransactionList page objects.
      def each_transaction_page(account_id, query = {}, opts: {}, &block)
        return enum_for(:each_transaction_page, account_id, query, opts: opts) unless block

        enum = Pagination::CursorEnumerator.new(start_query: query) do |q|
          list_transactions(account_id, q, opts: opts)
        end
        enum.each(&block)
      end

      private

      def with_cache(opts, ttl)
        return opts if opts.key?(:cache)

        opts.merge(cache: { ttl: ttl })
      end
    end
  end
end
