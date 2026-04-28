/**
 * Public type definitions for the Transport.
 * Lives in its own file to keep transport.ts under the 500 LOC budget.
 */

export interface RetryPolicy {
  maxAttempts: number;
  baseDelay: number;
  maxDelay: number;
  retryOn: ReadonlySet<number>;
}

export interface RateLimitSnapshot {
  limit: number | null;
  remaining: number | null;
  reset: number | null;
}

export interface CacheEntry {
  body: string;
  storedAt: number;
  ttlMs: number;
  contentType: string | null;
}

export interface CacheBackend {
  get(key: string): CacheEntry | null | Promise<CacheEntry | null>;
  set(key: string, entry: CacheEntry): void | Promise<void>;
  delete(key: string): void | Promise<void>;
  clear(): void | Promise<void>;
}

export interface CacheOptions {
  /** TTL in seconds. */
  ttl?: number;
  /** Pass `false` to bypass cache for a single call. */
  enabled?: boolean;
}

export interface LogEvent {
  phase: 'request' | 'response' | 'error' | 'retry';
  method: string;
  path: string;
  url: string;
  attempt: number;
  status?: number;
  requestId?: string | null;
  /** Pre-redacted authorization header. */
  authorization?: string;
  durationMs?: number;
  error?: unknown;
}

export type LogHook = (event: LogEvent) => void;

export interface TransportOptions {
  apiKey: string;
  baseUrl?: string;
  userAgent?: string;
  /** Connect timeout, ms. Default 5000. */
  connectTimeout?: number;
  /** Read timeout, ms. Default 30000. */
  readTimeout?: number;
  retryPolicy?: Partial<RetryPolicy>;
  cacheBackend?: CacheBackend;
  /** Custom fetch (tests). Defaults to globalThis.fetch. */
  fetch?: typeof fetch;
  log?: LogHook;
}

export interface RequestOptions {
  method?: string;
  path: string;
  query?: Record<string, string | number | boolean | null | undefined>;
  body?: unknown;
  headers?: Record<string, string>;
  idempotencyKey?: string;
  cache?: CacheOptions | false;
  /** Override for read timeout on this single call. */
  readTimeout?: number;
  /** Override for connect timeout on this single call. */
  connectTimeout?: number;
}

export interface ResponseEnvelope<T> {
  status: number;
  headers: Headers;
  data: T;
  requestId: string | null;
  rateLimit: RateLimitSnapshot;
}
