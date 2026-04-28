import type { Transport } from '../transport.js';

export class V3BalanceHistoryClient {
  constructor(private readonly _transport: Transport) {}

  async listForAccount(
    _accountId: string,
    _params: Record<string, unknown> = {},
  ): Promise<unknown> {
    throw new Error('not implemented');
  }
}
