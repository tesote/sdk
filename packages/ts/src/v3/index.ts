import { Transport, type TransportOptions } from '../transport.js';
import { V3AccountsClient } from './accounts.js';
import { V3BalanceHistoryClient } from './balance_history.js';
import { V3BatchesClient } from './batches.js';
import { V3CategoriesClient } from './categories.js';
import { V3ConnectionsClient } from './connections.js';
import { V3CounterpartiesClient } from './counterparties.js';
import { V3LegalEntitiesClient } from './legal_entities.js';
import { V3McpClient } from './mcp.js';
import { V3PaymentMethodsClient } from './payment_methods.js';
import { V3ReportsClient } from './reports.js';
import { V3StatusClient } from './status.js';
import { V3SyncSessionsClient } from './sync_sessions.js';
import { V3TransactionOrdersClient } from './transaction_orders.js';
import { V3TransactionsClient } from './transactions.js';
import { V3WebhooksClient } from './webhooks.js';
import { V3WorkspaceClient } from './workspace.js';

export type V3ClientOptions = TransportOptions;

export class V3Client {
  public readonly transport: Transport;
  public readonly accounts: V3AccountsClient;
  public readonly transactions: V3TransactionsClient;
  public readonly syncSessions: V3SyncSessionsClient;
  public readonly transactionOrders: V3TransactionOrdersClient;
  public readonly batches: V3BatchesClient;
  public readonly paymentMethods: V3PaymentMethodsClient;
  public readonly balanceHistory: V3BalanceHistoryClient;
  public readonly categories: V3CategoriesClient;
  public readonly counterparties: V3CounterpartiesClient;
  public readonly legalEntities: V3LegalEntitiesClient;
  public readonly connections: V3ConnectionsClient;
  public readonly webhooks: V3WebhooksClient;
  public readonly reports: V3ReportsClient;
  public readonly workspace: V3WorkspaceClient;
  public readonly mcp: V3McpClient;
  public readonly status: V3StatusClient;

  constructor(options: V3ClientOptions) {
    this.transport = new Transport(options);
    this.accounts = new V3AccountsClient(this.transport);
    this.transactions = new V3TransactionsClient(this.transport);
    this.syncSessions = new V3SyncSessionsClient(this.transport);
    this.transactionOrders = new V3TransactionOrdersClient(this.transport);
    this.batches = new V3BatchesClient(this.transport);
    this.paymentMethods = new V3PaymentMethodsClient(this.transport);
    this.balanceHistory = new V3BalanceHistoryClient(this.transport);
    this.categories = new V3CategoriesClient(this.transport);
    this.counterparties = new V3CounterpartiesClient(this.transport);
    this.legalEntities = new V3LegalEntitiesClient(this.transport);
    this.connections = new V3ConnectionsClient(this.transport);
    this.webhooks = new V3WebhooksClient(this.transport);
    this.reports = new V3ReportsClient(this.transport);
    this.workspace = new V3WorkspaceClient(this.transport);
    this.mcp = new V3McpClient(this.transport);
    this.status = new V3StatusClient(this.transport);
  }

  get lastRateLimit() {
    return this.transport.lastRateLimit;
  }
}

export { V3AccountsClient } from './accounts.js';
export { V3TransactionsClient } from './transactions.js';
export { V3SyncSessionsClient } from './sync_sessions.js';
export { V3TransactionOrdersClient } from './transaction_orders.js';
export { V3BatchesClient } from './batches.js';
export { V3PaymentMethodsClient } from './payment_methods.js';
export { V3BalanceHistoryClient } from './balance_history.js';
export { V3CategoriesClient } from './categories.js';
export { V3CounterpartiesClient } from './counterparties.js';
export { V3LegalEntitiesClient } from './legal_entities.js';
export { V3ConnectionsClient } from './connections.js';
export { V3WebhooksClient, verifyWebhookSignature } from './webhooks.js';
export { V3ReportsClient } from './reports.js';
export { V3WorkspaceClient } from './workspace.js';
export { V3McpClient } from './mcp.js';
export { V3StatusClient } from './status.js';
