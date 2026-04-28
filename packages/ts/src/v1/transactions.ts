import type { Transport } from '../transport.js';

export interface Transaction {
  id: string;
  account_id: string;
  amount: string | number;
  currency: string;
  posted_at?: string;
  [k: string]: unknown;
}

export class V1TransactionsClient {
  constructor(private readonly _transport: Transport) {}

  async listForAccount(
    _accountId: string,
    _params: { cursor?: string; limit?: number } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }

  async get(_id: string): Promise<Transaction> {
    throw new Error('not implemented');
  }
}
