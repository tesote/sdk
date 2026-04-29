# Changelog

All notable changes to the `com.tesote:sdk` artifact are documented here.
Format: [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

This project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 0.2.0 - 2026-04-28

### Added

- Full v1 + v2 resource layer. All 35 endpoints from the controller spec
  are implemented; no `UnsupportedOperationException` placeholders remain.
- v1 resource clients: `StatusClient` (status, whoami), `AccountsClient`
  (list, get), `TransactionsClient` (listForAccount, get).
- v2 resource clients: `StatusClient`, `AccountsClient` (+ `sync`),
  `TransactionsClient` (listForAccount, get, sync, syncLegacy, bulk,
  search, export), `SyncSessionsClient`, `TransactionOrdersClient`
  (list, get, create, submit, cancel), `BatchesClient` (create, show,
  approve, submit, cancel), `PaymentMethodsClient` (list, get, create,
  update, delete).
- Java records for every API model: `Account`, `Transaction`,
  `SyncTransaction`, `SyncSession`, `TransactionOrder`, `PaymentMethod`,
  `BatchSummary`, `Status`, `Whoami`, plus paginated wrappers
  (`PagePagination`, `CursorPagination`, `OffsetPage<T>`,
  `AccountsPage`, `TransactionsPage`, `SyncSessionsPage`) and request /
  response envelopes.
- New typed exceptions, one per API `error_code`: `AccountNotFoundException`,
  `TransactionNotFoundException`, `SyncSessionNotFoundException`,
  `PaymentMethodNotFoundException`, `TransactionOrderNotFoundException`,
  `BatchNotFoundException`, `BankConnectionNotFoundException`,
  `CategoryNotFoundException`, `CounterpartyNotFoundException`,
  `LegalEntityNotFoundException`, `WebhookNotFoundException`
  (all subclasses of new `NotFoundException`); `ValidationException`,
  `BatchValidationException`, `BankSubmissionException`;
  `InvalidCursorException`, `InvalidCountException`,
  `InvalidLimitException`, `InvalidQueryException`,
  `MissingDateRangeException`; `InvalidOrderStateException`,
  `SyncInProgressException`, `SyncRateLimitExceededException`,
  `BankUnderMaintenanceException`, `InternalErrorException`.
- `Transport.requestRaw(...)` for file-download endpoints (CSV / JSON
  export); `Transport.Options.jsonBody(Object)` and `query(Map)` helpers.
- Unit tests for every resource client (mocked via `MockWebServer`) plus
  expanded error-dispatcher coverage.

### Changed

- `Content-Type: application/json` is now sent on every POST/PUT/PATCH,
  even when the body is empty, matching the spec's 415 contract.
- `V1Client` / `V2Client` accessors (`accounts()`, `transactions()`, etc.)
  now return live resource clients instead of `UnsupportedOperationException`.

## 0.1.1 - 2026-04-28

### Changed

- Replaced `2.+` / `5.+` / `4.+` / `1.+` Gradle dynamic version notation
  with explicit Maven ranges (`[2.18,3)` etc). Sonatype Central Portal
  validation rejects POMs containing `+` in dependency versions, which
  failed the 0.1.0 deployment. Same dependency intent, Central-compatible
  POM.

## 0.1.0 - 2026-04-28 *(rejected by Central — see 0.1.1)*

Initial bootstrap.

### Added

- Gradle Kotlin DSL build with `java-library`, `maven-publish`, `signing`,
  and `com.gradleup.nmcp` (Sonatype Central Portal).
- Java 17 toolchain, tested on 17 and 21.
- `com.tesote.sdk.Transport` on `java.net.http.HttpClient` with:
  - Bearer token injection
  - Retries (3 attempts, exponential backoff with jitter, 250ms base, 8s cap)
  - Rate-limit header capture via `lastRateLimit()`
  - Idempotency-Key auto-generation (UUIDv4) for mutating methods
  - Request-id propagation into thrown exceptions
  - Bearer-token redaction to `Bearer ****<last4>`
  - Opt-in TTL response cache via pluggable `CacheBackend` interface
- Full typed error hierarchy under `com.tesote.sdk.errors`:
  `TesoteException`, `ApiException`, `UnauthorizedException`,
  `ApiKeyRevokedException`, `WorkspaceSuspendedException`,
  `AccountDisabledException`, `HistorySyncForbiddenException`,
  `MutationDuringPaginationException`, `UnprocessableContentException`,
  `InvalidDateRangeException`, `RateLimitExceededException`,
  `ServiceUnavailableException`, `TransportException`, `NetworkException`,
  `TimeoutException`, `TlsException`, `ConfigException`,
  `EndpointRemovedException`.
- Versioned client builders: `V1Client`, `V2Client`.
- v2 `accounts().list()` and `accounts().get(id)` wired end-to-end. Other
  resources stub with `UnsupportedOperationException`.
- JUnit 5 + MockWebServer test suite covering transport behaviors and error
  dispatch.

### Removed

- v3 client surface (will return as a separate release once the upstream OpenAPI is finalized).

### Notes

- Single runtime dependency: `jackson-databind`. Justified in README.md.
- No Java 21 features used in the core API.
