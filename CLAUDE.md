# CLAUDE.md

Public monorepo for the official client SDKs of the **equipo.tesote.com** API. Greenfield — only `.idea/` exists today; the shape below is the contract for scaffolding.

Repo name: `tesote-sdk`. npm uses scoped `@tesote/sdk`; other registries use `tesote-sdk`.

| Language   | Folder             | Package name               | Min version | Registry         |
|------------|--------------------|----------------------------|-------------|------------------|
| TypeScript | `packages/ts/`     | `@tesote/sdk`              | Node 18     | npm              |
| Python     | `packages/python/` | `tesote-sdk`               | Python 3.9  | PyPI             |
| Ruby       | `packages/ruby/`   | `tesote-sdk`               | Ruby 3.0    | RubyGems         |
| Java       | `packages/java/`   | `com.tesote:sdk`           | Java 17     | Maven Central    |
| PHP        | `packages/php/`    | `tesote/sdk`               | PHP 8.1     | Packagist        |
| Go         | `packages/go/`     | `github.com/tesote/sdk/go` | Go 1.21     | proxy.golang.org |

## Versions & deps

- **Min runtime** (table above): conservative floor. No features/syntax younger than the floor. No experimental.
- **Runtime deps: zero.** Stdlib only for HTTP, JSON, retries, caching. TS `fetch`, Python `urllib.request`+`json`, Ruby `Net::HTTP`+`json`, Java `java.net.http.HttpClient` (`jackson-databind` allowed if `jakarta.json` too awkward — only acceptable runtime dep), PHP ext-curl+`json_*`, Go `net/http`+`encoding/json`.
- **Dev/test/build deps**: latest stable, loose pins (`^x.y`, `~> x.y` — never `=x.y.z`). Actions: latest major (`actions/checkout@v4`). Test matrix: floor + latest LTS + current stable.
- Each language independently versioned, released, tested. No cross-language code sharing — duplicate idiomatically.
- **Semver**: patch is per-language only; minor and major land across all six in lockstep, gated by `parity-check.yml`.
- **Initial releases ship as `0.1.x`.** Pre-1.0 — surface may evolve from early-adopter feedback. `1.0.0` lands when the v3 OpenAPI is finalized upstream and 0.1.x has shipped stable for one cycle.

## API source of truth

Upstream OpenAPI lives in the platform repo (sibling dir, vendor a snapshot — never import at build time):

- v1: `../<platform>/engines/tesote_api/docs/openapi.yaml` — read-only accounts + transactions
- v2: `../<platform>/engines/tesote_api/docs/openapi_v2.yaml` — adds sync sessions, transaction orders, batches, payment methods, bulk, search
- v3: routes only (`engines/tesote_api/config/routes.rb`, controllers under `app/controllers/tesote_api/v3/`) — adds categories, counterparties, legal entities, connections, webhooks, reports, balance history, workspace, MCP. **No OpenAPI doc yet** — derive from controllers/serializers; flag missing pieces.

Expose only **client-facing, API-key-authenticated** endpoints (mounted under `TesoteApi::Engine` with `current_api_key` auth). No internal admin/session-cookie controllers.

## SDK shape

Versioned clients side-by-side; consumer picks (`V1Client`, `V2Client`, `V3Client`). Per-language signatures in [versioning.md](docs/architecture/versioning.md). **Back-compat is permanent.** v1 and v2 stay shipped after v3 lands. Removing or renaming a public symbol in any version = breaking. Don't.

One module/file per resource (accounts, transactions, sync_sessions, transaction_orders, batches, payment_methods, categories, counterparties, legal_entities, connections, webhooks, reports, balance_history, workspace). SOLID/SRP:

- Transport layer separate from resource clients (one HTTP client, swappable for tests).
- Resource clients thin: marshal params → call transport → deserialize **into typed model objects, not raw maps/hashes/dicts**. Per language: TS classes/interfaces, Python `@dataclass`, Ruby PORO classes (or `Struct`), Java records, PHP readonly classes with typed properties, Go structs. Field names match the API casing in the docs but follow each language's idiomatic casing in the public model (snake_case preserved on the wire, camelCase/PascalCase on the model where idiomatic).
- Errors typed (one class per `error_code`).
- **Transport owns cross-cutting**: pagination, retry (exponential backoff + jitter), rate-limit awareness, response caching (ETag / `Cache-Control` / opt-in TTL), idempotency keys for mutations, request-id propagation. Resource clients never reimplement.

### Error payload (every error class)

`error_code` · `message` · `http_status` · `request_id` (from `X-Request-Id`) · `retry_after` · `response_body` · `request_summary` (method + path + redacted query/body — never the bearer token).

One class per `error_code` (full table in [errors.md](docs/architecture/errors.md)). Don't collapse into a single `ApiError`. Transport-level failures get typed classes (`NetworkError`, `TimeoutError`) — never bubble up the underlying language exception.

## API contract

- **Auth**: `Authorization: Bearer <api_key>`. No other schemes.
- **Rate limits**: 200 req/min per API key, 400 req/min per IP. Headers: `X-RateLimit-{Limit,Remaining,Reset}`. On 429: `Retry-After`. Retry with backoff to a configurable cap → `RateLimitExceededError`.
- **Error envelope**: `{ "error": "...", "error_code": "...", "error_id": "uuid?", "retry_after": 60 }`. Map every `error_code` to a typed exception.
- **`X-Request-Id`** on every response — attach to thrown errors; accept a logger callback.
- **`Content-Type: application/json`** required on POST/PUT/PATCH (415 otherwise).
- **Polling**: v1/v2 are poll-based, not push. Document in each SDK's README with example code mirroring the OpenAPI's "Implementation Checklist". v3 adds webhooks — signature-verification helpers ship in the SDK.

## Tests

- Unit: mocked HTTP per resource client; full transport coverage.
- Integration: recorded cassette / VCR-style replay. Never the live API in CI.
- One smoke test per release against staging (`equipo-staging.tesote.com`); gated behind a secret so fork PRs skip.
- Cross-language coverage parity — missing test in PHP is as visible as missing test in TS.

## CI / release

Runners: **Blacksmith 2vcpu** (`runs-on: blacksmith-2vcpu-ubuntu-2204`). One workflow per language at `.github/workflows/<lang>.yml`, three jobs:

1. **`detect`** — `dorny/paths-filter@v3` sets `should_run` if `packages/<lang>/**` or `spec/**` changed. Other jobs `needs: detect`.
2. **`test`** — unit + integration replay. Matrix across supported versions.
3. **`release`** — `if: startsWith(github.ref, 'refs/tags/<lang>-v')`. Verify tag = manifest version, build, publish, GitHub Release.

Tags per-language so SDKs version independently. OIDC trusted publishers for npm/RubyGems/PyPI; Sonatype Central Portal user token for Maven; `PACKAGIST_TOKEN` for PHP; Go publishes via tag push (no token). See [release.md](docs/architecture/release.md).

## Documentation

End-user docs + API reference live at `../tesote.com` (`www.tesote.com/docs/sdk`) — link from each SDK's README; do not duplicate. Public-surface PRs update the matching doc page in that repo, same PR. End-user README: [`README.md`](README.md).

## Deep dives — `docs/architecture/`

- [versioning.md](docs/architecture/versioning.md) — v1/v2/v3 coexistence, back-compat
- [transport.md](docs/architecture/transport.md) — retries, caching, rate-limits, idempotency, pagination
- [errors.md](docs/architecture/errors.md) — typed-error taxonomy, "good error" definition
- [resources.md](docs/architecture/resources.md) — endpoint inventory by version
- [auth.md](docs/architecture/auth.md) — bearer-token rules
- [testing.md](docs/architecture/testing.md) — unit / replay / smoke layers
- [release.md](docs/architecture/release.md) — Blacksmith CI, tag-driven releases, OIDC

## Style

- Lead with the rule. Fragments over sentences. Tables for structured data.
- One logical change per commit; review every diff line.
- Files under ~500 LOC; split into smaller modules.
- No `rescue Exception` / catch-all handlers — typed errors only.
- No safe-navigation (`&.`, `?.`, `?:`) hiding nil — make nil explicit or refactor it out.
