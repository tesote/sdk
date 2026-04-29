import { describe, expect, it } from 'vitest';
import {
  AccountNotFoundError,
  BankConnectionNotFoundError,
  BankUnderMaintenanceError,
  SyncInProgressError,
  SyncRateLimitExceededError,
  UnauthorizedError,
} from '../src/errors.js';
import { V1Client, V2Client } from '../src/index.js';
import { callAt, getBody, getHeader, getMethod, jsonResponse, makeFetchMock } from './helpers.js';

const accountFixture = (id: string) => ({
  id,
  name: 'My Account',
  data: {
    masked_account_number: '****1234',
    currency: 'VES',
    transactions_data_current_as_of: null,
    balance_data_current_as_of: null,
    custom_user_provided_identifier: null,
  },
  bank: { name: 'Test Bank' },
  legal_entity: { id: null, legal_name: null },
  tesote_created_at: '2026-01-01T00:00:00Z',
  tesote_updated_at: '2026-01-01T00:00:00Z',
});

const pageEnvelope = (
  items: ReturnType<typeof accountFixture>[],
  page: number,
  totalPages: number,
) => ({
  total: items.length * totalPages,
  accounts: items,
  pagination: {
    current_page: page,
    per_page: items.length,
    total_pages: totalPages,
    total_count: items.length * totalPages,
  },
});

describe('V1 accounts', () => {
  it('list passes query params', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, pageEnvelope([accountFixture('a1')], 1, 1)),
    ]);
    const c = new V1Client({ apiKey: 'k', fetch });
    const r = await c.accounts.list({ page: 2, per_page: 10 });
    expect(r.accounts).toHaveLength(1);
    expect(calls[0]?.url).toContain('page=2');
    expect(calls[0]?.url).toContain('per_page=10');
  });

  it('get hits /v1/accounts/{id}', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, accountFixture('a1'))]);
    const c = new V1Client({ apiKey: 'k', fetch });
    const r = await c.accounts.get('a1');
    expect(r.id).toBe('a1');
    expect(calls[0]?.url).toContain('/v1/accounts/a1');
  });

  it('get → 404 maps to AccountNotFoundError', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(404, { error_code: 'ACCOUNT_NOT_FOUND', error: 'gone' }),
    ]);
    const c = new V1Client({ apiKey: 'k', fetch });
    await expect(c.accounts.get('a1')).rejects.toBeInstanceOf(AccountNotFoundError);
  });

  it('listAll iterates across all pages', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, pageEnvelope([accountFixture('a1'), accountFixture('a2')], 1, 2)),
      jsonResponse(200, pageEnvelope([accountFixture('a3'), accountFixture('a4')], 2, 2)),
    ]);
    const c = new V1Client({ apiKey: 'k', fetch });
    const ids: string[] = [];
    for await (const a of c.accounts.listAll({ per_page: 2 })) ids.push(a.id);
    expect(ids).toEqual(['a1', 'a2', 'a3', 'a4']);
    expect(calls).toHaveLength(2);
  });

  it('list → 401 maps to UnauthorizedError', async () => {
    const { fetch } = makeFetchMock([jsonResponse(401, { error_code: 'UNAUTHORIZED' })]);
    const c = new V1Client({ apiKey: 'k', fetch });
    await expect(c.accounts.list()).rejects.toBeInstanceOf(UnauthorizedError);
  });
});

describe('V2 accounts', () => {
  it('list hits /v2/accounts', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, pageEnvelope([accountFixture('a1')], 1, 1)),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.accounts.list();
    expect(calls[0]?.url).toContain('/v2/accounts');
  });

  it('sync POSTs with idempotency key + body', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(202, {
        message: 'Sync started',
        sync_session_id: 'ss1',
        status: 'pending',
        started_at: '2026-04-28T19:21:00Z',
      }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.accounts.sync('a1', { idempotencyKey: 'IDEMP-1' });
    expect(r.sync_session_id).toBe('ss1');
    expect(getMethod(callAt(calls, 0))).toBe('POST');
    expect(calls[0]?.url).toContain('/v2/accounts/a1/sync');
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toBe('IDEMP-1');
    expect(getBody(callAt(calls, 0))).toEqual({});
  });

  it('sync auto-generates idempotency key when none provided', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(202, {
        message: 'Sync started',
        sync_session_id: 'ss1',
        status: 'pending',
        started_at: 't',
      }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.accounts.sync('a1');
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toMatch(/^[0-9a-f-]{36}$/);
  });

  it('sync 409 → SyncInProgressError', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(409, { error_code: 'SYNC_IN_PROGRESS', error: 'busy' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.accounts.sync('a1')).rejects.toBeInstanceOf(SyncInProgressError);
  });

  it('sync 429 → SyncRateLimitExceededError', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(429, { error_code: 'SYNC_RATE_LIMIT_EXCEEDED', retry_after: 30 }),
      jsonResponse(429, { error_code: 'SYNC_RATE_LIMIT_EXCEEDED', retry_after: 30 }),
      jsonResponse(429, { error_code: 'SYNC_RATE_LIMIT_EXCEEDED', retry_after: 30 }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch, retryPolicy: { maxAttempts: 1 } });
    await expect(c.accounts.sync('a1')).rejects.toBeInstanceOf(SyncRateLimitExceededError);
  });

  it('sync 503 → BankUnderMaintenanceError', async () => {
    const { fetch } = makeFetchMock([jsonResponse(503, { error_code: 'BANK_UNDER_MAINTENANCE' })]);
    const c = new V2Client({ apiKey: 'k', fetch, retryPolicy: { maxAttempts: 1 } });
    await expect(c.accounts.sync('a1')).rejects.toBeInstanceOf(BankUnderMaintenanceError);
  });

  it('sync 404 BANK_CONNECTION_NOT_FOUND', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(404, { error_code: 'BANK_CONNECTION_NOT_FOUND' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.accounts.sync('a1')).rejects.toBeInstanceOf(BankConnectionNotFoundError);
  });
});
