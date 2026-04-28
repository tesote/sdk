/**
 * Single HTTP client built on native fetch.
 * Owns: bearer injection, retries, rate-limit parsing, idempotency keys,
 * request-id propagation, opt-in TTL response cache, timeouts.
 *
 * Mirrors docs/architecture/transport.md exactly.
 *
 * Helpers + types live in cache.ts, transport_internals.ts, transport_types.ts
 * to keep this file under the 500 LOC budget.
 */

import { InMemoryLRUCache } from './cache.js';
import {
  ConfigError,
  NetworkError,
  type RequestSummary,
  TimeoutError,
  mapApiError,
  redactBearer,
} from './errors.js';
import {
  DEFAULT_BASE_URL,
  DEFAULT_RETRY_POLICY,
  MUTATING_METHODS,
  backoffDelay,
  bodyShape,
  buildQueryString,
  buildUserAgent,
  classifyFetchError,
  isTypedSdkError,
  readRateLimit,
  safeJsonParse,
  sleep,
  uuidv4,
} from './transport_internals.js';
import type {
  CacheBackend,
  CacheOptions,
  LogHook,
  RateLimitSnapshot,
  RequestOptions,
  ResponseEnvelope,
  RetryPolicy,
  TransportOptions,
} from './transport_types.js';

export { DEFAULT_BASE_URL, DEFAULT_RETRY_POLICY, SDK_VERSION } from './transport_internals.js';
export { InMemoryLRUCache } from './cache.js';
export type {
  CacheBackend,
  CacheEntry,
  CacheOptions,
  LogEvent,
  LogHook,
  RateLimitSnapshot,
  RequestOptions,
  ResponseEnvelope,
  RetryPolicy,
  TransportOptions,
} from './transport_types.js';

export class Transport {
  public lastRateLimit: RateLimitSnapshot | null = null;

  private readonly apiKey: string;
  private readonly baseUrl: string;
  private readonly userAgent: string;
  private readonly connectTimeout: number;
  private readonly readTimeout: number;
  private readonly retryPolicy: RetryPolicy;
  private readonly cacheBackend: CacheBackend;
  private readonly fetchImpl: typeof fetch;
  private readonly log: LogHook | undefined;
  private readonly redactedAuth: string;

  constructor(opts: TransportOptions) {
    if (typeof opts.apiKey !== 'string' || opts.apiKey.length === 0) {
      throw new ConfigError({ errorCode: 'CONFIG_MISSING_API_KEY', message: 'apiKey is required' });
    }
    const fetchFn =
      opts.fetch ?? (typeof globalThis.fetch === 'function' ? globalThis.fetch : undefined);
    if (fetchFn === undefined) {
      throw new ConfigError({
        errorCode: 'CONFIG_MISSING_FETCH',
        message: 'global fetch is unavailable; Node 18+ required or pass options.fetch',
      });
    }
    this.apiKey = opts.apiKey;
    this.baseUrl = (opts.baseUrl ?? DEFAULT_BASE_URL).replace(/\/+$/, '');
    this.userAgent = buildUserAgent(opts.userAgent);
    this.connectTimeout = opts.connectTimeout ?? 5_000;
    this.readTimeout = opts.readTimeout ?? 30_000;
    this.retryPolicy = { ...DEFAULT_RETRY_POLICY, ...(opts.retryPolicy ?? {}) };
    this.cacheBackend = opts.cacheBackend ?? new InMemoryLRUCache();
    this.fetchImpl = fetchFn.bind(globalThis);
    this.log = opts.log;
    this.redactedAuth = redactBearer(this.apiKey);
  }

  public async request<T = unknown>(opts: RequestOptions): Promise<ResponseEnvelope<T>> {
    const method = (opts.method ?? 'GET').toUpperCase();
    const qs = buildQueryString(opts.query);
    const path = opts.path.startsWith('/') ? opts.path : `/${opts.path}`;
    const url = `${this.baseUrl}${path}${qs}`;
    const isMutation = MUTATING_METHODS.has(method);

    const cacheOpts: CacheOptions | false = opts.cache ?? false;
    const cacheEnabled =
      method === 'GET' &&
      cacheOpts !== false &&
      cacheOpts.enabled !== false &&
      typeof cacheOpts.ttl === 'number' &&
      cacheOpts.ttl > 0;
    const cacheKey = cacheEnabled ? this.cacheKey(method, path, qs) : null;

    if (cacheKey !== null) {
      const hit = await this.cacheBackend.get(cacheKey);
      if (hit !== null) return this.responseFromCache<T>(hit);
    }

    const headers = this.buildHeaders(opts, isMutation);
    const bodyString = this.encodeBody(opts.body, headers);
    const summary = this.buildSummary(method, path, opts);

    return await this.runWithRetries<T>({
      url,
      method,
      path,
      headers,
      bodyString,
      summary,
      cacheKey,
      cacheOpts,
      isMutation,
      readTimeout: opts.readTimeout ?? this.readTimeout,
      connectTimeout: opts.connectTimeout ?? this.connectTimeout,
    });
  }

  private buildHeaders(opts: RequestOptions, isMutation: boolean): Headers {
    const headers = new Headers(opts.headers ?? {});
    headers.set('Authorization', `Bearer ${this.apiKey}`);
    headers.set('Accept', 'application/json');
    headers.set('User-Agent', this.userAgent);
    if (isMutation) {
      headers.set('Idempotency-Key', opts.idempotencyKey ?? uuidv4());
    }
    return headers;
  }

  private encodeBody(body: unknown, headers: Headers): string | undefined {
    if (body === undefined || body === null) return undefined;
    const json = JSON.stringify(body);
    if (!headers.has('Content-Type')) headers.set('Content-Type', 'application/json');
    return json;
  }

  private buildSummary(method: string, path: string, opts: RequestOptions): RequestSummary {
    const summary: RequestSummary = {
      method,
      path,
      authorization: this.redactedAuth,
    };
    if (opts.query !== undefined) summary.query = { ...opts.query };
    const shape = bodyShape(opts.body);
    if (shape !== undefined) summary.bodyShape = shape;
    return summary;
  }

  private responseFromCache<T>(hit: {
    body: string;
    contentType: string | null;
  }): ResponseEnvelope<T> {
    const data = safeJsonParse(hit.body) as T;
    const headers = new Headers();
    if (hit.contentType !== null) headers.set('content-type', hit.contentType);
    return {
      status: 200,
      headers,
      data,
      requestId: null,
      rateLimit: this.lastRateLimit ?? { limit: null, remaining: null, reset: null },
    };
  }

  private async runWithRetries<T>(args: {
    url: string;
    method: string;
    path: string;
    headers: Headers;
    bodyString: string | undefined;
    summary: RequestSummary;
    cacheKey: string | null;
    cacheOpts: CacheOptions | false;
    isMutation: boolean;
    readTimeout: number;
    connectTimeout: number;
  }): Promise<ResponseEnvelope<T>> {
    const policy = this.retryPolicy;
    let attempt = 0;
    let lastErr: unknown;

    while (attempt < policy.maxAttempts) {
      attempt += 1;
      const started = Date.now();
      this.log?.({
        phase: 'request',
        method: args.method,
        path: args.path,
        url: args.url,
        attempt,
        authorization: this.redactedAuth,
      });

      const controller = new AbortController();
      // why: native fetch lacks separate connect/read timeouts; use the larger of
      // the two as the abort, matching the architecture doc's "configurable" intent.
      const totalTimeoutMs = Math.max(args.readTimeout, args.connectTimeout);
      const timer = setTimeout(() => controller.abort(), totalTimeoutMs);

      try {
        const fetchInit: RequestInit = {
          method: args.method,
          headers: args.headers,
          signal: controller.signal,
        };
        if (args.bodyString !== undefined) fetchInit.body = args.bodyString;
        const res = await this.fetchImpl(args.url, fetchInit);
        clearTimeout(timer);

        const requestId = res.headers.get('X-Request-Id');
        const rateLimit = readRateLimit(res.headers);
        this.lastRateLimit = rateLimit;
        const text = await res.text();

        this.log?.({
          phase: 'response',
          method: args.method,
          path: args.path,
          url: args.url,
          attempt,
          status: res.status,
          requestId,
          authorization: this.redactedAuth,
          durationMs: Date.now() - started,
        });

        if (res.status >= 200 && res.status < 300) {
          const data = (safeJsonParse(text) ?? null) as T;
          await this.afterSuccess(args, text, res.headers.get('content-type'));
          return { status: res.status, headers: res.headers, data, requestId, rateLimit };
        }

        const apiErr = mapApiError({
          httpStatus: res.status,
          requestId,
          retryAfterHeader: res.headers.get('Retry-After'),
          responseBody: text,
          parsedBody: safeJsonParse(text),
          requestSummary: args.summary,
          attempts: attempt,
        });

        const shouldRetry = attempt < policy.maxAttempts && policy.retryOn.has(res.status);
        if (!shouldRetry) throw apiErr;

        // why: respect Retry-After when the server tells us to wait.
        const ra = apiErr.retryAfter;
        const delay =
          ra !== null && ra >= 0
            ? Math.min(policy.maxDelay, ra * 1000)
            : backoffDelay(attempt, policy);

        this.log?.({
          phase: 'retry',
          method: args.method,
          path: args.path,
          url: args.url,
          attempt,
          status: res.status,
          requestId,
          authorization: this.redactedAuth,
        });

        lastErr = apiErr;
        await sleep(delay);
      } catch (err) {
        clearTimeout(timer);
        const handled = await this.handleAttemptError({
          err,
          attempt,
          totalTimeoutMs,
          isMutation: args.isMutation,
          method: args.method,
          path: args.path,
          url: args.url,
          summary: args.summary,
        });
        if (handled.retry) {
          lastErr = handled.error;
          continue;
        }
        throw handled.error;
      }
    }

    throw lastErr ?? new NetworkError({ errorCode: 'NETWORK_ERROR', message: 'retries exhausted' });
  }

  private async afterSuccess(
    args: {
      method: string;
      path: string;
      cacheKey: string | null;
      cacheOpts: CacheOptions | false;
      isMutation: boolean;
    },
    text: string,
    contentType: string | null,
  ): Promise<void> {
    if (
      args.cacheKey !== null &&
      args.cacheOpts !== false &&
      typeof args.cacheOpts.ttl === 'number'
    ) {
      await this.cacheBackend.set(args.cacheKey, {
        body: text,
        storedAt: Date.now(),
        ttlMs: args.cacheOpts.ttl * 1000,
        contentType,
      });
    }
    if (args.isMutation) await this.cacheBackend.clear();
  }

  private async handleAttemptError(args: {
    err: unknown;
    attempt: number;
    totalTimeoutMs: number;
    isMutation: boolean;
    method: string;
    path: string;
    url: string;
    summary: RequestSummary;
  }): Promise<{ retry: boolean; error: unknown }> {
    const { err, attempt, totalTimeoutMs, isMutation, method, path, url, summary } = args;
    const policy = this.retryPolicy;

    if (err instanceof Error && err.name === 'AbortError') {
      const timeoutErr = new TimeoutError({
        errorCode: 'TIMEOUT',
        message: `request timed out after ${totalTimeoutMs}ms`,
        requestSummary: summary,
        attempts: attempt,
        cause: err,
      });
      this.log?.({
        phase: 'error',
        method,
        path,
        url,
        attempt,
        authorization: this.redactedAuth,
        error: timeoutErr,
      });
      // why: idempotent methods may retry on timeout; mutations must not (the
      // request may have committed server-side).
      if (attempt < policy.maxAttempts && !isMutation) {
        await sleep(backoffDelay(attempt, policy));
        return { retry: true, error: timeoutErr };
      }
      return { retry: false, error: timeoutErr };
    }

    if (isTypedSdkError(err)) {
      this.log?.({
        phase: 'error',
        method,
        path,
        url,
        attempt,
        authorization: this.redactedAuth,
        error: err,
      });
      return { retry: false, error: err };
    }

    const wrapped = classifyFetchError(err);
    const Ctor = wrapped.constructor as new (
      f: ConstructorParameters<typeof NetworkError>[0],
    ) => NetworkError;
    const finalErr = new Ctor({
      errorCode: wrapped.errorCode,
      message: wrapped.message,
      requestSummary: summary,
      attempts: attempt,
      cause: err,
    });
    this.log?.({
      phase: 'error',
      method,
      path,
      url,
      attempt,
      authorization: this.redactedAuth,
      error: finalErr,
    });
    if (attempt < policy.maxAttempts && !isMutation) {
      await sleep(backoffDelay(attempt, policy));
      return { retry: true, error: finalErr };
    }
    return { retry: false, error: finalErr };
  }

  private cacheKey(method: string, path: string, qs: string): string {
    // why: include API-key tail to prevent cross-tenant cache bleed.
    return `${method} ${path}${qs} :: ${this.apiKey.slice(-6)}`;
  }
}
