# Changelog

All notable changes to `tesote-sdk` (Ruby) are documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versioning is
[SemVer](https://semver.org/) per the SDK's back-compat policy.

## 0.1.0 - 2026-04-28

### Added

- Initial Ruby SDK scaffold.
- Versioned clients: `TesoteSdk::V1::Client`, `V2::Client`.
- Stdlib-only HTTP transport (`Net::HTTP`) with bearer auth, retries
  (exp-backoff + jitter, 3 attempts default), rate-limit header capture
  (`client.last_rate_limit`), auto-generated `Idempotency-Key` for mutations,
  request-id propagation, opt-in TTL LRU cache via `CacheBackend`,
  bearer-token redaction in logs.
- Typed error hierarchy mapping every documented `error_code` to its own
  subclass of `TesoteSdk::ApiError`; transport-level `NetworkError`,
  `TimeoutError`, `TlsError`.
- Wired endpoints: `accounts.list` and `accounts.get` on every version. All
  other resource methods raise `NotImplementedError` until subsequent releases.

### Removed

- v3 client surface (will return as a separate release once the upstream OpenAPI is finalized).
