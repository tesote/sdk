# Changelog

All notable changes to `tesote-sdk` (Ruby) are documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versioning is
[SemVer](https://semver.org/) per the SDK's back-compat policy.

## 0.2.0 - 2026-04-28

### Added

- Full v1 + v2 endpoint surface (35 endpoints total). v1 retains 6
  read-only endpoints (status, whoami, accounts, accounts.transactions,
  transactions.get); v2 adds 29 endpoints across accounts, transactions,
  sync_sessions, transaction_orders, batches, and payment_methods.
- Typed PORO/Struct models under `TesoteSdk::Models`: `Account`,
  `Transaction`, `SyncTransaction`, `SyncResult`, `SyncSession`,
  `TransactionOrder`, `PaymentMethod`, `BatchSummary`, `BatchCreateResult`,
  `BulkResult`, `SearchResult`, `OffsetPage`, `Pagination`, `Whoami`,
  `StatusResult`, plus nested types. All are forward-compatible with
  unknown fields.
- Cursor and offset pagination helpers (`Pagination::CursorEnumerator`,
  `Pagination::OffsetEnumerator`) — used by `Accounts#each_transaction_page`,
  `Transactions#each_page_for_account`, and `SyncSessions#each_page`.
- New typed errors: `AccountNotFoundError`, `TransactionNotFoundError`,
  `SyncSessionNotFoundError`, `PaymentMethodNotFoundError`,
  `TransactionOrderNotFoundError`, `BatchNotFoundError`,
  `CategoryNotFoundError`, `CounterpartyNotFoundError`,
  `LegalEntityNotFoundError`, `WebhookNotFoundError`,
  `BankConnectionNotFoundError`, `SyncInProgressError`,
  `InvalidOrderStateError`, `MissingDateRangeError`, `InvalidCursorError`,
  `InvalidCountError`, `InvalidLimitError`, `InvalidQueryError`,
  `ValidationError`, `BatchValidationError`, `BankSubmissionError`,
  `SyncRateLimitExceededError`, `InternalServerError`,
  `BankUnderMaintenanceError`. All registered against their server
  `error_code`.
- Transport: `request_unversioned` for the unversioned `/status` and
  `/whoami` endpoints, and `request_raw` for file-download responses
  (CSV/JSON export) returning a `RawResponse` with body, content-type,
  and content-disposition.

### Changed

- `V1::Accounts#list`, `#get`, and `#list_transactions` now return typed
  `Models::AccountList`, `Models::Account`, and `Models::TransactionList`
  instead of raw hashes. Same for v2.
- The `v2/stubs.rb` placeholder file was removed; each resource now lives
  in its own file (`v2/transactions.rb`, `v2/sync_sessions.rb`,
  `v2/transaction_orders.rb`, `v2/batches.rb`, `v2/payment_methods.rb`,
  `v2/status.rb`).

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
