# Changelog

All notable changes to `tesote/sdk` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
adheres to semver per [docs/architecture/release.md](../../docs/architecture/release.md).

## 0.2.0 - 2026-04-28

### Added
- Full v1 + v2 resource surface (35 endpoints): `V1\Status`, `V1\Accounts`
  (list/get/listTransactions), `V1\Transactions` (get); `V2\Status`,
  `V2\Accounts` (list/get/sync/listTransactions/syncTransactions/
  exportTransactions), `V2\Transactions` (get/sync/bulk/search),
  `V2\SyncSessions` (listForAccount/get), `V2\TransactionOrders`
  (listForAccount/get/create/submit/cancel), `V2\Batches`
  (create/get/approve/submit/cancel), `V2\PaymentMethods`
  (list/get/create/update/delete).
- Typed readonly model classes under `Tesote\Sdk\Models\` for every
  payload in the v1/v2 spec.
- New typed exception classes covering every remaining `error_code`
  (account/transaction/payment-method/order/batch not-found,
  invalid-cursor/count/limit/query, missing-date-range, sync-in-progress,
  sync-rate-limit-exceeded, bank-under-maintenance, bank-connection-not-
  found, validation, invalid-order-state, bank-submission-error,
  batch-validation-error, internal-error).
- `Transport::requestRaw()` for endpoints returning non-JSON bodies
  (transactions export).
- PHPUnit coverage: one test file per resource, plus the existing
  Transport / Errors suites.

### Removed
- `Tesote\Sdk\NotImplemented` — every resource is now wired.

## 0.1.1 - 2026-04-28

### Changed

- Dual-tag releases (`php-vX.Y.Z` + `vX.Y.Z`) so Packagist parses the
  version and indexes it alongside the cross-language tag.

## 0.1.0 - 2026-04-28

### Added
- Initial scaffold of the PHP SDK.
- `Tesote\Sdk\Transport` built on ext-curl: bearer injection, exponential
  backoff with jitter retries, rate-limit header capture, idempotency-key
  auto-generation for mutations, request-id propagation, opt-in TTL cache
  via `CacheBackend`.
- Typed exception hierarchy under `Tesote\Sdk\Errors\` covering every
  `error_code` from the API plus transport-level (`NetworkException`,
  `TimeoutException`, `TlsException`).
- `V1\Client`, `V2\Client` exposing per-resource sub-clients.
- `accounts.list()` and `accounts.get()` wired end-to-end on both
  versions; remaining resources stubbed via `NotImplemented`.
- PHPUnit, PHPStan (level 8) and php-cs-fixer dev tooling.

### Removed
- v3 client surface (will return as a separate release once the upstream OpenAPI is finalized).
