require_relative '../transport'
require_relative 'accounts'
require_relative 'stubs'

module TesoteSdk
  module V2
    class Client
      VERSION_SEGMENT = 'v2'.freeze

      attr_reader :transport

      def initialize(api_key:, **transport_options)
        @transport = Transport.new(
          api_key: api_key,
          version_segment: VERSION_SEGMENT,
          **transport_options
        )
      end

      def accounts
        @accounts ||= Accounts.new(transport)
      end

      def transactions
        @transactions ||= Transactions.new(transport)
      end

      def sync_sessions
        @sync_sessions ||= SyncSessions.new(transport)
      end

      def transaction_orders
        @transaction_orders ||= TransactionOrders.new(transport)
      end

      def batches
        @batches ||= Batches.new(transport)
      end

      def payment_methods
        @payment_methods ||= PaymentMethods.new(transport)
      end

      def status
        @status ||= Status.new(transport)
      end

      def last_rate_limit
        transport.last_rate_limit
      end

      def last_request_id
        transport.last_request_id
      end
    end
  end
end
