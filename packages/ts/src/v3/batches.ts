import type { Transport } from '../transport.js';

export class V3BatchesClient {
  constructor(private readonly _transport: Transport) {}

  async create(
    _body: Record<string, unknown>,
    _opts: { idempotencyKey?: string } = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }

  async get(_id: string): Promise<unknown> {
    throw new Error('not implemented');
  }

  async approve(_id: string, _opts: { idempotencyKey?: string } = {}): Promise<unknown> {
    throw new Error('not implemented');
  }

  async submit(_id: string, _opts: { idempotencyKey?: string } = {}): Promise<unknown> {
    throw new Error('not implemented');
  }

  async cancel(_id: string, _opts: { idempotencyKey?: string } = {}): Promise<unknown> {
    throw new Error('not implemented');
  }
}
