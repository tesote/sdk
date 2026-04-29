import type {
  SyncSession,
  SyncSessionListResponse,
  SyncSessionStatus,
} from '../models/sync_session.js';
import type { Transport } from '../transport.js';

export interface SyncSessionListParams {
  limit?: number;
  offset?: number;
  status?: SyncSessionStatus;
}

export class V2SyncSessionsClient {
  constructor(private readonly transport: Transport) {}

  /** GET /v2/accounts/{id}/sync_sessions */
  async list(
    accountId: string,
    params: SyncSessionListParams = {},
  ): Promise<SyncSessionListResponse> {
    const res = await this.transport.request<SyncSessionListResponse>({
      method: 'GET',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/sync_sessions`,
      query: { ...params },
    });
    return res.data;
  }

  /** Async iterator paging through every sync session for an account. */
  async *listAll(
    accountId: string,
    params: SyncSessionListParams = {},
  ): AsyncGenerator<SyncSession, void, void> {
    let offset = params.offset ?? 0;
    const limit = params.limit ?? 50;
    while (true) {
      const page = await this.list(accountId, { ...params, limit, offset });
      for (const s of page.sync_sessions) yield s;
      if (!page.has_more) return;
      offset += page.sync_sessions.length;
    }
  }

  /** GET /v2/accounts/{id}/sync_sessions/{session_id} */
  async get(accountId: string, sessionId: string): Promise<SyncSession> {
    const res = await this.transport.request<SyncSession>({
      method: 'GET',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/sync_sessions/${encodeURIComponent(sessionId)}`,
    });
    return res.data;
  }
}
