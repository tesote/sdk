import { describe, expect, it } from 'vitest';
import { UnauthorizedError } from '../src/errors.js';
import { V1Client, V2Client } from '../src/index.js';
import { callAt, getMethod, jsonResponse, makeFetchMock } from './helpers.js';

describe('V1 status', () => {
  it('GET /status — anonymous probe', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { status: 'ok', authenticated: false }),
    ]);
    const c = new V1Client({ apiKey: 'k', fetch });
    const r = await c.status.status();
    expect(r).toEqual({ status: 'ok', authenticated: false });
    expect(getMethod(callAt(calls, 0))).toBe('GET');
    expect(calls[0]?.url).toContain('/status');
  });

  it('GET /whoami — returns client envelope', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(200, { client: { id: 'cid', name: 'acme', type: 'workspace' } }),
    ]);
    const c = new V1Client({ apiKey: 'k', fetch });
    const r = await c.status.whoami();
    expect(r.client.type).toBe('workspace');
    expect(r.client.id).toBe('cid');
  });

  it('whoami → 401 maps to UnauthorizedError', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(401, { error_code: 'UNAUTHORIZED', error: 'no' }),
    ]);
    const c = new V1Client({ apiKey: 'k', fetch });
    await expect(c.status.whoami()).rejects.toBeInstanceOf(UnauthorizedError);
  });
});

describe('V2 status', () => {
  it('GET /v2/status', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { status: 'ok', authenticated: false }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.status.status();
    expect(r.status).toBe('ok');
    expect(calls[0]?.url).toContain('/v2/status');
  });

  it('GET /v2/whoami', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { client: { id: 'cid', name: 'a', type: 'user' } }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.status.whoami();
    expect(r.client.type).toBe('user');
    expect(calls[0]?.url).toContain('/v2/whoami');
  });
});
