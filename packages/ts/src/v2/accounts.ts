import type { Transport } from '../transport.js';
import type { Account, AccountListParams, AccountListResponse } from '../v1/accounts.js';

export class V2AccountsClient {
  constructor(private readonly transport: Transport) {}

  async list(params: AccountListParams = {}): Promise<AccountListResponse> {
    const res = await this.transport.request<AccountListResponse>({
      method: 'GET',
      path: '/v2/accounts',
      query: { ...params },
    });
    return res.data;
  }

  async get(id: string): Promise<Account> {
    const res = await this.transport.request<Account>({
      method: 'GET',
      path: `/v2/accounts/${encodeURIComponent(id)}`,
    });
    return res.data;
  }

  async sync(_id: string, _opts: { idempotencyKey?: string } = {}): Promise<unknown> {
    throw new Error('not implemented');
  }
}
