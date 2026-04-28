import type { Transport } from '../transport.js';

export class V3ReportsClient {
  constructor(private readonly _transport: Transport) {}

  async cashFlow(_params: Record<string, unknown> = {}): Promise<unknown> {
    throw new Error('not implemented');
  }
}
