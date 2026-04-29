import type { BulkResponse } from '../models/bulk.js';
import type { SearchResult } from '../models/search.js';
import type { SyncResult } from '../models/sync_transaction.js';
import type {
  CursorPagination,
  Transaction,
  TransactionListResponse,
} from '../models/transaction.js';
import type { Transport } from '../transport.js';

export interface V2TransactionListParams {
  start_date?: string;
  end_date?: string;
  scope?: string;
  page?: number;
  per_page?: number;
  transactions_after_id?: string;
  transactions_before_id?: string;
  transaction_date_after?: string;
  transaction_date_before?: string;
  created_after?: string;
  updated_after?: string;
  amount_min?: number;
  amount_max?: number;
  amount?: number;
  status?: string;
  category_id?: string;
  counterparty_id?: string;
  q?: string;
  type?: string;
  reference_code?: string;
}

export interface V2TransactionExportParams extends V2TransactionListParams {
  format?: 'csv' | 'json';
}

export interface ExportResponse {
  /** CSV body or pretty-printed JSON body, depending on `format`. */
  body: string;
  contentType: string | null;
  /** Suggested filename parsed from Content-Disposition (best-effort). */
  filename: string | null;
}

export interface SyncRequestOptions {
  count?: number;
  cursor?: string | 'now' | null;
  options?: {
    include_running_balance?: boolean;
  };
}

export interface BulkRequest {
  account_ids: string[];
  page?: number;
  per_page?: number;
  limit?: number;
  offset?: number;
}

export interface SearchParams extends V2TransactionListParams {
  q: string;
  account_id?: string;
  limit?: number;
  offset?: number;
}

const DISPOSITION_FILENAME = /filename\*?=(?:UTF-8'')?"?([^";]+)"?/i;

function parseFilename(disposition: string | null): string | null {
  if (disposition === null) return null;
  const match = DISPOSITION_FILENAME.exec(disposition);
  if (match === null) return null;
  const raw = match[1];
  if (raw === undefined) return null;
  try {
    return decodeURIComponent(raw);
  } catch {
    return raw;
  }
}

export class V2TransactionsClient {
  constructor(private readonly transport: Transport) {}

  /** GET /v2/accounts/{id}/transactions */
  async listForAccount(
    accountId: string,
    params: V2TransactionListParams = {},
  ): Promise<TransactionListResponse> {
    const res = await this.transport.request<TransactionListResponse>({
      method: 'GET',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/transactions`,
      query: { ...params },
      cache: { ttl: 60 },
    });
    return res.data;
  }

  /** Cursor-following async iterator over /v2/accounts/{id}/transactions. */
  async *listAllForAccount(
    accountId: string,
    params: V2TransactionListParams = {},
  ): AsyncGenerator<Transaction, void, void> {
    let after = params.transactions_after_id;
    while (true) {
      const page = await this.listForAccount(accountId, {
        ...params,
        ...(after !== undefined ? { transactions_after_id: after } : {}),
      });
      for (const tx of page.transactions) yield tx;
      const pg: CursorPagination = page.pagination;
      if (!pg.has_more || pg.after_id === null) return;
      after = pg.after_id;
    }
  }

  /** GET /v2/transactions/{id} */
  async get(id: string): Promise<Transaction> {
    const res = await this.transport.request<Transaction>({
      method: 'GET',
      path: `/v2/transactions/${encodeURIComponent(id)}`,
      cache: { ttl: 300 },
    });
    return res.data;
  }

  /**
   * GET /v2/accounts/{id}/transactions/export — CSV or JSON file download.
   * The body is returned as a string; up to 10,000 transactions per call.
   */
  async export(accountId: string, params: V2TransactionExportParams = {}): Promise<ExportResponse> {
    const res = await this.transport.request<string>({
      method: 'GET',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/transactions/export`,
      query: { ...params },
      headers: { Accept: '*/*' },
    });
    const contentType = res.headers.get('content-type');
    const body = typeof res.data === 'string' ? res.data : JSON.stringify(res.data);
    return {
      body,
      contentType,
      filename: parseFilename(res.headers.get('content-disposition')),
    };
  }

  /**
   * POST /v2/accounts/{id}/transactions/sync — Plaid-style flattened sync.
   */
  async sync(
    accountId: string,
    body: SyncRequestOptions = {},
    opts: { idempotencyKey?: string } = {},
  ): Promise<SyncResult> {
    const res = await this.transport.request<SyncResult>({
      method: 'POST',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/transactions/sync`,
      body: { ...body },
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }

  /**
   * POST /v2/transactions/sync (legacy non-nested route). Identical request/response
   * to {@link sync} but takes account context in the body rather than the path.
   */
  async syncLegacy(
    body: SyncRequestOptions & { account_id?: string } = {},
    opts: { idempotencyKey?: string } = {},
  ): Promise<SyncResult> {
    const res = await this.transport.request<SyncResult>({
      method: 'POST',
      path: '/v2/transactions/sync',
      body: { ...body },
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }

  /** POST /v2/transactions/bulk */
  async bulk(body: BulkRequest, opts: { idempotencyKey?: string } = {}): Promise<BulkResponse> {
    const res = await this.transport.request<BulkResponse>({
      method: 'POST',
      path: '/v2/transactions/bulk',
      body,
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }

  /** GET /v2/transactions/search */
  async search(params: SearchParams): Promise<SearchResult> {
    const res = await this.transport.request<SearchResult>({
      method: 'GET',
      path: '/v2/transactions/search',
      query: { ...params },
    });
    return res.data;
  }
}
