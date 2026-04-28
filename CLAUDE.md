# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

Public monorepo for the official client SDKs of the **equipo.tesote.com** API. Greenfield — only `.idea/` exists today; everything below is the agreed shape, not what's on disk yet. When you scaffold, follow it.

Languages shipped from here. Repo-level name is `tesote-sdk`; the npm package uses the scoped form `@tesote/sdk`. Other registries use `tesote-sdk`.

| Language | Folder              | Package name              | Min version | Registry        |
|----------|---------------------|---------------------------|-------------|-----------------|
| TypeScript | `packages/ts/`     | `@tesote/sdk`             | Node 18     | npm             |
| Python   | `packages/python/` | `tesote-sdk`              | Python 3.9  | PyPI            |
| Ruby     | `packages/ruby/`   | `tesote-sdk`              | Ruby 3.0    | RubyGems        |
| Java     | `packages/java/`   | `com.tesote:sdk`          | Java 17     | Maven Central   |
| PHP      | `packages/php/`    | `tesote/sdk`              | PHP 8.1     | Packagist       |
| Go       | `packages/go/`     | `github.com/tesote/sdk/go` | Go 1.21    | proxy.golang.org |

Two axes, don't conflate them:

- **Min runtime version** (above): conservative, set the floor low so the SDK works on widely-deployed runtimes. No syntax/features younger than the floor. No experimental features. No clever tricks.
- **Dependencies + tooling**: pin to **latest stable**. `httpx`, `faraday`, `guzzle`, `jackson`, `vitest`, `pytest`, `rspec`, `junit-jupiter`, `phpunit`, etc. — newest released version at scaffold time. GitHub Actions: latest major (`actions/checkout@v4`, `actions/setup-node@v4`, `actions/setup-python@v5`, etc.). Test the SDK on a matrix that includes the latest LTS *and* the current stable for that language.

Each language is independently versioned, released, and tested. No cross-language code sharing — duplicate idiomatically rather than abstract.

## Source of truth for the API

Upstream OpenAPI specs live in the platform repo (sibling dir, do not import from it at build time — vendor a snapshot):

- v1: `../<platform>/engines/tesote_api/docs/openapi.yaml` — read-only accounts + transactions
- v2: `../<platform>/engines/tesote_api/docs/openapi_v2.yaml` — adds sync sessions, transaction orders, batches, payment methods, bulk, search
- v3: routes only (`engines/tesote_api/config/routes.rb`, controllers under `app/controllers/tesote_api/v3/`) — adds categories, counterparties, legal entities, connections, webhooks, reports, balance history, workspace, MCP. **No OpenAPI doc yet** — derive types from the controllers/serializers and flag missing pieces.

Only expose **client-facing, API-key-authenticated** endpoints — anything mounted under `TesoteApi::Engine` with `current_api_key` auth. Do not surface internal admin/session-cookie controllers.

## SDK shape (all languages)

Versioned clients live side-by-side; consumers pick a version explicitly:

```ts
import { V1Client, V2Client, V3Client } from '@tesote/sdk'
```

```python
from tesote_sdk import V1Client, V2Client, V3Client
```

```ruby
require 'tesote_sdk'
TesoteSdk::V2::Client.new(api_key: ...)
```

```go
import "github.com/tesote/sdk/go/v3"
```

**Back-compat is permanent.** v1 and v2 stay shipped even after v3 lands. Removing or renaming a public symbol in any version is a breaking change — don't.

### Modular layout per SDK

One module/file per resource (accounts, transactions, sync_sessions, transaction_orders, batches, payment_methods, categories, counterparties, legal_entities, connections, webhooks, reports, balance_history, workspace). Per SOLID/SRP:

- Transport layer separate from resource clients (one HTTP client, swappable for tests).
- Resource clients are thin: marshal params, call transport, deserialize.
- Errors are typed (one class per `error_code` from the API — see below).
- **Transport owns the cross-cutting concerns**: pagination, retry with exponential backoff + jitter, rate-limit awareness, response caching (ETag / `Cache-Control` / opt-in TTL for read-only resources), idempotency keys for mutating endpoints, request-id propagation. Resource clients never reimplement them.

### Raise good errors

Every error the SDK throws must carry enough context to debug without re-running the request. The minimum payload on every error class:

- `error_code` (string, from the API envelope)
- `message` (human-readable, from API or synthesized for transport errors)
- `http_status` (int)
- `request_id` (from `X-Request-Id`)
- `retry_after` (int seconds, when present)
- `response_body` (raw, for unexpected shapes)
- `request_summary` (method + path + redacted query/body — never log the bearer token)

Error class per `error_code` (table below). Don't collapse them into a single `ApiError` — callers should `catch RateLimitExceededError` or `catch WorkspaceSuspendedError` distinctly. Transport-level failures (DNS, TLS, timeout, connection reset) get their own typed classes (`NetworkError`, `TimeoutError`) — never bubble up the underlying language exception.

## API contract clients must implement

### Auth
`Authorization: Bearer <api_key>` on every request. No other auth schemes.

### Rate limits
- 200 req/min per API key, 400 req/min per IP.
- Headers on every response: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`.
- On 429: `Retry-After` header (seconds). SDKs should retry with backoff up to a configurable cap, surfacing the typed `RateLimitExceeded` error if the cap is hit.

### Error envelope
```json
{ "error": "...", "error_code": "...", "error_id": "uuid?", "retry_after": 60 }
```

Map each `error_code` to a typed exception. Known codes (from `engines/tesote_api/app/lib/tesote_api/error_codes.rb`):

| HTTP | error_code              | Class name suggestion          |
|------|-------------------------|--------------------------------|
| 401  | `UNAUTHORIZED`          | `UnauthorizedError`            |
| 401  | `API_KEY_REVOKED`       | `ApiKeyRevokedError`           |
| 403  | `WORKSPACE_SUSPENDED`   | `WorkspaceSuspendedError`      |
| 403  | `ACCOUNT_DISABLED`      | `AccountDisabledError`         |
| 403  | `HISTORY_SYNC_FORBIDDEN`| `HistorySyncForbiddenError`    |
| 409  | `MUTATION_CONFLICT`     | `MutationDuringPaginationError`|
| 422  | `UNPROCESSABLE_CONTENT` | `ApiError`                     |
| 422  | `INVALID_DATE_RANGE`    | `InvalidDateRangeError`        |
| 429  | `RATE_LIMIT_EXCEEDED`   | `RateLimitExceededError`       |
| 503  | (pause mode)            | `ServiceUnavailableError`      |

### Other headers
- `X-Request-Id` on every response — SDKs should attach it to thrown errors and accept a logger callback.
- `Content-Type: application/json` required on POST/PUT/PATCH (415 otherwise).

### Polling model (v1/v2)
Architecture is poll-based, not push. Document this in each SDK's README and provide example code mirroring the upstream OpenAPI's "Implementation Checklist". v3 adds webhooks — webhook signature verification helpers belong in the SDK.

## Tests

Each language has its own runner. Common rules:

- Unit-test the transport (mocked HTTP) and each resource client (fixture-based).
- Integration tests hit a recorded cassette / VCR-style replay — never the live API in CI.
- One smoke test per release that hits a sandbox API key against staging (`equipo-staging.tesote.com`); gated behind a secret so PRs from forks skip it.
- Aim for parity in test coverage across languages so a missing test in PHP is as visible as a missing test in TS.

## CI / release

Runners: **Blacksmith 2vcpu** for both test and release jobs (`runs-on: blacksmith-2vcpu-ubuntu-2204`). Per-language workflows under `.github/workflows/`:

- `<lang>-test.yml` — runs on PRs touching `<lang>/**` or shared spec.
- `<lang>-release.yml` — triggered by tag `<lang>-vX.Y.Z`, builds + publishes to the language registry.

Tags are per-language so SDKs version independently. Release secrets per registry stored in repo secrets (`NPM_TOKEN`, `PYPI_TOKEN`, `RUBYGEMS_TOKEN`, `MAVEN_*`, `PACKAGIST_TOKEN`, none for Go — Go publishes on tag push via the proxy).

## Documentation

Usage docs and API reference both live in the marketing/docs site at `../tesote.com` (do **not** duplicate them here — link from each SDK's README). When changing an SDK's public surface, also update the corresponding doc page in that repo in the same PR.

## CI shape (per language)

One workflow file per language under `.github/workflows/<lang>.yml`, three jobs:

1. **`detect`** — `dorny/paths-filter@v3` sets `should_run` true if `packages/<lang>/**` or `spec/**` changed. Other jobs `needs: detect` and skip on false.
2. **`test`** — unit tests (mocked HTTP) + integration tests (recorded cassettes/replay). Matrix on the language's supported version range.
3. **`release`** — `if: startsWith(github.ref, 'refs/tags/<lang>-v')`. Verify tag matches package version, build, publish, GitHub Release.

All jobs `runs-on: blacksmith-2vcpu-ubuntu-2204`.

## Deep dives

Architecture details live under `docs/architecture/`. Read before scaffolding:

- [versioning.md](docs/architecture/versioning.md) — v1/v2/v3 coexistence, back-compat policy
- [transport.md](docs/architecture/transport.md) — retries, caching, rate-limits, idempotency, pagination
- [errors.md](docs/architecture/errors.md) — typed-error taxonomy, "good error" definition
- [resources.md](docs/architecture/resources.md) — endpoint inventory by version
- [auth.md](docs/architecture/auth.md) — bearer-token rules
- [testing.md](docs/architecture/testing.md) — unit / replay / smoke layers
- [release.md](docs/architecture/release.md) — Blacksmith CI + tag-driven releases

End-user README: [`README.md`](README.md). End-user docs: `www.tesote.com/docs/sdk` (sibling `tesote.com` repo).

## Style

- Lead with the rule. Fragments over sentences. Tables for structured data.
- One logical change per commit; review every diff line.
- Files under ~500 LOC; split into smaller helper classes/modules instead.
- No `rescue Exception` / catch-all error handlers — typed errors only.
- No safe-navigation (`&.`, `?.`, `?:`) hiding nil — make nil explicit or refactor it out.
