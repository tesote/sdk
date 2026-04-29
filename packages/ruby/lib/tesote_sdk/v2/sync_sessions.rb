module TesoteSdk
  module V2
    class SyncSessions
      def initialize(transport)
        @transport = transport
      end

      # GET /v2/accounts/{id}/sync_sessions
      def list(account_id, query = {}, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?

        body = @transport.request('GET', "accounts/#{account_id}/sync_sessions", query: query, opts: opts)
        items_raw = (body && body['sync_sessions']) || []
        Models::OffsetPage.new(
          items: Models::SyncSession.from_array(items_raw),
          limit: body && body['limit'],
          offset: body && body['offset'],
          has_more: body && body['has_more']
        )
      end

      # Walks pages via offset pagination; yields OffsetPage per call.
      def each_page(account_id, query = {}, opts: {}, page_size: 50, &block)
        return enum_for(:each_page, account_id, query, opts: opts, page_size: page_size) unless block

        enum = Pagination::OffsetEnumerator.new(start_query: query, limit: page_size) do |q|
          list(account_id, q, opts: opts)
        end
        enum.each(&block)
      end

      # GET /v2/accounts/{id}/sync_sessions/{session_id}
      def get(account_id, session_id, opts: {})
        raise ArgumentError, 'account_id is required' if account_id.nil? || account_id.to_s.empty?
        raise ArgumentError, 'session_id is required' if session_id.nil? || session_id.to_s.empty?

        body = @transport.request('GET', "accounts/#{account_id}/sync_sessions/#{session_id}", opts: opts)
        Models::SyncSession.from_hash(body)
      end
    end
  end
end
