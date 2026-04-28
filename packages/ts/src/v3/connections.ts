import type { Transport } from '../transport.js';

export class V3ConnectionsClient {
  constructor(private readonly _transport: Transport) {}

  async list(_params: Record<string, unknown> = {}): Promise<unknown> {
    throw new Error('not implemented');
  }

  async get(_id: string): Promise<unknown> {
    throw new Error('not implemented');
  }

  async status(_id: string): Promise<unknown> {
    throw new Error('not implemented');
  }
}
