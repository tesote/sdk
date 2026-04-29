import { Transport, type TransportOptions } from '../transport.js';
import { V2AccountsClient } from './accounts.js';
import { V2BatchesClient } from './batches.js';
import { V2PaymentMethodsClient } from './payment_methods.js';
import { V2StatusClient } from './status.js';
import { V2SyncSessionsClient } from './sync_sessions.js';
import { V2TransactionOrdersClient } from './transaction_orders.js';
import { V2TransactionsClient } from './transactions.js';

export type V2ClientOptions = TransportOptions;

export class V2Client {
  public readonly transport: Transport;
  public readonly accounts: V2AccountsClient;
  public readonly transactions: V2TransactionsClient;
  public readonly syncSessions: V2SyncSessionsClient;
  public readonly transactionOrders: V2TransactionOrdersClient;
  public readonly batches: V2BatchesClient;
  public readonly paymentMethods: V2PaymentMethodsClient;
  public readonly status: V2StatusClient;

  constructor(options: V2ClientOptions) {
    this.transport = new Transport(options);
    this.accounts = new V2AccountsClient(this.transport);
    this.transactions = new V2TransactionsClient(this.transport);
    this.syncSessions = new V2SyncSessionsClient(this.transport);
    this.transactionOrders = new V2TransactionOrdersClient(this.transport);
    this.batches = new V2BatchesClient(this.transport);
    this.paymentMethods = new V2PaymentMethodsClient(this.transport);
    this.status = new V2StatusClient(this.transport);
  }

  get lastRateLimit() {
    return this.transport.lastRateLimit;
  }
}

export { V2AccountsClient } from './accounts.js';
export { V2TransactionsClient } from './transactions.js';
export { V2SyncSessionsClient } from './sync_sessions.js';
export { V2TransactionOrdersClient } from './transaction_orders.js';
export { V2BatchesClient } from './batches.js';
export { V2PaymentMethodsClient } from './payment_methods.js';
export { V2StatusClient } from './status.js';

export type {
  V2TransactionListParams,
  V2TransactionExportParams,
  ExportResponse,
  SyncRequestOptions,
  BulkRequest,
  SearchParams,
} from './transactions.js';
export type { SyncSessionListParams } from './sync_sessions.js';
export type {
  TransactionOrderListParams,
  TransactionOrderCreateInput,
  TransactionOrderSubmitOptions,
} from './transaction_orders.js';
export type { BatchOrderInput, BatchCreateInput } from './batches.js';
export type {
  PaymentMethodListParams,
  PaymentMethodCreateInput,
  PaymentMethodUpdateInput,
} from './payment_methods.js';
