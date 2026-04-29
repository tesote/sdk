import { Transport, type TransportOptions } from '../transport.js';
import { V1AccountsClient } from './accounts.js';
import { V1StatusClient } from './status.js';
import { V1TransactionsClient } from './transactions.js';

export type V1ClientOptions = TransportOptions;

export class V1Client {
  public readonly transport: Transport;
  public readonly accounts: V1AccountsClient;
  public readonly transactions: V1TransactionsClient;
  public readonly status: V1StatusClient;

  constructor(options: V1ClientOptions) {
    this.transport = new Transport(options);
    this.accounts = new V1AccountsClient(this.transport);
    this.transactions = new V1TransactionsClient(this.transport);
    this.status = new V1StatusClient(this.transport);
  }

  get lastRateLimit() {
    return this.transport.lastRateLimit;
  }
}

export { V1AccountsClient } from './accounts.js';
export { V1TransactionsClient } from './transactions.js';
export { V1StatusClient } from './status.js';
export type {
  Account,
  AccountListParams,
  AccountListResponse,
  PageBasedPagination,
} from './accounts.js';
export type {
  Transaction,
  TransactionListParams,
  TransactionListResponse,
} from './transactions.js';
