require_relative 'tesote_sdk/version'
require_relative 'tesote_sdk/errors'
require_relative 'tesote_sdk/transport'

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
    autoload :Transactions,       'tesote_sdk/v2/stubs'
    autoload :SyncSessions,       'tesote_sdk/v2/stubs'
    autoload :TransactionOrders,  'tesote_sdk/v2/stubs'
    autoload :Batches,            'tesote_sdk/v2/stubs'
    autoload :PaymentMethods,     'tesote_sdk/v2/stubs'
    autoload :Status,             'tesote_sdk/v2/stubs'
  end
end
