import { describe, expect, it, vi } from 'vitest';
import { V3Client } from '../../src/v3/index.js';

interface FetchCall {
  url: string;
  init: RequestInit;
}

function makeFetchMock(responses: ReadonlyArray<Response>): {
  fetch: typeof fetch;
  calls: FetchCall[];
} {
  const calls: FetchCall[] = [];
  let i = 0;
  const fn = vi.fn(async (url: string | URL | Request, init: RequestInit = {}) => {
    calls.push({ url: String(url), init });
    const next = responses[i++];
    if (next === undefined) throw new Error('fetch over-called');
    return next;
  }) as unknown as typeof fetch;
  return { fetch: fn, calls };
}

function jsonResponse(status: number, body: unknown): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'content-type': 'application/json' },
  });
}

describe('V3Client.accounts', () => {
  it('list() hits GET /v3/accounts and returns parsed body', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { data: [{ id: 'acct_1', name: 'Ops', currency: 'EUR' }] }),
    ]);
    const client = new V3Client({ apiKey: 'k', fetch });
    const res = await client.accounts.list();
    expect(calls[0]?.init.method).toBe('GET');
    expect(calls[0]?.url).toBe('https://equipo.tesote.com/api/v3/accounts');
    expect(res.data[0]?.id).toBe('acct_1');
  });

  it('list() forwards cursor + limit as query', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, { data: [] })]);
    const client = new V3Client({ apiKey: 'k', fetch });
    await client.accounts.list({ cursor: 'c1', limit: 50 });
    expect(calls[0]?.url).toContain('cursor=c1');
    expect(calls[0]?.url).toContain('limit=50');
  });

  it('get(id) URL-encodes the id', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { id: 'acct/1', name: 'x', currency: 'EUR' }),
    ]);
    const client = new V3Client({ apiKey: 'k', fetch });
    const acct = await client.accounts.get('acct/1');
    expect(calls[0]?.url).toBe('https://equipo.tesote.com/api/v3/accounts/acct%2F1');
    expect(acct.id).toBe('acct/1');
  });

  it('get(id) propagates X-Request-Id into errors', async () => {
    const errResp = new Response(JSON.stringify({ error: 'gone', error_code: 'UNAUTHORIZED' }), {
      status: 401,
      headers: { 'content-type': 'application/json', 'x-request-id': 'req-acc-1' },
    });
    const { fetch } = makeFetchMock([errResp]);
    const client = new V3Client({ apiKey: 'k', fetch });
    const err = await client.accounts.get('acct_1').catch((e: unknown) => e);
    expect(err).toMatchObject({
      errorCode: 'UNAUTHORIZED',
      requestId: 'req-acc-1',
      httpStatus: 401,
    });
  });

  it('exposes lastRateLimit via the V3Client', async () => {
    const r = new Response(JSON.stringify({ data: [] }), {
      status: 200,
      headers: {
        'content-type': 'application/json',
        'x-ratelimit-limit': '200',
        'x-ratelimit-remaining': '12',
        'x-ratelimit-reset': '1700000000',
      },
    });
    const { fetch } = makeFetchMock([r]);
    const client = new V3Client({ apiKey: 'k', fetch });
    await client.accounts.list();
    expect(client.lastRateLimit).toEqual({ limit: 200, remaining: 12, reset: 1700000000 });
  });
});
