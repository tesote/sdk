import type { Transaction, TransactionListResponse } from '../models/transaction.js';
import type { Transport } from '../transport.js';

export type { Transaction, TransactionListResponse };

export interface TransactionListParams {
  start_date?: string;
  end_date?: string;
  scope?: string;
  page?: number;
  per_page?: number;
  transactions_after_id?: string;
  transactions_before_id?: string;
}

export class V1TransactionsClient {
  constructor(private readonly transport: Transport) {}

  async listForAccount(
    accountId: string,
    params: TransactionListParams = {},
  ): Promise<TransactionListResponse> {
    const res = await this.transport.request<TransactionListResponse>({
      method: 'GET',
      path: `/v1/accounts/${encodeURIComponent(accountId)}/transactions`,
      query: { ...params },
    });
    return res.data;
  }

  /**
   * Async iterator that follows cursor pagination across every page of
   * /v1/accounts/{id}/transactions. Yields transactions one at a time.
   */
  async *listAllForAccount(
    accountId: string,
    params: TransactionListParams = {},
  ): AsyncGenerator<Transaction, void, void> {
    let after = params.transactions_after_id;
    while (true) {
      const page = await this.listForAccount(accountId, {
        ...params,
        ...(after !== undefined ? { transactions_after_id: after } : {}),
      });
      for (const tx of page.transactions) yield tx;
      if (!page.pagination.has_more || page.pagination.after_id === null) return;
      after = page.pagination.after_id;
    }
  }

  async get(id: string): Promise<Transaction> {
    const res = await this.transport.request<Transaction>({
      method: 'GET',
      path: `/v1/transactions/${encodeURIComponent(id)}`,
      cache: { ttl: 300 },
    });
    return res.data;
  }
}
