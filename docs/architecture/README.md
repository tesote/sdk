# Architecture

These docs define how every `tesote-sdk` (TS, Python, Ruby, Java, PHP, Go) is structured. Read in order:

1. [versioning.md](versioning.md) — how v1/v2/v3 clients coexist; back-compat policy
2. [transport.md](transport.md) — HTTP layer: retries, caching, rate-limits, idempotency, pagination
3. [errors.md](errors.md) — typed-error taxonomy; what "raise a good error" means
4. [resources.md](resources.md) — endpoint inventory across API versions
5. [auth.md](auth.md) — bearer-token auth and key-type rules
6. [testing.md](testing.md) — unit / replay / smoke layers; cross-language parity
7. [release.md](release.md) — Blacksmith 2vcpu CI + per-language tag-driven releases

## Non-negotiables

- **All six languages keep public-API parity.** A method existing in TS must exist in Go, named idiomatically. A breaking change in one is a breaking change in all.
- **All API versions stay shipped forever.** v1 and v2 do not get removed when v3 lands.
- **Transport owns cross-cutting concerns.** Resource clients are thin marshal/unmarshal layers. Caching, retry, rate-limit, idempotency, pagination, request-id propagation — all in transport.
- **Errors are typed and rich.** One class per `error_code`; every error carries `request_id`, `http_status`, `retry_after`, `response_body`, redacted `request_summary`.
- **Only client-facing endpoints.** Anything that requires session-cookie auth, admin scope, or internal IPs stays out of the SDK.

## Repo layout

```
sdk/
├── ts/        python/    ruby/    java/    php/    go/
├── docs/
│   └── architecture/   ← this directory
├── spec/                ← vendored OpenAPI snapshots (v1, v2; v3 derived)
└── .github/workflows/   ← per-language test + release pipelines
```

## What is not in this repo

- API server source — lives in the platform repo (sibling dir).
- End-user docs — live at `www.tesote.com/docs/sdk` (sibling `tesote.com` repo). SDK READMEs link out; do not duplicate.
- Internal admin tooling, MCP server transport details beyond the public `/v3/mcp` endpoint, webhook delivery infra.
