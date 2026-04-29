import { describe, expect, it } from 'vitest';
import {
  AccountNotFoundError,
  BankConnectionNotFoundError,
  SyncSessionNotFoundError,
} from '../src/errors.js';
import { V2Client } from '../src/index.js';
import { jsonResponse, makeFetchMock } from './helpers.js';

const session = (id: string, status: 'completed' | 'failed' | 'started' = 'completed') => ({
  id,
  status,
  started_at: '2026-04-28T12:00:00Z',
  completed_at: status === 'completed' ? '2026-04-28T12:00:30Z' : null,
  transactions_synced: 5,
  accounts_count: 1,
  error: status === 'failed' ? { type: 'NetError', message: 'down' } : null,
  performance: null,
});

describe('V2 syncSessions', () => {
  it('list builds query and yields sessions', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, {
        sync_sessions: [session('s1'), session('s2')],
        limit: 50,
        offset: 0,
        has_more: false,
      }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.syncSessions.list('a1', { status: 'completed' });
    expect(r.sync_sessions).toHaveLength(2);
    expect(calls[0]?.url).toContain('/v2/accounts/a1/sync_sessions');
    expect(calls[0]?.url).toContain('status=completed');
  });

  it('list → 404 ACCOUNT_NOT_FOUND', async () => {
    const { fetch } = makeFetchMock([jsonResponse(404, { error_code: 'ACCOUNT_NOT_FOUND' })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.syncSessions.list('a1')).rejects.toBeInstanceOf(AccountNotFoundError);
  });

  it('list → 404 BANK_CONNECTION_NOT_FOUND', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(404, { error_code: 'BANK_CONNECTION_NOT_FOUND' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.syncSessions.list('a1')).rejects.toBeInstanceOf(BankConnectionNotFoundError);
  });

  it('get returns session', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, session('s1'))]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.syncSessions.get('a1', 's1');
    expect(r.id).toBe('s1');
    expect(calls[0]?.url).toContain('/v2/accounts/a1/sync_sessions/s1');
  });

  it('get → 404 SYNC_SESSION_NOT_FOUND', async () => {
    const { fetch } = makeFetchMock([jsonResponse(404, { error_code: 'SYNC_SESSION_NOT_FOUND' })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.syncSessions.get('a1', 's1')).rejects.toBeInstanceOf(SyncSessionNotFoundError);
  });

  it('listAll iterates offset pagination', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, {
        sync_sessions: [session('s1'), session('s2')],
        limit: 2,
        offset: 0,
        has_more: true,
      }),
      jsonResponse(200, {
        sync_sessions: [session('s3')],
        limit: 2,
        offset: 2,
        has_more: false,
      }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const ids: string[] = [];
    for await (const s of c.syncSessions.listAll('a1', { limit: 2 })) ids.push(s.id);
    expect(ids).toEqual(['s1', 's2', 's3']);
    expect(calls[1]?.url).toContain('offset=2');
  });
});
