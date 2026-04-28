module TesoteSdk
  module V3
    # Resource clients documented in resources.md but not wired in 0.1.0.
    # Each method raises NotImplementedError so the surface is discoverable.

    class StubBase
      def initialize(transport)
        @transport = transport
      end
    end

    class Transactions < StubBase
      def list_for_account(_account_id, _query = {}, opts: {})
        raise NotImplementedError, 'V3::Transactions#list_for_account not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V3::Transactions#get not wired in 0.1.0'
      end

      def export(_account_id, _params = {}, opts: {})
        raise NotImplementedError, 'V3::Transactions#export not wired in 0.1.0'
      end

      def sync(_account_id, opts: {})
        raise NotImplementedError, 'V3::Transactions#sync not wired in 0.1.0'
      end

      def bulk(_payload, opts: {})
        raise NotImplementedError, 'V3::Transactions#bulk not wired in 0.1.0'
      end

      def search(_query = {}, opts: {})
        raise NotImplementedError, 'V3::Transactions#search not wired in 0.1.0'
      end
    end

    class SyncSessions < StubBase
      def list(_account_id, _query = {}, opts: {})
        raise NotImplementedError, 'V3::SyncSessions#list not wired in 0.1.0'
      end

      def get(_account_id, _id, opts: {})
        raise NotImplementedError, 'V3::SyncSessions#get not wired in 0.1.0'
      end
    end

    class TransactionOrders < StubBase
      def list(_account_id, _query = {}, opts: {})
        raise NotImplementedError, 'V3::TransactionOrders#list not wired in 0.1.0'
      end

      def get(_account_id, _id, opts: {})
        raise NotImplementedError, 'V3::TransactionOrders#get not wired in 0.1.0'
      end

      def create(_account_id, _payload, opts: {})
        raise NotImplementedError, 'V3::TransactionOrders#create not wired in 0.1.0'
      end

      def submit(_account_id, _id, opts: {})
        raise NotImplementedError, 'V3::TransactionOrders#submit not wired in 0.1.0'
      end

      def cancel(_account_id, _id, opts: {})
        raise NotImplementedError, 'V3::TransactionOrders#cancel not wired in 0.1.0'
      end
    end

    class Batches < StubBase
      def create(_payload, opts: {})
        raise NotImplementedError, 'V3::Batches#create not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V3::Batches#get not wired in 0.1.0'
      end

      def approve(_id, opts: {})
        raise NotImplementedError, 'V3::Batches#approve not wired in 0.1.0'
      end

      def submit(_id, opts: {})
        raise NotImplementedError, 'V3::Batches#submit not wired in 0.1.0'
      end

      def cancel(_id, opts: {})
        raise NotImplementedError, 'V3::Batches#cancel not wired in 0.1.0'
      end
    end

    class PaymentMethods < StubBase
      def list(_query = {}, opts: {})
        raise NotImplementedError, 'V3::PaymentMethods#list not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V3::PaymentMethods#get not wired in 0.1.0'
      end

      def create(_payload, opts: {})
        raise NotImplementedError, 'V3::PaymentMethods#create not wired in 0.1.0'
      end

      def update(_id, _payload, opts: {})
        raise NotImplementedError, 'V3::PaymentMethods#update not wired in 0.1.0'
      end

      def delete(_id, opts: {})
        raise NotImplementedError, 'V3::PaymentMethods#delete not wired in 0.1.0'
      end
    end

    class Categories < StubBase
      def list(_query = {}, opts: {})
        raise NotImplementedError, 'V3::Categories#list not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V3::Categories#get not wired in 0.1.0'
      end

      def create(_payload, opts: {})
        raise NotImplementedError, 'V3::Categories#create not wired in 0.1.0'
      end

      def update(_id, _payload, opts: {})
        raise NotImplementedError, 'V3::Categories#update not wired in 0.1.0'
      end

      def delete(_id, opts: {})
        raise NotImplementedError, 'V3::Categories#delete not wired in 0.1.0'
      end
    end

    class Counterparties < StubBase
      def list(_query = {}, opts: {})
        raise NotImplementedError, 'V3::Counterparties#list not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V3::Counterparties#get not wired in 0.1.0'
      end

      def create(_payload, opts: {})
        raise NotImplementedError, 'V3::Counterparties#create not wired in 0.1.0'
      end

      def update(_id, _payload, opts: {})
        raise NotImplementedError, 'V3::Counterparties#update not wired in 0.1.0'
      end

      def delete(_id, opts: {})
        raise NotImplementedError, 'V3::Counterparties#delete not wired in 0.1.0'
      end
    end

    class LegalEntities < StubBase
      def list(_query = {}, opts: {})
        raise NotImplementedError, 'V3::LegalEntities#list not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V3::LegalEntities#get not wired in 0.1.0'
      end
    end

    class Connections < StubBase
      def list(_query = {}, opts: {})
        raise NotImplementedError, 'V3::Connections#list not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V3::Connections#get not wired in 0.1.0'
      end

      def status(_id, opts: {})
        raise NotImplementedError, 'V3::Connections#status not wired in 0.1.0'
      end
    end

    class Webhooks < StubBase
      def list(_query = {}, opts: {})
        raise NotImplementedError, 'V3::Webhooks#list not wired in 0.1.0'
      end

      def get(_id, opts: {})
        raise NotImplementedError, 'V3::Webhooks#get not wired in 0.1.0'
      end

      def create(_payload, opts: {})
        raise NotImplementedError, 'V3::Webhooks#create not wired in 0.1.0'
      end

      def update(_id, _payload, opts: {})
        raise NotImplementedError, 'V3::Webhooks#update not wired in 0.1.0'
      end

      def delete(_id, opts: {})
        raise NotImplementedError, 'V3::Webhooks#delete not wired in 0.1.0'
      end
    end

    class Reports < StubBase
      def cash_flow(_query = {}, opts: {})
        raise NotImplementedError, 'V3::Reports#cash_flow not wired in 0.1.0'
      end
    end

    class BalanceHistory < StubBase
      def list_for_account(_account_id, _query = {}, opts: {})
        raise NotImplementedError, 'V3::BalanceHistory#list_for_account not wired in 0.1.0'
      end
    end

    class Workspace < StubBase
      def get(opts: {})
        raise NotImplementedError, 'V3::Workspace#get not wired in 0.1.0'
      end
    end

    class Mcp < StubBase
      def handle(_payload, opts: {})
        raise NotImplementedError, 'V3::Mcp#handle not wired in 0.1.0'
      end
    end

    class Status < StubBase
      def status(opts: {})
        raise NotImplementedError, 'V3::Status#status not wired in 0.1.0'
      end

      def whoami(opts: {})
        raise NotImplementedError, 'V3::Status#whoami not wired in 0.1.0'
      end
    end
  end
end
