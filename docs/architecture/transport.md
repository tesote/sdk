# Transport

The only place that touches the network. Resource clients call `transport.request(method, path, params, body, opts)` and get parsed JSON or a typed error. Everything below is the transport's job — never duplicate in a resource client.

## Zero runtime deps

Customers don't want a wire-protocol SDK pulling in transitive deps that conflict with their lock files. Stdlib only.

| Language   | HTTP                          | JSON                                                  |
|------------|-------------------------------|-------------------------------------------------------|
| TypeScript | native `fetch` (Node 18+)     | native `JSON`                                          |
| Python     | `urllib.request`              | `json`                                                 |
| Ruby       | `Net::HTTP`                   | `json`                                                 |
| Java       | `java.net.http.HttpClient`    | `jakarta.json` preferred; `jackson-databind` allowed if too awkward — only acceptable runtime dep across all languages |
| PHP        | ext-curl                      | `json_*`                                               |
| Go         | `net/http`                    | `encoding/json`                                        |

## Responsibilities

| Concern              | Behavior |
|----------------------|----------|
| Auth                 | Inject `Authorization: Bearer <api_key>` on every request |
| Base URL             | Configurable; default `https://equipo.tesote.com/api`; per-version path appended (`/v3/...`) |
| Content negotiation  | `Accept: application/json`; `Content-Type: application/json` on POST/PUT/PATCH |
| User-Agent           | `tesote-sdk-<lang>/<sdk_version> (<runtime>/<runtime_version>)` — captured at client construction |
| Request ID           | Surface response `X-Request-Id` to callers and into every thrown error |
| Timeouts             | Connect 5s, read 30s; both configurable via the language's stdlib mechanism |
| Keep-alive           | HTTP/1.1 keep-alive enabled; reuse one connection pool per `Client` instance. Per language: TS share one `Agent`/`Dispatcher`; Python single `http.client.HTTPSConnection` reused or `urllib` opener with persistent connection; Ruby `Net::HTTP.start` block / persistent instance; Java single `HttpClient` (it pools internally); PHP curl multi-handle reused; Go single `*http.Client` with `Transport.MaxIdleConnsPerHost`. |
| Lifecycle            | Each `Client` is OOP-instantiable; create as many as needed in one process. Idiomatic close: TS `client.close()` (releases agent + cache), Python `client.close()` plus `__enter__`/`__exit__` for `with`, Ruby `client.close`, Java `client.close()` (implements `AutoCloseable`), PHP `$client->close()` (`__destruct` fallback), Go `client.Close()` (releases idle connections). After `close`, calls raise `ConfigError` (or language-equivalent "closed client" error). |
| Caching              | See below |
| Retries              | See below |
| Rate-limit awareness | See below |
| Idempotency          | See below |
| Pagination           | See below |
| Logging hook         | Single callback (request → response or error). Never log the bearer token; redact to `Bearer <last4>` |

## Retries

Default: 3 retries, exponential backoff with jitter (`min(cap, base * 2^attempt) ± rand`), `base = 250ms`, `cap = 8s`. Configurable via `RetryPolicy { maxAttempts, baseDelay, maxDelay, retryOn }`.

Retry on:
- `429` — wait `Retry-After` if present, otherwise backoff.
- `502 / 503 / 504` — backoff.
- Network errors: connection reset, DNS failure, idempotent-method timeout.

Never retry on:
- `4xx` other than 429.
- Non-idempotent methods (POST without an idempotency key) on read timeouts — surface to the caller; the request may have succeeded server-side.

When retries exhausted: raise the typed error from the last attempt with `attempts: n` attached.

## Rate-limit awareness

- Read `X-RateLimit-Remaining` / `X-RateLimit-Reset` from every response; expose as `client.lastRateLimit`.
- On `429`, raise `RateLimitExceededError` only after retries exhausted. While retrying, sleep `Retry-After` seconds.
- Optional opt-in: `RateLimiter` mode that proactively slows requests when `Remaining` < threshold (e.g. 10). Off by default.

API limits: 200 req/min per API key, 400 req/min per IP.

## Caching

Two layers, both opt-in:

1. **HTTP-conditional** — for any `GET`, send `If-None-Match` / `If-Modified-Since` from the previous response when present; treat `304` as a cache hit.
2. **TTL response cache** — in-memory LRU keyed on `(method, path, sorted_query, accept_header)`. Default off; opt-in per resource (`client.accounts.list({ cache: { ttl: 30 } })`). Mutations on the same resource path bust matching keys.

Both layers must:
- Be pluggable — accept a `CacheBackend` interface so users can drop in Redis/memcached.
- Key by API-key-id-hash to avoid cross-tenant bleed; disabled when `Authorization` could change scope.
- Bypass when caller passes `cache: false`.

## Idempotency

Every mutating endpoint (POST/PUT/PATCH/DELETE) accepts an optional `idempotencyKey` argument. Transport:

- Sends as `Idempotency-Key: <uuid>` header.
- Auto-generates UUIDv4 for SDK-driven retries when caller didn't pass one.
- Caches the in-flight response for ~24h so a retry returns the same result instead of double-creating.

Endpoints mutating without a natural idempotency key (e.g. `POST /v3/accounts/:id/sync`) still accept the header — server dedupes server-side.

## Pagination

API uses cursor pagination. Transport exposes:

- `client.transactions.list({ ... })` → one page + cursor metadata.
- `client.transactions.listAll({ ... })` → async iterator/generator walking all pages.

Mid-iteration mutations to the underlying dataset surface as `MutationDuringPaginationError` (HTTP 409). Consumer choice: restart with the new cursor or abort.

## What the transport does NOT do

- Schema validation of request bodies — type the public methods and trust them. Server returns 422 → SDK turns into a typed error.
- Business logic — no "smart" merging of paginated results, no auto-creating dependent resources.
- Persistence — no on-disk cache by default. Pluggable backend, off by default.
