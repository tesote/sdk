# Changelog

All notable changes to `tesote-sdk` (Python) are listed here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project adheres to semver per the monorepo's `docs/architecture/versioning.md`.

## 0.2.0 - 2026-04-28

### Added

- Full v1 + v2 endpoint surface (35 endpoints). v1: `accounts.list/get`,
  `transactions.list_for_account/get`, `status.status/whoami`. v2:
  `accounts.list/get/sync`, `transactions.list_for_account/get/sync/sync_legacy/bulk/search/export`,
  `sync_sessions.list/get`, `transaction_orders.list/get/create/submit/cancel`,
  `batches.create/show/approve/submit/cancel`,
  `payment_methods.list/get/create/update/delete`, `status.status/whoami`.
- Frozen `@dataclass` models for every payload in `tesote_sdk.models`:
  `Account`, `Transaction`, `SyncTransaction`, `SyncDelta`, `SyncResult`,
  `SyncSession`, `TransactionOrder` (+ attempt / destination / source / fee),
  `PaymentMethod`, `BatchSummary`, `BatchCreateResult`, `BatchApproveResult`,
  `BatchSubmitResult`, `BatchCancelResult`, `AccountSyncStarted`,
  `StatusResponse`, `WhoAmI`, plus list/cursor wrappers.
- 20+ new typed exceptions mapped from `error_code`:
  `AccountNotFoundError`, `TransactionNotFoundError`, `SyncSessionNotFoundError`,
  `PaymentMethodNotFoundError`, `TransactionOrderNotFoundError`,
  `BatchNotFoundError`, `BankConnectionNotFoundError`, `SyncInProgressError`,
  `SyncRateLimitExceededError`, `BankUnderMaintenanceError`, `ValidationError`,
  `BatchValidationError`, `InvalidOrderStateError`, `InvalidCursorError`,
  `InvalidCountError`, `InvalidLimitError`, `InvalidQueryError`,
  `MissingDateRangeError`, `BankSubmissionError`, `InternalError`,
  `NotFoundError`.
- Cursor- and offset-pagination iterators (`iter_*` generators) on every
  list method.

### Removed

- `tesote_sdk/v2/_stubs.py` — every method is now wired end-to-end.

## 0.1.0 - 2026-04-28

Initial scaffold.

### Added

- `V1Client`, `V2Client` — versioned clients, side-by-side per the monorepo's versioning policy.
- Single `Transport` on stdlib `urllib.request` + `json`. Zero runtime dependencies.
  - Bearer-token auth, automatic injection.
  - Retries with exponential backoff + full jitter; defaults: 3 attempts, base 250ms, cap 8s. Configurable via `RetryPolicy`. Retries on 429 / 502 / 503 / 504 / network errors. Never retries 4xx other than 429, never retries non-idempotent timeouts without an idempotency key.
  - Auto-generated `Idempotency-Key: <uuid4>` on POST/PUT/PATCH/DELETE.
  - Rate-limit header capture exposed as `client.last_rate_limit`.
  - Request-id propagation into every typed error.
  - Opt-in TTL LRU cache via `CacheBackend` Protocol; `InMemoryLRUCache` shipped.
  - Logger callback hook with bearer-token redaction (`Bearer <last4>`).
  - Configurable `connect_timeout` (5s) and `read_timeout` (30s).
  - Default base URL `https://equipo.tesote.com/api`.
- Full typed error hierarchy in `tesote_sdk.errors`. One class per `error_code`. Required fields: `error_code`, `message`, `http_status`, `request_id`, `error_id`, `retry_after`, `response_body`, `request_summary`, `attempts`. `__cause__` preserved.
- v1: `accounts.list`, `accounts.get` wired end-to-end.
- v2: `accounts.list`, `accounts.get` wired end-to-end. Other resources stubbed.
- Test suite covering transport (mocked `urllib.request.urlopen`) and error mapping.
- `pyproject.toml` with `hatchling` build backend, optional `[test]` extras (`pytest`, `mypy`, `ruff`), no version pins.

### Removed

- v3 client surface (will return as a separate release once the upstream OpenAPI is finalized).
