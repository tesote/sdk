# Architecture

How every `tesote-sdk` (TS, Python, Ruby, Java, PHP, Go) is structured. Read in order:

1. [versioning.md](versioning.md) — v1/v2 coexistence; back-compat policy
2. [transport.md](transport.md) — HTTP layer: retries, caching, rate-limits, idempotency, pagination
3. [errors.md](errors.md) — typed-error taxonomy; "good error" definition
4. [resources.md](resources.md) — endpoint inventory across API versions
5. [auth.md](auth.md) — bearer-token auth, key-type rules
6. [testing.md](testing.md) — unit / replay / smoke layers; cross-language parity
7. [release.md](release.md) — Blacksmith 2vcpu CI, per-language tag-driven releases, OIDC

## Non-negotiables

- All six languages keep public-API parity. Method in TS → must exist in Go, named idiomatically. Breaking change in one = breaking change in all.
- All API versions stay shipped forever. v1 does not get removed when later versions land.
- Transport owns cross-cutting concerns. Resource clients are thin marshal/unmarshal. Caching, retry, rate-limit, idempotency, pagination, request-id propagation — all in transport.
- Errors are typed and rich. One class per `error_code`; every error carries `request_id`, `http_status`, `retry_after`, `response_body`, redacted `request_summary`.
- Only client-facing endpoints. Anything requiring session-cookie auth, admin scope, or internal IPs stays out.

## Repo layout

```
sdk/
├── packages/
│   ├── ts/    python/    ruby/    java/    php/    go/
├── docs/
│   └── architecture/   ← this directory
├── spec/                ← vendored OpenAPI snapshots (v1, v2)
└── .github/workflows/   ← per-language test + release pipelines
```

## Not in this repo

- API server source — platform repo (sibling dir).
- End-user docs — `www.tesote.com/docs/sdk` (sibling `tesote.com` repo). SDK READMEs link out.
- Internal admin tooling, webhook delivery infra.
