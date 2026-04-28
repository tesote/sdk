# Transport

The transport is the only place that touches the network. Resource clients call `transport.request(method, path, params, body, opts)` and get back parsed JSON or a typed error. Everything else listed here is the transport's job — never duplicate it in a resource client.

## Responsibilities

| Concern              | Behavior |
|----------------------|----------|
| Auth                 | Inject `Authorization: Bearer <api_key>` on every request |
| Base URL             | Configurable; default `https://equipo.tesote.com/api`; per-version path appended (`/v3/...`) |
| Content negotiation  | `Accept: application/json`; `Content-Type: application/json` on POST/PUT/PATCH |
| User-Agent           | `tesote-sdk-<lang>/<sdk_version> (<runtime>)` — opaque to the server, but useful in support tickets |
| Request ID           | Surface response `X-Request-Id` to callers and into every thrown error |
| Caching              | See below |
| Retries              | See below |
| Rate-limit awareness | See below |
| Idempotency          | See below |
| Pagination           | See below |
| Logging hook         | Single callback (request → response or error). Never log the bearer token; redact to `Bearer <last4>` |

## Retries

Default: 3 retries, exponential backoff with jitter (e.g. `min(cap, base * 2^attempt) ± rand`), `base = 250ms`, `cap = 8s`. Configurable via `RetryPolicy { maxAttempts, baseDelay, maxDelay, retryOn }`.

Retry on:
- `429` — wait `Retry-After` if present, otherwise backoff.
- `502 / 503 / 504` — backoff.
- Network errors: connection reset, DNS failure, idempotent-method timeout.

**Never retry on:**
- `4xx` other than 429.
- Non-idempotent methods (POST without an idempotency key) on read timeouts — surface to the caller, the request may have succeeded server-side.

When retries are exhausted, raise the typed error from the last attempt with `attempts: n` attached.

## Rate-limit awareness

- Read `X-RateLimit-Remaining` and `X-RateLimit-Reset` from every response; expose as `client.lastRateLimit`.
- On `429`, raise `RateLimitExceededError` only after retries are exhausted. While retrying, sleep for `Retry-After` seconds.
- Optional opt-in: `RateLimiter` mode that *proactively* slows requests when `Remaining` falls under a threshold (e.g. 10). Off by default.

Limits to design against (from the API): 200 req/min per API key, 400 req/min per IP.

## Caching

Two layers, both opt-in:

1. **HTTP-conditional** — for any `GET`, send `If-None-Match` / `If-Modified-Since` from the previous response when present, treat `304` as a cache hit.
2. **TTL response cache** — in-memory LRU keyed on `(method, path, sorted_query, accept_header)`. Default off; opt-in per resource (`client.accounts.list({ cache: { ttl: 30 } })`). Mutations on the same resource path bust matching keys.

Both layers must be:
- Pluggable — accept a `CacheBackend` interface so users can drop in Redis/memcached.
- Disabled when `Authorization` could change scope; key the cache by API-key-id-hash to avoid cross-tenant bleed.
- Bypassed when the caller passes `cache: false`.

## Idempotency

Every mutating endpoint (POST/PUT/PATCH/DELETE) accepts an optional `idempotencyKey` argument. The transport:

- Sends it as `Idempotency-Key: <uuid>` header.
- Auto-generates one (UUIDv4) for SDK-driven retries when the caller didn't pass one.
- Caches the in-flight response for ~24h so a retry returns the same result instead of double-creating.

Endpoints that mutate without a natural idempotency key (e.g. `POST /v3/accounts/:id/sync`) still accept the header — the server dedupes server-side.

## Pagination

API uses cursor pagination. Transport exposes:

- `client.transactions.list({ ... })` → returns one page + cursor metadata.
- `client.transactions.listAll({ ... })` → async iterator / generator that walks all pages.

Mid-iteration mutations to the underlying dataset surface as `MutationDuringPaginationError` (HTTP 409 from the API). Consumer choice: restart with the new cursor or abort.

## Timeouts

Defaults: `connectTimeout = 5s`, `readTimeout = 30s`. Both configurable. The default applies even to `listAll()` — each individual page request gets its own timeout, not a wall-clock for the whole walk.

## What the transport does NOT do

- Schema validation of request bodies — type the public methods and trust them. Server returns 422 for bad input, which the SDK turns into a typed error.
- Business logic — no "smart" merging of paginated results into custom shapes, no auto-creating dependent resources.
- Persistence — no on-disk cache by default. Pluggable backend, off by default.
