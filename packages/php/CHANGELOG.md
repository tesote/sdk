# Changelog

All notable changes to `tesote/sdk` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
adheres to semver per [docs/architecture/release.md](../../docs/architecture/release.md).

## 0.1.0 - unreleased

### Added
- Initial scaffold of the PHP SDK.
- `Tesote\Sdk\Transport` built on ext-curl: bearer injection, exponential
  backoff with jitter retries, rate-limit header capture, idempotency-key
  auto-generation for mutations, request-id propagation, opt-in TTL cache
  via `CacheBackend`.
- Typed exception hierarchy under `Tesote\Sdk\Errors\` covering every
  `error_code` from the API plus transport-level (`NetworkException`,
  `TimeoutException`, `TlsException`).
- `V1\Client`, `V2\Client`, `V3\Client` exposing per-resource sub-clients.
- `accounts.list()` and `accounts.get()` wired end-to-end on all three
  versions; remaining resources stubbed via `NotImplemented`.
- PHPUnit, PHPStan (level 8) and php-cs-fixer dev tooling.
