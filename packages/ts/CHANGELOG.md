# Changelog

All notable changes to `@tesote/sdk` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
adheres to semver.

## 0.1.0 - unreleased

### Added

- Initial scaffold: `V1Client`, `V2Client`, `V3Client`.
- Native-`fetch` transport with bearer auth, exponential-backoff retries,
  rate-limit header capture (`client.lastRateLimit`), idempotency-key
  auto-generation for mutations, request-id propagation into thrown errors,
  opt-in TTL response cache (in-memory LRU; `CacheBackend` interface).
- Full typed-error hierarchy mirroring `docs/architecture/errors.md`.
- `accounts.list()` and `accounts.get(id)` wired end-to-end on every API
  version. Other methods stubbed with signatures.
- `verifyWebhookSignature` exported from `@tesote/sdk` and `@tesote/sdk/v3`
  (stub pending platform spec).
