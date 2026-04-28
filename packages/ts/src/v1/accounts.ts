import type { Transport } from '../transport.js';

export interface Account {
  id: string;
  name: string;
  currency: string;
  balance?: string | number;
  [k: string]: unknown;
}

export interface AccountListParams {
  cursor?: string;
  limit?: number;
}

export interface AccountListResponse {
  data: Account[];
  next_cursor?: string | null;
  [k: string]: unknown;
}

export class V1AccountsClient {
  constructor(private readonly transport: Transport) {}

  async list(params: AccountListParams = {}): Promise<AccountListResponse> {
    const res = await this.transport.request<AccountListResponse>({
      method: 'GET',
      path: '/v1/accounts',
      query: { ...params },
    });
    return res.data;
  }

  async get(id: string): Promise<Account> {
    const res = await this.transport.request<Account>({
      method: 'GET',
      path: `/v1/accounts/${encodeURIComponent(id)}`,
    });
    return res.data;
  }
}
