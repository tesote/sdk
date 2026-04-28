# Changelog

All notable changes to this SDK are documented here. Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Versioning: [SemVer](https://semver.org/spec/v2.0.0.html).

## 0.1.0 - unreleased

### Added

- Initial Go SDK scaffold under `packages/go/`.
- `tesote.Client` transport with bearer auth, exp-backoff retries (3 attempts, 250ms base, 8s cap), rate-limit header capture (`LastRateLimit()`), auto-generated UUIDv4 idempotency keys for mutations, request-id propagation, opt-in TTL LRU cache via `CacheBackend`, bearer redaction.
- Typed error hierarchy mirroring `docs/architecture/errors.md`: sentinels (`ErrUnauthorized`, `ErrRateLimitExceeded`, etc.) plus rich `*APIError` subtypes (`*RateLimitExceededError`, `*WorkspaceSuspendedError`, ...) and transport errors (`*NetworkError`, `*TimeoutError`, `*TLSError`).
- `v1`, `v2`, `v3` packages with per-resource service stubs.
- `v3.AccountsService.List` and `v3.AccountsService.Get` wired end-to-end.
- Test coverage for transport (httptest), errors (every `error_code` -> typed mapping), and v3 accounts.
