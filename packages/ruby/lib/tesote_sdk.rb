require_relative 'tesote_sdk/version'
require_relative 'tesote_sdk/errors'
require_relative 'tesote_sdk/transport'
require_relative 'tesote_sdk/models'
require_relative 'tesote_sdk/pagination'

module TesoteSdk
  module V1
    autoload :Client,       'tesote_sdk/v1/client'
    autoload :Accounts,     'tesote_sdk/v1/accounts'
    autoload :Transactions, 'tesote_sdk/v1/transactions'
    autoload :Status,       'tesote_sdk/v1/status'
  end

  module V2
    autoload :Client,             'tesote_sdk/v2/client'
    autoload :Accounts,           'tesote_sdk/v2/accounts'
    autoload :Transactions,       'tesote_sdk/v2/transactions'
    autoload :SyncSessions,       'tesote_sdk/v2/sync_sessions'
    autoload :TransactionOrders,  'tesote_sdk/v2/transaction_orders'
    autoload :Batches,            'tesote_sdk/v2/batches'
    autoload :PaymentMethods,     'tesote_sdk/v2/payment_methods'
    autoload :Status,             'tesote_sdk/v2/status'
  end
end
