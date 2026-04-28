import type { Transport } from '../transport.js';

export class V2PaymentMethodsClient {
  constructor(private readonly _transport: Transport) {}

  async list(_params: Record<string, unknown> = {}): Promise<unknown> {
    throw new Error('not implemented');
  }

  async get(_id: string): Promise<unknown> {
    throw new Error('not implemented');
  }

  async create(
    _body: Record<string, unknown>,
    _opts: { idempotencyKey?: string } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }

  async update(
    _id: string,
    _body: Record<string, unknown>,
    _opts: { idempotencyKey?: string } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }

  async delete(_id: string, _opts: { idempotencyKey?: string } = {}): Promise<unknown> {
    throw new Error('not implemented');
  }
}
