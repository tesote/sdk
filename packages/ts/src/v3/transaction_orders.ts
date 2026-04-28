import type { Transport } from '../transport.js';

export class V3TransactionOrdersClient {
  constructor(private readonly _transport: Transport) {}

  async list(_accountId: string, _params: Record<string, unknown> = {}): Promise<unknown> {
    throw new Error('not implemented');
  }

  async get(_accountId: string, _orderId: string): Promise<unknown> {
    throw new Error('not implemented');
  }

  async create(
    _accountId: string,
    _body: Record<string, unknown>,
    _opts: { idempotencyKey?: string } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }

  async submit(
    _accountId: string,
    _orderId: string,
    _opts: { idempotencyKey?: string } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }

  async cancel(
    _accountId: string,
    _orderId: string,
    _opts: { idempotencyKey?: string } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }
}
