module TesoteSdk
  module V3
    class Accounts
      def initialize(transport)
        @transport = transport
      end

      def list(query = {}, opts: {})
        @transport.request('GET', 'accounts', query: query, opts: opts)
      end

      def get(id, opts: {})
        raise ArgumentError, 'id is required' if id.nil? || id.to_s.empty?

        @transport.request('GET', "accounts/#{id}", opts: opts)
      end

      def sync(_id, opts: {})
        raise NotImplementedError, 'V3::Accounts#sync not wired in 0.1.0'
      end
    end
  end
end
