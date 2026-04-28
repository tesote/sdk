import type { Transport } from '../transport.js';

export class V3TransactionsClient {
  constructor(private readonly _transport: Transport) {}

  async listForAccount(
    _accountId: string,
    _params: Record<string, unknown> = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }

  async get(_id: string): Promise<unknown> {
    throw new Error('not implemented');
  }

  async export(_params: Record<string, unknown>): Promise<unknown> {
    throw new Error('not implemented');
  }

  async sync(_params: Record<string, unknown>): Promise<unknown> {
    throw new Error('not implemented');
  }

  async bulk(
    _items: ReadonlyArray<unknown>,
    _opts: { idempotencyKey?: string } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }

  async search(_params: Record<string, unknown>): Promise<unknown> {
    throw new Error('not implemented');
  }
}
