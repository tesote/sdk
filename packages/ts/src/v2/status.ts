import type { Transport } from '../transport.js';

export class V2StatusClient {
  constructor(private readonly _transport: Transport) {}

  async status(): Promise<unknown> {
    throw new Error('not implemented');
  }

  async whoami(): Promise<unknown> {
    throw new Error('not implemented');
  }
}
