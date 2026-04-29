module TesoteSdk
  # Typed PORO models for API responses. Wire format is preserved on the right —
  # attribute names match the JSON snake_case so callers can round-trip raw hashes
  # via .from_hash without surprise.
  #
  # All model structs are keyword-init and tolerate unknown keys (forward-compat
  # with future fields the SDK has not learned about yet).
  module Models # rubocop:disable Metrics/ModuleLength
    module FromHash
      # why: API may grow new fields between SDK releases — drop unknown keys
      # silently rather than crash. Clients can still get raw via response_body.
      def from_hash(hash)
        return nil if hash.nil?
        return hash if hash.is_a?(self)

        hash = hash.transform_keys(&:to_s) if hash.is_a?(Hash)
        known = members.map(&:to_s)
        attrs = {}
        known.each do |key|
          attrs[key.to_sym] = build_field(key, hash[key])
        end
        new(**attrs)
      end

      def from_array(arr)
        return [] if arr.nil?

        arr.map { |h| from_hash(h) }
      end

      # Override in subclasses that wrap nested objects.
      def build_field(_key, value)
        value
      end
    end

    Bank = Struct.new(:name, keyword_init: true) do
      extend FromHash
    end

    LegalEntity = Struct.new(:id, :legal_name, keyword_init: true) do
      extend FromHash
    end

    AccountData = Struct.new(
      :masked_account_number,
      :currency,
      :transactions_data_current_as_of,
      :balance_data_current_as_of,
      :custom_user_provided_identifier,
      :balance_cents,
      :available_balance_cents,
      keyword_init: true
    ) do
      extend FromHash
    end

    Account = Struct.new(
      :id,
      :name,
      :data,
      :bank,
      :legal_entity,
      :tesote_created_at,
      :tesote_updated_at,
      keyword_init: true
    ) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'data' then AccountData.from_hash(value)
        when 'bank' then Bank.from_hash(value)
        when 'legal_entity' then LegalEntity.from_hash(value)
        else value
        end
      end
    end

    TransactionCategory = Struct.new(
      :name,
      :external_category_code,
      :created_at,
      :updated_at,
      keyword_init: true
    ) do
      extend FromHash
    end

    Counterparty = Struct.new(:id, :name, keyword_init: true) do
      extend FromHash
    end

    TransactionData = Struct.new(
      :amount_cents,
      :currency,
      :description,
      :transaction_date,
      :created_at,
      :created_at_date,
      :note,
      :external_service_id,
      :running_balance_cents,
      keyword_init: true
    ) do
      extend FromHash
    end

    Transaction = Struct.new(
      :id,
      :status,
      :data,
      :tesote_imported_at,
      :tesote_updated_at,
      :transaction_categories,
      :counterparty,
      keyword_init: true
    ) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'data' then TransactionData.from_hash(value)
        when 'transaction_categories' then TransactionCategory.from_array(value)
        when 'counterparty' then Counterparty.from_hash(value)
        else value
        end
      end
    end

    SyncTransaction = Struct.new(
      :transaction_id,
      :account_id,
      :amount,
      :iso_currency_code,
      :unofficial_currency_code,
      :date,
      :datetime,
      :name,
      :merchant_name,
      :pending,
      :category,
      :running_balance_cents,
      keyword_init: true
    ) do
      extend FromHash
    end

    SyncRemoval = Struct.new(:transaction_id, :account_id, keyword_init: true) do
      extend FromHash
    end

    SyncResult = Struct.new(
      :added,
      :modified,
      :removed,
      :next_cursor,
      :has_more,
      keyword_init: true
    ) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'added', 'modified' then SyncTransaction.from_array(value)
        when 'removed' then SyncRemoval.from_array(value)
        else value
        end
      end
    end

    SyncStartResult = Struct.new(
      :message,
      :sync_session_id,
      :status,
      :started_at,
      keyword_init: true
    ) do
      extend FromHash
    end

    SyncSessionError = Struct.new(:type, :message, keyword_init: true) do
      extend FromHash
    end

    SyncSessionPerformance = Struct.new(
      :total_duration,
      :complexity_score,
      :sync_speed_score,
      keyword_init: true
    ) do
      extend FromHash
    end

    SyncSession = Struct.new(
      :id,
      :status,
      :started_at,
      :completed_at,
      :transactions_synced,
      :accounts_count,
      :error,
      :performance,
      keyword_init: true
    ) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'error' then SyncSessionError.from_hash(value)
        when 'performance' then SyncSessionPerformance.from_hash(value)
        else value
        end
      end
    end

    SourceAccount = Struct.new(:id, :name, :payment_method_id, keyword_init: true) do
      extend FromHash
    end

    Destination = Struct.new(
      :payment_method_id,
      :counterparty_id,
      :counterparty_name,
      keyword_init: true
    ) do
      extend FromHash
    end

    Fee = Struct.new(:amount, :currency, keyword_init: true) do
      extend FromHash
    end

    TesoteTransactionRef = Struct.new(:id, :status, keyword_init: true) do
      extend FromHash
    end

    LatestAttempt = Struct.new(
      :id,
      :status,
      :attempt_number,
      :external_reference,
      :submitted_at,
      :completed_at,
      :error_code,
      :error_message,
      keyword_init: true
    ) do
      extend FromHash
    end

    TransactionOrder = Struct.new(
      :id,
      :status,
      :amount,
      :currency,
      :description,
      :reference,
      :external_reference,
      :idempotency_key,
      :batch_id,
      :scheduled_for,
      :approved_at,
      :submitted_at,
      :completed_at,
      :failed_at,
      :cancelled_at,
      :source_account,
      :destination,
      :fee,
      :execution_strategy,
      :tesote_transaction,
      :latest_attempt,
      :created_at,
      :updated_at,
      keyword_init: true
    ) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'source_account' then SourceAccount.from_hash(value)
        when 'destination' then Destination.from_hash(value)
        when 'fee' then Fee.from_hash(value)
        when 'tesote_transaction' then TesoteTransactionRef.from_hash(value)
        when 'latest_attempt' then LatestAttempt.from_hash(value)
        else value
        end
      end
    end

    TesoteAccountRef = Struct.new(:id, :name, keyword_init: true) do
      extend FromHash
    end

    PaymentMethod = Struct.new(
      :id,
      :method_type,
      :currency,
      :label,
      :details,
      :verified,
      :verified_at,
      :last_used_at,
      :counterparty,
      :tesote_account,
      :created_at,
      :updated_at,
      keyword_init: true
    ) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'counterparty' then Counterparty.from_hash(value)
        when 'tesote_account' then TesoteAccountRef.from_hash(value)
        else value
        end
      end
    end

    BatchSummary = Struct.new(
      :batch_id,
      :total_orders,
      :total_amount_cents,
      :amount_currency,
      :statuses,
      :batch_status,
      :created_at,
      :orders,
      keyword_init: true
    ) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'orders' then TransactionOrder.from_array(value)
        else value
        end
      end
    end

    BatchCreateResult = Struct.new(:batch_id, :orders, :errors, keyword_init: true) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'orders' then TransactionOrder.from_array(value)
        else value
        end
      end
    end

    Pagination = Struct.new(
      :current_page,
      :per_page,
      :total_pages,
      :total_count,
      :has_more,
      :after_id,
      :before_id,
      keyword_init: true
    ) do
      extend FromHash
    end

    AccountList = Struct.new(:accounts, :total, :pagination, keyword_init: true) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'accounts' then Account.from_array(value)
        when 'pagination' then Pagination.from_hash(value)
        else value
        end
      end
    end

    TransactionList = Struct.new(:transactions, :total, :pagination, keyword_init: true) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'transactions' then Transaction.from_array(value)
        when 'pagination' then Pagination.from_hash(value)
        else value
        end
      end
    end

    OffsetPage = Struct.new(:items, :limit, :offset, :has_more, keyword_init: true) do
      extend FromHash
    end

    SearchResult = Struct.new(:transactions, :total, keyword_init: true) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'transactions' then Transaction.from_array(value)
        else value
        end
      end
    end

    BulkAccountResult = Struct.new(:account_id, :transactions, :pagination, keyword_init: true) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'transactions' then Transaction.from_array(value)
        when 'pagination' then Pagination.from_hash(value)
        else value
        end
      end
    end

    BulkResult = Struct.new(:bulk_results, keyword_init: true) do
      extend FromHash

      def self.build_field(key, value)
        case key
        when 'bulk_results' then BulkAccountResult.from_array(value)
        else value
        end
      end
    end

    Whoami = Struct.new(:client, keyword_init: true) do
      extend FromHash
    end

    StatusResult = Struct.new(:status, :authenticated, keyword_init: true) do
      extend FromHash
    end

    BatchApproveResult = Struct.new(:approved, :failed, keyword_init: true) do
      extend FromHash
    end

    BatchSubmitResult = Struct.new(:enqueued, :failed, keyword_init: true) do
      extend FromHash
    end

    BatchCancelResult = Struct.new(:cancelled, :skipped, :errors, keyword_init: true) do
      extend FromHash
    end
  end
end
