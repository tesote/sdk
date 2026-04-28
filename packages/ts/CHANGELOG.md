# Changelog

All notable changes to `@tesote/sdk` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
adheres to semver.

## 0.1.0 - unreleased

### Added

- Initial scaffold: `V1Client`, `V2Client`.
- Native-`fetch` transport with bearer auth, exponential-backoff retries,
  rate-limit header capture (`client.lastRateLimit`), idempotency-key
  auto-generation for mutations, request-id propagation into thrown errors,
  opt-in TTL response cache (in-memory LRU; `CacheBackend` interface).
- Full typed-error hierarchy mirroring `docs/architecture/errors.md`.
- `accounts.list()` and `accounts.get(id)` wired end-to-end on every API
  version. Other methods stubbed with signatures.

### Removed

- v3 client surface (will return as a separate release once the upstream
  OpenAPI is finalized).
