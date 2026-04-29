import type { Account, AccountListResponse } from '../models/account.js';
import type { SyncStartResponse } from '../models/sync_session.js';
import type { Transport } from '../transport.js';
import type { AccountListParams } from '../v1/accounts.js';

export class V2AccountsClient {
  constructor(private readonly transport: Transport) {}

  async list(params: AccountListParams = {}): Promise<AccountListResponse> {
    const res = await this.transport.request<AccountListResponse>({
      method: 'GET',
      path: '/v2/accounts',
      query: { ...params },
      cache: { ttl: 60 },
    });
    return res.data;
  }

  /**
   * Async iterator over every page of /v2/accounts. Yields accounts one at a time.
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
      path: `/v2/accounts/${encodeURIComponent(id)}`,
      cache: { ttl: 300 },
    });
    return res.data;
  }

  /**
   * POST /v2/accounts/{id}/sync — fire a bank sync. Returns a SyncStartResponse;
   * poll sync_sessions for completion.
   */
  async sync(id: string, opts: { idempotencyKey?: string } = {}): Promise<SyncStartResponse> {
    const res = await this.transport.request<SyncStartResponse>({
      method: 'POST',
      path: `/v2/accounts/${encodeURIComponent(id)}/sync`,
      body: {},
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }
}
