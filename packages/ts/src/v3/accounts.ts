import type { CacheOptions, Transport } from '../transport.js';
import type { Account, AccountListParams, AccountListResponse } from '../v1/accounts.js';

export interface V3AccountListParams extends AccountListParams {
  cache?: CacheOptions | false;
}

export class V3AccountsClient {
  constructor(private readonly transport: Transport) {}

  async list(params: V3AccountListParams = {}): Promise<AccountListResponse> {
    const { cache, ...query } = params;
    const res = await this.transport.request<AccountListResponse>({
      method: 'GET',
      path: '/v3/accounts',
      query: { ...query },
      ...(cache === undefined ? {} : { cache }),
    });
    return res.data;
  }

  async get(id: string, opts: { cache?: CacheOptions | false } = {}): Promise<Account> {
    const res = await this.transport.request<Account>({
      method: 'GET',
      path: `/v3/accounts/${encodeURIComponent(id)}`,
      ...(opts.cache === undefined ? {} : { cache: opts.cache }),
    });
    return res.data;
  }

  async sync(_id: string, _opts: { idempotencyKey?: string } = {}): Promise<unknown> {
    throw new Error('not implemented');
  }
}
