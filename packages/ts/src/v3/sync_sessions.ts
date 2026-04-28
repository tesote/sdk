import type { Transport } from '../transport.js';

export class V3SyncSessionsClient {
  constructor(private readonly _transport: Transport) {}

  async list(_accountId: string, _params: Record<string, unknown> = {}): Promise<unknown> {
    throw new Error('not implemented');
  }

  async get(_accountId: string, _sessionId: string): Promise<unknown> {
    throw new Error('not implemented');
  }
}
