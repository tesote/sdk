import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import {
  ConfigError,
  RateLimitExceededError,
  ServiceUnavailableError,
  TesoteError,
} from '../src/errors.js';
import { InMemoryLRUCache, Transport } from '../src/transport.js';

interface FetchCall {
  url: string;
  init: RequestInit;
}

function makeFetchMock(responses: ReadonlyArray<Response | (() => Response | Promise<Response>)>): {
  fetch: typeof fetch;
  calls: FetchCall[];
} {
  const calls: FetchCall[] = [];
  let i = 0;
  const fn = vi.fn(async (url: string | URL | Request, init: RequestInit = {}) => {
    calls.push({ url: String(url), init });
    const next = responses[i++];
    if (next === undefined) throw new Error(`fetch called more times than mocked (${i})`);
    return typeof next === 'function' ? await next() : next;
  }) as unknown as typeof fetch;
  return { fetch: fn, calls };
}

function jsonResponse(
  status: number,
  body: unknown,
  headers: Record<string, string> = {},
): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'content-type': 'application/json', ...headers },
  });
}

describe('Transport — construction', () => {
  it('throws ConfigError when apiKey missing', () => {
    expect(() => new Transport({ apiKey: '' } as never)).toThrow(ConfigError);
  });
});

describe('Transport — bearer + headers', () => {
  it('sends Authorization, Accept, User-Agent on every request', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, { ok: true })]);
    const t = new Transport({ apiKey: 'sk_test_abcd1234', fetch });
    await t.request({ method: 'GET', path: '/v3/accounts' });
    expect(calls).toHaveLength(1);
    const headers = new Headers(calls[0]?.init.headers as HeadersInit);
    expect(headers.get('authorization')).toBe('Bearer sk_test_abcd1234');
    expect(headers.get('accept')).toBe('application/json');
    expect(headers.get('user-agent')).toMatch(/^@tesote\/sdk-ts\/0\.1\.0 \(node\//);
  });

  it('hits the default base URL when none provided', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, {})]);
    const t = new Transport({ apiKey: 'k', fetch });
    await t.request({ method: 'GET', path: '/v3/accounts' });
    expect(calls[0]?.url).toBe('https://equipo.tesote.com/api/v3/accounts');
  });

  it('strips trailing slashes from baseUrl and prepends path', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, {})]);
    const t = new Transport({ apiKey: 'k', baseUrl: 'https://x.example.com/api/', fetch });
    await t.request({ method: 'GET', path: 'v3/accounts' });
    expect(calls[0]?.url).toBe('https://x.example.com/api/v3/accounts');
  });

  it('serializes query params alphabetically', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, {})]);
    const t = new Transport({ apiKey: 'k', fetch });
    await t.request({
      method: 'GET',
      path: '/v3/accounts',
      query: { z: '1', a: '2', skipped: null },
    });
    expect(calls[0]?.url).toContain('?a=2&z=1');
    expect(calls[0]?.url).not.toContain('skipped');
  });
});

describe('Transport — retries', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });
  afterEach(() => {
    vi.useRealTimers();
  });

  it('retries on 429 then succeeds; honors Retry-After', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(
        429,
        { error: 'slow down', error_code: 'RATE_LIMIT_EXCEEDED', retry_after: 1 },
        { 'retry-after': '1', 'x-request-id': 'req-1' },
      ),
      jsonResponse(200, { ok: true }, { 'x-request-id': 'req-2' }),
    ]);
    const t = new Transport({
      apiKey: 'k',
      fetch,
      retryPolicy: { maxAttempts: 3, baseDelay: 10, maxDelay: 100 },
    });
    const promise = t.request({ method: 'GET', path: '/v3/accounts' });
    await vi.advanceTimersByTimeAsync(2000);
    const res = await promise;
    expect(calls).toHaveLength(2);
    expect(res.status).toBe(200);
    expect(res.requestId).toBe('req-2');
  });

  it('retries on 502 then surfaces the final error after exhaustion', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(502, {}, { 'x-request-id': 'r1' }),
      jsonResponse(502, {}, { 'x-request-id': 'r2' }),
      jsonResponse(502, {}, { 'x-request-id': 'r3' }),
    ]);
    const t = new Transport({
      apiKey: 'k',
      fetch,
      retryPolicy: { maxAttempts: 3, baseDelay: 1, maxDelay: 5 },
    });
    const promise = t.request({ method: 'GET', path: '/v3/accounts' }).catch((e: unknown) => e);
    await vi.advanceTimersByTimeAsync(50);
    const err = (await promise) as TesoteError;
    expect(err).toBeInstanceOf(TesoteError);
    expect(err.attempts).toBe(3);
    expect(err.requestId).toBe('r3');
    expect(calls).toHaveLength(3);
  });

  it('raises RateLimitExceededError after exhaustion on 429', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(429, { error_code: 'RATE_LIMIT_EXCEEDED', retry_after: 0 }),
      jsonResponse(429, { error_code: 'RATE_LIMIT_EXCEEDED', retry_after: 0 }),
      jsonResponse(429, { error_code: 'RATE_LIMIT_EXCEEDED', retry_after: 0 }),
    ]);
    const t = new Transport({
      apiKey: 'k',
      fetch,
      retryPolicy: { maxAttempts: 3, baseDelay: 1, maxDelay: 5 },
    });
    const promise = t.request({ method: 'GET', path: '/v3/accounts' }).catch((e: unknown) => e);
    await vi.advanceTimersByTimeAsync(50);
    const err = await promise;
    expect(err).toBeInstanceOf(RateLimitExceededError);
  });

  it('raises ServiceUnavailableError on 503 after retries', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(503, {}),
      jsonResponse(503, {}),
      jsonResponse(503, {}),
    ]);
    const t = new Transport({
      apiKey: 'k',
      fetch,
      retryPolicy: { maxAttempts: 3, baseDelay: 1, maxDelay: 5 },
    });
    const promise = t.request({ method: 'GET', path: '/v3/accounts' }).catch((e: unknown) => e);
    await vi.advanceTimersByTimeAsync(50);
    expect(await promise).toBeInstanceOf(ServiceUnavailableError);
  });

  it('does not retry on 400', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(400, { error: 'bad', error_code: 'UNPROCESSABLE_CONTENT' }),
    ]);
    const t = new Transport({
      apiKey: 'k',
      fetch,
      retryPolicy: { maxAttempts: 5, baseDelay: 1, maxDelay: 5 },
    });
    await expect(t.request({ method: 'GET', path: '/v3/accounts' })).rejects.toBeInstanceOf(
      TesoteError,
    );
    expect(calls).toHaveLength(1);
  });
});

describe('Transport — rate-limit capture', () => {
  it('exposes lastRateLimit from response headers', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(
        200,
        {},
        {
          'x-ratelimit-limit': '200',
          'x-ratelimit-remaining': '199',
          'x-ratelimit-reset': '1700000000',
        },
      ),
    ]);
    const t = new Transport({ apiKey: 'k', fetch });
    await t.request({ method: 'GET', path: '/v3/accounts' });
    expect(t.lastRateLimit).toEqual({ limit: 200, remaining: 199, reset: 1700000000 });
  });
});

describe('Transport — request-id propagation into errors', () => {
  it('attaches X-Request-Id to thrown errors', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(
        422,
        { error: 'bad', error_code: 'UNPROCESSABLE_CONTENT' },
        { 'x-request-id': 'abc-123' },
      ),
    ]);
    const t = new Transport({ apiKey: 'k', fetch });
    const err = (await t
      .request({ method: 'GET', path: '/v3/accounts' })
      .catch((e: unknown) => e)) as TesoteError;
    expect(err.requestId).toBe('abc-123');
    expect(err.requestSummary?.method).toBe('GET');
    expect(err.requestSummary?.path).toBe('/v3/accounts');
  });
});

describe('Transport — idempotency-key', () => {
  it('auto-generates an Idempotency-Key for POST when none provided', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, {})]);
    const t = new Transport({ apiKey: 'k', fetch });
    await t.request({ method: 'POST', path: '/v3/accounts/x/sync' });
    const headers = new Headers(calls[0]?.init.headers as HeadersInit);
    const key = headers.get('idempotency-key');
    expect(key).toBeTruthy();
    expect(key).toMatch(/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/);
  });

  it('uses the caller-supplied idempotencyKey when provided', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, {})]);
    const t = new Transport({ apiKey: 'k', fetch });
    await t.request({ method: 'POST', path: '/v3/x', idempotencyKey: 'my-key' });
    const headers = new Headers(calls[0]?.init.headers as HeadersInit);
    expect(headers.get('idempotency-key')).toBe('my-key');
  });

  it('does not set Idempotency-Key on GET', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, {})]);
    const t = new Transport({ apiKey: 'k', fetch });
    await t.request({ method: 'GET', path: '/v3/accounts' });
    const headers = new Headers(calls[0]?.init.headers as HeadersInit);
    expect(headers.get('idempotency-key')).toBeNull();
  });
});

describe('Transport — body serialization', () => {
  it('JSON-encodes object bodies and sets Content-Type', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, {})]);
    const t = new Transport({ apiKey: 'k', fetch });
    await t.request({ method: 'POST', path: '/v3/x', body: { name: 'a' } });
    const headers = new Headers(calls[0]?.init.headers as HeadersInit);
    expect(headers.get('content-type')).toBe('application/json');
    expect(calls[0]?.init.body).toBe('{"name":"a"}');
  });
});

describe('Transport — TTL cache', () => {
  it('returns cached response on repeat GET within TTL', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { hit: 1 }),
      jsonResponse(200, { hit: 2 }),
    ]);
    const t = new Transport({ apiKey: 'k', fetch, cacheBackend: new InMemoryLRUCache() });
    const r1 = await t.request<{ hit: number }>({
      method: 'GET',
      path: '/v3/accounts',
      cache: { ttl: 60 },
    });
    const r2 = await t.request<{ hit: number }>({
      method: 'GET',
      path: '/v3/accounts',
      cache: { ttl: 60 },
    });
    expect(r1.data.hit).toBe(1);
    expect(r2.data.hit).toBe(1);
    expect(calls).toHaveLength(1);
  });

  it('busts cache after a mutation', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { hit: 1 }),
      jsonResponse(200, { mutated: true }),
      jsonResponse(200, { hit: 2 }),
    ]);
    const t = new Transport({ apiKey: 'k', fetch });
    await t.request({ method: 'GET', path: '/v3/accounts', cache: { ttl: 60 } });
    await t.request({ method: 'POST', path: '/v3/accounts/x/sync' });
    const r3 = await t.request<{ hit: number }>({
      method: 'GET',
      path: '/v3/accounts',
      cache: { ttl: 60 },
    });
    expect(r3.data.hit).toBe(2);
    expect(calls).toHaveLength(3);
  });

  it('cache: false bypasses cache for that call', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { hit: 1 }),
      jsonResponse(200, { hit: 2 }),
    ]);
    const t = new Transport({ apiKey: 'k', fetch });
    await t.request({ method: 'GET', path: '/v3/accounts', cache: { ttl: 60 } });
    await t.request({ method: 'GET', path: '/v3/accounts', cache: false });
    expect(calls).toHaveLength(2);
  });
});

describe('Transport — log hook redaction', () => {
  it('emits redacted bearer (Bearer <last4>), never the full key', async () => {
    const events: { authorization?: string }[] = [];
    const { fetch } = makeFetchMock([jsonResponse(200, {})]);
    const t = new Transport({
      apiKey: 'sk_supersecret_1234',
      fetch,
      log: (e) => events.push(e),
    });
    await t.request({ method: 'GET', path: '/v3/accounts' });
    for (const e of events) {
      expect(e.authorization).toBe('Bearer 1234');
      expect(e.authorization).not.toContain('supersecret');
    }
    expect(events.length).toBeGreaterThan(0);
  });
});
