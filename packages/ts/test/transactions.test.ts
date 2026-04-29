import { describe, expect, it } from 'vitest';
import {
  AccountNotFoundError,
  HistorySyncForbiddenError,
  InvalidCountError,
  InvalidCursorError,
  InvalidDateRangeError,
  TransactionNotFoundError,
  UnprocessableContentError,
} from '../src/errors.js';
import { V1Client, V2Client } from '../src/index.js';
import {
  callAt,
  getBody,
  getHeader,
  getMethod,
  jsonResponse,
  makeFetchMock,
  rawResponse,
} from './helpers.js';

const txFixture = (id: string) => ({
  id,
  status: 'posted',
  data: {
    amount_cents: 1000,
    currency: 'VES',
    description: 'coffee',
    transaction_date: '2026-04-28',
    created_at: '2026-04-28T12:00:00Z',
    created_at_date: '2026-04-28',
    note: null,
    external_service_id: null,
  },
  tesote_imported_at: '2026-04-28T12:00:00Z',
  tesote_updated_at: '2026-04-28T12:00:00Z',
  transaction_categories: [],
  counterparty: null,
});

const cursorPage = (
  ids: string[],
  hasMore: boolean,
  afterId: string | null,
  beforeId: string | null,
) => ({
  total: ids.length,
  transactions: ids.map(txFixture),
  pagination: { has_more: hasMore, per_page: ids.length, after_id: afterId, before_id: beforeId },
});

describe('V1 transactions', () => {
  it('listForAccount serializes filters', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, cursorPage(['t1'], false, 't1', 't1')),
    ]);
    const c = new V1Client({ apiKey: 'k', fetch });
    await c.transactions.listForAccount('a1', { start_date: '2026-04-01', per_page: 10 });
    expect(calls[0]?.url).toContain('/v1/accounts/a1/transactions');
    expect(calls[0]?.url).toContain('start_date=2026-04-01');
    expect(calls[0]?.url).toContain('per_page=10');
  });

  it('listForAccount → 422 INVALID_DATE_RANGE', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(422, { error_code: 'INVALID_DATE_RANGE', error: 'bad' }),
    ]);
    const c = new V1Client({ apiKey: 'k', fetch });
    await expect(c.transactions.listForAccount('a1', { start_date: 'x' })).rejects.toBeInstanceOf(
      InvalidDateRangeError,
    );
  });

  it('listAllForAccount follows after_id cursor', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, cursorPage(['t1', 't2'], true, 't2', 't1')),
      jsonResponse(200, cursorPage(['t3'], false, 't3', 't3')),
    ]);
    const c = new V1Client({ apiKey: 'k', fetch });
    const ids: string[] = [];
    for await (const t of c.transactions.listAllForAccount('a1')) ids.push(t.id);
    expect(ids).toEqual(['t1', 't2', 't3']);
    expect(calls[1]?.url).toContain('transactions_after_id=t2');
  });

  it('get → 404 maps to TransactionNotFoundError', async () => {
    const { fetch } = makeFetchMock([jsonResponse(404, { error_code: 'TRANSACTION_NOT_FOUND' })]);
    const c = new V1Client({ apiKey: 'k', fetch });
    await expect(c.transactions.get('t1')).rejects.toBeInstanceOf(TransactionNotFoundError);
  });

  it('get returns Transaction', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, txFixture('t1'))]);
    const c = new V1Client({ apiKey: 'k', fetch });
    const t = await c.transactions.get('t1');
    expect(t.id).toBe('t1');
    expect(calls[0]?.url).toContain('/v1/transactions/t1');
  });
});

describe('V2 transactions', () => {
  it('listForAccount hits /v2 path with filters', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, cursorPage(['t1'], false, 't1', 't1')),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.transactions.listForAccount('a1', { q: 'coffee', amount_min: 100 });
    expect(calls[0]?.url).toContain('/v2/accounts/a1/transactions');
    expect(calls[0]?.url).toContain('q=coffee');
    expect(calls[0]?.url).toContain('amount_min=100');
  });

  it('listForAccount → 404 ACCOUNT_NOT_FOUND', async () => {
    const { fetch } = makeFetchMock([jsonResponse(404, { error_code: 'ACCOUNT_NOT_FOUND' })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.transactions.listForAccount('missing')).rejects.toBeInstanceOf(
      AccountNotFoundError,
    );
  });

  it('get returns v1-shape transaction', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, txFixture('t9'))]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const t = await c.transactions.get('t9');
    expect(t.id).toBe('t9');
    expect(calls[0]?.url).toContain('/v2/transactions/t9');
  });

  it('export returns raw CSV body and parses filename', async () => {
    const csv = 'Transaction ID,Date\nt1,2026-04-28\n';
    const { fetch, calls } = makeFetchMock([
      rawResponse(200, csv, {
        'content-type': 'text/csv',
        'content-disposition': 'attachment; filename="transactions_a1_2026-04-28.csv"',
      }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.transactions.export('a1', { format: 'csv' });
    expect(r.body).toBe(csv);
    expect(r.contentType).toBe('text/csv');
    expect(r.filename).toBe('transactions_a1_2026-04-28.csv');
    expect(calls[0]?.url).toContain('format=csv');
  });

  it('sync POSTs to nested path with idempotency-key', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, {
        added: [],
        modified: [],
        removed: [],
        next_cursor: null,
        has_more: false,
      }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.transactions.sync('a1', { count: 50, cursor: 'now' }, { idempotencyKey: 'i1' });
    expect(getMethod(callAt(calls, 0))).toBe('POST');
    expect(calls[0]?.url).toContain('/v2/accounts/a1/transactions/sync');
    expect(getBody(callAt(calls, 0))).toEqual({ count: 50, cursor: 'now' });
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toBe('i1');
    expect(getHeader(callAt(calls, 0), 'content-type')).toBe('application/json');
  });

  it('sync 422 INVALID_COUNT', async () => {
    const { fetch } = makeFetchMock([jsonResponse(422, { error_code: 'INVALID_COUNT' })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.transactions.sync('a1', { count: 99999 })).rejects.toBeInstanceOf(
      InvalidCountError,
    );
  });

  it('sync 422 INVALID_CURSOR', async () => {
    const { fetch } = makeFetchMock([jsonResponse(422, { error_code: 'INVALID_CURSOR' })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.transactions.sync('a1', { cursor: 'garbage' })).rejects.toBeInstanceOf(
      InvalidCursorError,
    );
  });

  it('sync 403 HISTORY_SYNC_FORBIDDEN', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(403, { error_code: 'HISTORY_SYNC_FORBIDDEN', error: 'enable feature' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.transactions.sync('a1', { cursor: null })).rejects.toBeInstanceOf(
      HistorySyncForbiddenError,
    );
  });

  it('syncLegacy hits /v2/transactions/sync', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, {
        added: [],
        modified: [],
        removed: [],
        next_cursor: null,
        has_more: false,
      }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.transactions.syncLegacy({ account_id: 'a1', count: 10 });
    expect(calls[0]?.url).toContain('/v2/transactions/sync');
    expect(getBody(callAt(calls, 0))).toEqual({ account_id: 'a1', count: 10 });
  });

  it('bulk POSTs account_ids body', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, { bulk_results: [] })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.transactions.bulk({ account_ids: ['a1', 'a2'], per_page: 25 });
    expect(getBody(callAt(calls, 0))).toEqual({ account_ids: ['a1', 'a2'], per_page: 25 });
    expect(calls[0]?.url).toContain('/v2/transactions/bulk');
  });

  it('bulk → 422 UNPROCESSABLE_CONTENT', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(422, { error_code: 'UNPROCESSABLE_CONTENT', error: 'too many' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.transactions.bulk({ account_ids: [] })).rejects.toBeInstanceOf(
      UnprocessableContentError,
    );
  });

  it('search query builds URL', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, { transactions: [], total: 0 })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.transactions.search({ q: 'starbucks', limit: 50 });
    expect(calls[0]?.url).toContain('/v2/transactions/search');
    expect(calls[0]?.url).toContain('q=starbucks');
    expect(calls[0]?.url).toContain('limit=50');
  });

  it('listAllForAccount drives cursor pagination', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(200, cursorPage(['t1', 't2'], true, 't2', 't1')),
      jsonResponse(200, cursorPage(['t3'], false, 't3', 't3')),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const ids: string[] = [];
    for await (const t of c.transactions.listAllForAccount('a1')) ids.push(t.id);
    expect(ids).toEqual(['t1', 't2', 't3']);
  });

  it('POST without server complains 415 — surfaced as UnprocessableContentError when error_code maps so', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(415, { error_code: 'UNPROCESSABLE_CONTENT', error: 'need json' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    // why: 415 means missing Content-Type. Real SDK always sends it on bodies, but the
    // server-side code is exercised here to confirm the mapping pass-through.
    await expect(c.transactions.bulk({ account_ids: ['a1'] })).rejects.toBeInstanceOf(
      UnprocessableContentError,
    );
  });
});
