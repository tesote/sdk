# Changelog

All notable changes to `@tesote.com/sdk` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
adheres to semver.

## 0.2.0 - 2026-04-28

### Added

- Full v1+v2 resource surface (35 endpoints):
  - v1: `accounts.list/get/listAll`, `transactions.listForAccount/get/listAllForAccount`,
    `status.status/whoami`.
  - v2: `accounts.list/get/sync/listAll`, `transactions.listForAccount/get/export/sync/syncLegacy/bulk/search/listAllForAccount`,
    `syncSessions.list/get/listAll`, `transactionOrders.list/get/create/submit/cancel/listAll`,
    `batches.create/get/approve/submit/cancel`,
    `paymentMethods.list/get/create/update/delete/listAll`, `status.status/whoami`.
- Typed model interfaces for every payload: `Account`, `Transaction`, `SyncTransaction`,
  `SyncResult`, `SyncSession`, `TransactionOrder`, `PaymentMethod`, `BatchSummary`,
  `BulkResult`, `SearchResult`, plus pagination envelopes and request inputs.
- Cursor- and offset-pagination async iterators (`listAll*`).
- Typed errors for every API `error_code`: `AccountNotFoundError`, `TransactionNotFoundError`,
  `SyncSessionNotFoundError`, `PaymentMethodNotFoundError`, `TransactionOrderNotFoundError`,
  `BatchNotFoundError`, `BankConnectionNotFoundError`, `InvalidOrderStateError`,
  `SyncInProgressError`, `InvalidCursorError`, `InvalidCountError`, `InvalidLimitError`,
  `InvalidQueryError`, `MissingDateRangeError`, `BankSubmissionError`, `ValidationError`,
  `BatchValidationError`, `SyncRateLimitExceededError`, `BankUnderMaintenanceError`,
  `InternalServerError`, plus a generic `NotFoundError` base.

## 0.1.1 - 2026-04-28

### Changed

- First release published via GitHub Actions OIDC trusted publisher (no
  manual `npm publish`). 0.1.0 was published locally as the bootstrap.

## 0.1.0 - 2026-04-28

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
