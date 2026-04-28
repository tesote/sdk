module TesoteSdk
  module V1
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
    end
  end
end
