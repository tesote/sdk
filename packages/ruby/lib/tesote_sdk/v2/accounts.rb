module TesoteSdk
  module V2
    class Accounts
      INDEX_CACHE_TTL = 60
      SHOW_CACHE_TTL = 300

      def initialize(transport)
        @transport = transport
      end

      # GET /v2/accounts
      def list(query = {}, opts: {})
        body = @transport.request('GET', 'accounts', query: query, opts: with_cache(opts, INDEX_CACHE_TTL))
        Models::AccountList.from_hash(body)
      end

      # GET /v2/accounts/{id}
      def get(id, opts: {})
        raise ArgumentError, 'id is required' if id.nil? || id.to_s.empty?

        body = @transport.request('GET', "accounts/#{id}", opts: with_cache(opts, SHOW_CACHE_TTL))
        Models::Account.from_hash(body)
      end

      # POST /v2/accounts/{id}/sync
      # Triggers an async sync; returns SyncStartResult (status: pending).
      def sync(id, opts: {})
        raise ArgumentError, 'id is required' if id.nil? || id.to_s.empty?

        body = @transport.request('POST', "accounts/#{id}/sync", body: {}, opts: opts)
        Models::SyncStartResult.from_hash(body)
      end

      private

      def with_cache(opts, ttl)
        return opts if opts.key?(:cache)

        opts.merge(cache: { ttl: ttl })
      end
    end
  end
end
