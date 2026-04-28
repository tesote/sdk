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

  module V3
    autoload :Client,             'tesote_sdk/v3/client'
    autoload :Accounts,           'tesote_sdk/v3/accounts'
    autoload :Transactions,       'tesote_sdk/v3/stubs'
    autoload :SyncSessions,       'tesote_sdk/v3/stubs'
    autoload :TransactionOrders,  'tesote_sdk/v3/stubs'
    autoload :Batches,            'tesote_sdk/v3/stubs'
    autoload :PaymentMethods,     'tesote_sdk/v3/stubs'
    autoload :Categories,         'tesote_sdk/v3/stubs'
    autoload :Counterparties,     'tesote_sdk/v3/stubs'
    autoload :LegalEntities,      'tesote_sdk/v3/stubs'
    autoload :Connections,        'tesote_sdk/v3/stubs'
    autoload :Webhooks,           'tesote_sdk/v3/stubs'
    autoload :Reports,            'tesote_sdk/v3/stubs'
    autoload :BalanceHistory,     'tesote_sdk/v3/stubs'
    autoload :Workspace,          'tesote_sdk/v3/stubs'
    autoload :Mcp,                'tesote_sdk/v3/stubs'
    autoload :Status,             'tesote_sdk/v3/stubs'
    autoload :WebhooksSignature,  'tesote_sdk/v3/webhooks_signature'
  end
end

# Eagerly load the v3 webhook helper so TesoteSdk::V3.verify_webhook_signature
# is callable without first instantiating a client.
require_relative 'tesote_sdk/v3/webhooks_signature'
