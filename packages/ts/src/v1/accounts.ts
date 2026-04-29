import type { Account, AccountListResponse, PageBasedPagination } from '../models/account.js';
import type { Transport } from '../transport.js';

export type { Account, AccountListResponse, PageBasedPagination };

export interface AccountListParams {
  page?: number;
  per_page?: number;
  include?: string;
  sort?: string;
}

export class V1AccountsClient {
  constructor(private readonly transport: Transport) {}

  async list(params: AccountListParams = {}): Promise<AccountListResponse> {
    const res = await this.transport.request<AccountListResponse>({
      method: 'GET',
      path: '/v1/accounts',
      query: { ...params },
      cache: { ttl: 60 },
    });
    return res.data;
  }

  /**
   * Async iterator over every page of /v1/accounts. Yields accounts one at a time.
   */
  async *listAll(params: AccountListParams = {}): AsyncGenerator<Account, void, void> {
    let page = params.page ?? 1;
    const perPage = params.per_page ?? 50;
    while (true) {
      const res = await this.list({ ...params, page, per_page: perPage });
      for (const account of res.accounts) yield account;
      if (page >= res.pagination.total_pages) return;
      page += 1;
    }
  }

  async get(id: string): Promise<Account> {
    const res = await this.transport.request<Account>({
      method: 'GET',
      path: `/v1/accounts/${encodeURIComponent(id)}`,
      cache: { ttl: 300 },
    });
    return res.data;
  }
}
