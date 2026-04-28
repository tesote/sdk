import type { Transport } from '../transport.js';

export class V3WorkspaceClient {
  constructor(private readonly _transport: Transport) {}

  async get(): Promise<unknown> {
    throw new Error('not implemented');
  }
}
