# Changelog

All notable changes to this SDK are documented here. Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Versioning: [SemVer](https://semver.org/spec/v2.0.0.html).

## 0.2.0 - 2026-04-28

### Added

- Full v1 + v2 endpoint surface (35 endpoints). v1: `accounts.List/Get`,
  `transactions.ListForAccount/Get`, `status.Status/Whoami`. v2: `accounts.List/Get/Sync`,
  `transactions.ListForAccount/Get/Sync/SyncLegacy/Bulk/Search/Export`,
  `sync_sessions.List/Get`, `transaction_orders.List/Get/Create/Submit/Cancel`,
  `batches.Create/Show/Approve/Submit/Cancel`,
  `payment_methods.List/Get/Create/Update/Delete`, `status.Status/Whoami`.
- Typed model structs in `models.go` for every payload: `Account`, `Transaction`,
  `SyncTransaction`, `SyncResult`, `SyncSession`, `TransactionOrder`,
  `PaymentMethod`, `BatchSummary`, plus envelope and pagination types — all
  with `json:"snake_case"` tags.
- 21 new typed errors mapped from `error_code`: `*AccountNotFoundError`,
  `*TransactionNotFoundError`, `*SyncSessionNotFoundError`,
  `*PaymentMethodNotFoundError`, `*TransactionOrderNotFoundError`,
  `*BatchNotFoundError`, `*BankConnectionNotFoundError`, `*InvalidCursorError`,
  `*InvalidCountError`, `*InvalidLimitError`, `*InvalidQueryError`,
  `*MissingDateRangeError`, `*SyncInProgressError`, `*SyncRateLimitExceededError`,
  `*BankUnderMaintenanceError`, `*ValidationError`, `*InvalidOrderStateError`,
  `*BankSubmissionError`, `*BatchValidationError`, `*InternalError`, plus
  `Err*` sentinels for each.
- `Transport.RequestRaw` for non-JSON responses (CSV/JSON export).

### Changed

- All `ErrNotImplemented` returns replaced with real implementations across
  `v1` and `v2` packages.

## 0.1.1 - 2026-04-28

### Changed

- Module relocated from `packages/go/` to repo-root `go/` so
  `proxy.golang.org` can resolve `github.com/tesote/sdk/go`. Submodule tags
  follow the Go-toolchain `go/vX.Y.Z` format. No source-level changes from
  0.1.0; consumers can pin either version.

## 0.1.0 - 2026-04-28

### Added

- Initial Go SDK scaffold under `go/`.
- `tesote.Client` transport with bearer auth, exp-backoff retries (3 attempts, 250ms base, 8s cap), rate-limit header capture (`LastRateLimit()`), auto-generated UUIDv4 idempotency keys for mutations, request-id propagation, opt-in TTL LRU cache via `CacheBackend`, bearer redaction.
- Typed error hierarchy mirroring `docs/architecture/errors.md`: sentinels (`ErrUnauthorized`, `ErrRateLimitExceeded`, etc.) plus rich `*APIError` subtypes (`*RateLimitExceededError`, `*WorkspaceSuspendedError`, ...) and transport errors (`*NetworkError`, `*TimeoutError`, `*TLSError`).
- `v1`, `v2` packages with per-resource service stubs.
- Test coverage for transport (httptest) and errors (every `error_code` -> typed mapping).

### Removed

- v3 client surface (will return as a separate release once the upstream OpenAPI is finalized).
