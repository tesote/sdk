module TesoteSdk
  module V1
    class Transactions
      def initialize(transport)
        @transport = transport
      end

      def list_for_account(_account_id, _query = {}, opts: {})
        raise NotImplementedError, 'V1::Transactions#list_for_account not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V1::Transactions#get not wired in 0.1.0'
      end
    end
  end
end
