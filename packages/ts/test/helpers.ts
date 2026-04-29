import { expect, vi } from 'vitest';

export interface FetchCall {
  url: string;
  init: RequestInit;
}

export function makeFetchMock(
  responses: ReadonlyArray<Response | (() => Response | Promise<Response>)>,
): {
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

export function jsonResponse(
  status: number,
  body: unknown,
  headers: Record<string, string> = {},
): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'content-type': 'application/json', ...headers },
  });
}

export function rawResponse(
  status: number,
  body: string,
  headers: Record<string, string>,
): Response {
  return new Response(body, { status, headers });
}

export function noContent(status = 204, headers: Record<string, string> = {}): Response {
  // why: Response constructor forbids a non-null body on 204/205/304.
  return new Response(null, { status, headers });
}

/** Pull the i-th call, asserting it exists. Centralizes the runtime check so
 *  call sites stay free of non-null assertions (forbidden by biome). */
export function callAt(calls: FetchCall[], idx: number): FetchCall {
  const c = calls[idx];
  expect(c, `call #${idx} missing`).toBeDefined();
  return c as FetchCall;
}

export function getHeader(call: FetchCall, name: string): string | null {
  return new Headers(call.init.headers as HeadersInit).get(name);
}

export function getMethod(call: FetchCall): string {
  return (call.init.method ?? 'GET').toUpperCase();
}

export function getBody(call: FetchCall): unknown {
  const body = call.init.body;
  if (body === undefined || body === null) return undefined;
  if (typeof body === 'string') {
    try {
      return JSON.parse(body) as unknown;
    } catch {
      return body;
    }
  }
  return body;
}
