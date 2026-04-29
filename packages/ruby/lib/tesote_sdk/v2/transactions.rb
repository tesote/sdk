module TesoteSdk
  module V2
    # Wraps:
    # - GET    /v2/accounts/{id}/transactions
    # - GET    /v2/accounts/{id}/transactions/export
    # - POST   /v2/accounts/{id}/transactions/sync
    # - POST   /v2/transactions/sync (legacy)
    # - GET    /v2/transactions/{id}
    # - POST   /v2/transactions/bulk
    # - GET    /v2/transactions/search
    class Transactions
      INDEX_CACHE_TTL = 60
      SHOW_CACHE_TTL = 300

      def initialize(transport)
        @transport = transport
      end

      # GET /v2/accounts/{id}/transactions
      def list_for_account(account_id, query = {}, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?

        merged = opts.key?(:cache) ? opts : opts.merge(cache: { ttl: INDEX_CACHE_TTL })
        body = @transport.request('GET', "accounts/#{account_id}/transactions", query: query, opts: merged)
        Models::TransactionList.from_hash(body)
      end

      # Cursor pagination over list_for_account; yields TransactionList pages.
      def each_page_for_account(account_id, query = {}, opts: {}, &block)
        return enum_for(:each_page_for_account, account_id, query, opts: opts) unless block

        enum = Pagination::CursorEnumerator.new(start_query: query) do |q|
          list_for_account(account_id, q, opts: opts)
        end
        enum.each(&block)
      end

      # GET /v2/accounts/{id}/transactions/export
      # Returns a Transport::RawResponse (file body, Content-Type, filename).
      def export(account_id, query = {}, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?

        @transport.request_raw('GET', "accounts/#{account_id}/transactions/export", query: query, opts: opts)
      end

      # POST /v2/accounts/{id}/transactions/sync
      # body: { count:, cursor:, options: } — all optional per spec.
      def sync(account_id, count: nil, cursor: nil, options: nil, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?

        payload = build_sync_payload(count: count, cursor: cursor, options: options)
        body = @transport.request('POST', "accounts/#{account_id}/transactions/sync", body: payload, opts: opts)
        Models::SyncResult.from_hash(body)
      end

      # POST /v2/transactions/sync (legacy, non-nested)
      def sync_legacy(count: nil, cursor: nil, options: nil, opts: {})
        payload = build_sync_payload(count: count, cursor: cursor, options: options)
        body = @transport.request('POST', 'transactions/sync', body: payload, opts: opts)
        Models::SyncResult.from_hash(body)
      end

      # GET /v2/transactions/{id}
      def get(id, opts: {})
        raise ArgumentError, 'id is required' if id.nil? || id.to_s.empty?

        merged = opts.key?(:cache) ? opts : opts.merge(cache: { ttl: SHOW_CACHE_TTL })
        body = @transport.request('GET', "transactions/#{id}", opts: merged)
        Models::Transaction.from_hash(body)
      end

      # POST /v2/transactions/bulk
      def bulk(account_ids:, page: nil, per_page: nil, limit: nil, offset: nil, opts: {})
        raise ArgumentError, 'account_ids is required and must not be empty' if account_ids.nil? || account_ids.empty?

        payload = {
          account_ids: account_ids,
          page: page,
          per_page: per_page,
          limit: limit,
          offset: offset
        }.compact
        body = @transport.request('POST', 'transactions/bulk', body: payload, opts: opts)
        Models::BulkResult.from_hash(body)
      end

      # GET /v2/transactions/search
      # Pass `q:` (required) plus any optional filters as keyword args
      # (account_id, limit, offset, status, type, start_date, etc.).
      # `opts:` is the transport options hash.
      def search(opts: {}, **filters)
        filters = filters.transform_keys(&:to_sym)
        q_value = filters[:q]
        raise ArgumentError, 'q is required' if q_value.nil? || q_value.to_s.empty?

        body = @transport.request('GET', 'transactions/search', query: filters.compact, opts: opts)
        Models::SearchResult.from_hash(body)
      end

      private

      def build_sync_payload(count:, cursor:, options:)
        payload = {}
        payload[:count] = count unless count.nil?
        payload[:cursor] = cursor unless cursor.nil?
        payload[:options] = options unless options.nil?
        payload
      end
    end
  end
end
