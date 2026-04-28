module TesoteSdk
  module V2
    # Resource clients whose endpoints are documented but not wired in 0.1.0.
    # Each method raises NotImplementedError so the public surface is discoverable
    # while back-compat constraints stay honest.

    class StubBase
      def initialize(transport)
        @transport = transport
      end
    end

    class Transactions < StubBase
      def list_for_account(_account_id, _query = {}, opts: {})
        raise NotImplementedError, 'V2::Transactions#list_for_account not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V2::Transactions#get not wired in 0.1.0'
      end

      def export(_account_id, _params = {}, opts: {})
        raise NotImplementedError, 'V2::Transactions#export not wired in 0.1.0'
      end

      def sync(_account_id, opts: {})
        raise NotImplementedError, 'V2::Transactions#sync not wired in 0.1.0'
      end

      def bulk(_payload, opts: {})
        raise NotImplementedError, 'V2::Transactions#bulk not wired in 0.1.0'
      end

      def search(_query = {}, opts: {})
        raise NotImplementedError, 'V2::Transactions#search not wired in 0.1.0'
      end
    end

    class SyncSessions < StubBase
      def list(_account_id, _query = {}, opts: {})
        raise NotImplementedError, 'V2::SyncSessions#list not wired in 0.1.0'
      end

      def get(_account_id, _id, opts: {})
        raise NotImplementedError, 'V2::SyncSessions#get not wired in 0.1.0'
      end
    end

    class TransactionOrders < StubBase
      def list(_account_id, _query = {}, opts: {})
        raise NotImplementedError, 'V2::TransactionOrders#list not wired in 0.1.0'
      end

      def get(_account_id, _id, opts: {})
        raise NotImplementedError, 'V2::TransactionOrders#get not wired in 0.1.0'
      end

      def create(_account_id, _payload, opts: {})
        raise NotImplementedError, 'V2::TransactionOrders#create not wired in 0.1.0'
      end

      def submit(_account_id, _id, opts: {})
        raise NotImplementedError, 'V2::TransactionOrders#submit not wired in 0.1.0'
      end

      def cancel(_account_id, _id, opts: {})
        raise NotImplementedError, 'V2::TransactionOrders#cancel not wired in 0.1.0'
      end
    end

    class Batches < StubBase
      def create(_payload, opts: {})
        raise NotImplementedError, 'V2::Batches#create not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V2::Batches#get not wired in 0.1.0'
      end

      def approve(_id, opts: {})
        raise NotImplementedError, 'V2::Batches#approve not wired in 0.1.0'
      end

      def submit(_id, opts: {})
        raise NotImplementedError, 'V2::Batches#submit not wired in 0.1.0'
      end

      def cancel(_id, opts: {})
        raise NotImplementedError, 'V2::Batches#cancel not wired in 0.1.0'
      end
    end

    class PaymentMethods < StubBase
      def list(_query = {}, opts: {})
        raise NotImplementedError, 'V2::PaymentMethods#list not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V2::PaymentMethods#get not wired in 0.1.0'
      end

      def create(_payload, opts: {})
        raise NotImplementedError, 'V2::PaymentMethods#create not wired in 0.1.0'
      end

      def update(_id, _payload, opts: {})
        raise NotImplementedError, 'V2::PaymentMethods#update not wired in 0.1.0'
      end

      def delete(_id, opts: {})
        raise NotImplementedError, 'V2::PaymentMethods#delete not wired in 0.1.0'
      end
    end

    class Status < StubBase
      def status(opts: {})
        raise NotImplementedError, 'V2::Status#status not wired in 0.1.0'
      end

      def whoami(opts: {})
        raise NotImplementedError, 'V2::Status#whoami not wired in 0.1.0'
      end
    end
  end
end
