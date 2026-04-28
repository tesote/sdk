/**
 * Internal helpers for the Transport: UA, error classification, UUID,
 * sleep, backoff, querystring, rate-limit parsing, JSON safety.
 * Kept private to keep transport.ts under the 500 LOC budget.
 */

import { NetworkError, type TimeoutError, TlsError } from './errors.js';
import type { RateLimitSnapshot, RetryPolicy } from './transport_types.js';

export const SDK_VERSION = '0.1.0';
export const DEFAULT_BASE_URL = 'https://equipo.tesote.com/api';

const DEFAULT_RETRY_STATUSES = new Set([429, 502, 503, 504]);

export const MUTATING_METHODS = new Set(['POST', 'PUT', 'PATCH', 'DELETE']);

export const DEFAULT_RETRY_POLICY: RetryPolicy = {
  maxAttempts: 3,
  baseDelay: 250,
  maxDelay: 8_000,
  retryOn: DEFAULT_RETRY_STATUSES,
};

export function buildUserAgent(custom: string | undefined): string {
  if (custom !== undefined && custom.length > 0) return custom;
  const node = typeof process !== 'undefined' ? process.version : 'unknown';
  return `@tesote/sdk-ts/${SDK_VERSION} (node/${node})`;
}

export function classifyFetchError(err: unknown): NetworkError | TimeoutError | TlsError {
  if (err instanceof Error) {
    const msg = err.message.toLowerCase();
    if (
      msg.includes('certificate') ||
      msg.includes('cert') ||
      msg.includes('tls') ||
      msg.includes('ssl')
    ) {
      return new TlsError({ errorCode: 'TLS_ERROR', message: err.message, cause: err });
    }
  }
  return new NetworkError({
    errorCode: 'NETWORK_ERROR',
    message: err instanceof Error ? err.message : 'network error',
    cause: err,
  });
}

export function uuidv4(): string {
  // why: crypto.randomUUID is in Node 14.17+; safe for the Node 18 floor.
  if (typeof globalThis.crypto?.randomUUID === 'function') {
    return globalThis.crypto.randomUUID();
  }
  // why: defensive fallback; never expected to run on Node 18+.
  const bytes = new Uint8Array(16);
  for (let i = 0; i < 16; i++) bytes[i] = Math.floor(Math.random() * 256);
  if (bytes[6] !== undefined) bytes[6] = (bytes[6] & 0x0f) | 0x40;
  if (bytes[8] !== undefined) bytes[8] = (bytes[8] & 0x3f) | 0x80;
  const hex: string[] = [];
  for (const b of bytes) hex.push(b.toString(16).padStart(2, '0'));
  return `${hex.slice(0, 4).join('')}-${hex.slice(4, 6).join('')}-${hex.slice(6, 8).join('')}-${hex.slice(8, 10).join('')}-${hex.slice(10, 16).join('')}`;
}

export function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(() => resolve(), ms));
}

export function backoffDelay(attempt: number, policy: RetryPolicy): number {
  const base = policy.baseDelay * 2 ** (attempt - 1);
  const capped = Math.min(policy.maxDelay, base);
  // why: full jitter, per AWS architecture blog; spreads thundering herd.
  return Math.floor(Math.random() * capped);
}

export function buildQueryString(
  query: Record<string, string | number | boolean | null | undefined> | undefined,
): string {
  if (query === undefined) return '';
  const params = new URLSearchParams();
  const keys = Object.keys(query).sort();
  for (const k of keys) {
    const v = query[k];
    if (v === null || v === undefined) continue;
    params.append(k, String(v));
  }
  const s = params.toString();
  return s.length === 0 ? '' : `?${s}`;
}

export function readRateLimit(headers: Headers): RateLimitSnapshot {
  const num = (h: string): number | null => {
    const v = headers.get(h);
    if (v === null) return null;
    const n = Number(v);
    return Number.isFinite(n) ? n : null;
  };
  return {
    limit: num('X-RateLimit-Limit'),
    remaining: num('X-RateLimit-Remaining'),
    reset: num('X-RateLimit-Reset'),
  };
}

export function safeJsonParse(text: string): unknown {
  if (text.length === 0) return null;
  try {
    return JSON.parse(text) as unknown;
  } catch {
    return null;
  }
}

export function bodyShape(body: unknown): string | undefined {
  if (body === undefined || body === null) return undefined;
  if (Array.isArray(body)) return `array(${body.length})`;
  if (typeof body === 'object') return `object(${Object.keys(body as object).length} keys)`;
  return typeof body;
}

const SDK_ERROR_NAMES = new Set([
  'TesoteError',
  'ApiError',
  'UnauthorizedError',
  'ApiKeyRevokedError',
  'WorkspaceSuspendedError',
  'AccountDisabledError',
  'HistorySyncForbiddenError',
  'MutationDuringPaginationError',
  'UnprocessableContentError',
  'InvalidDateRangeError',
  'RateLimitExceededError',
  'ServiceUnavailableError',
  'TransportError',
  'NetworkError',
  'TimeoutError',
  'TlsError',
  'ConfigError',
  'EndpointRemovedError',
]);

export function isTypedSdkError(err: unknown): boolean {
  if (!(err instanceof Error)) return false;
  // why: avoid an import cycle with errors.ts; check by name (set in TesoteError).
  return SDK_ERROR_NAMES.has(err.name);
}
