# Changelog

All notable changes to the `Tesote.Sdk` NuGet package are documented here.

This project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 0.2.0 - 2026-04-28

Full v1 + v2 resource surface.

### Added

- All 35 documented v1 + v2 endpoints implemented, with typed model records
  in `Tesote.Sdk.Models` (Account, Transaction, SyncTransaction, SyncSession,
  TransactionOrder, PaymentMethod, BatchSummary, etc.) — `JsonPropertyName`
  attributes preserve snake_case on the wire.
- Resource clients: `V1.{StatusClient, AccountsClient, TransactionsClient}`
  and `V2.{StatusClient, AccountsClient, TransactionsClient,
  SyncSessionsClient, TransactionOrdersClient, BatchesClient,
  PaymentMethodsClient}` — accessible as properties on `V1Client` / `V2Client`.
- New typed exceptions for every documented `error_code`:
  `AccountNotFoundException`, `TransactionNotFoundException`,
  `SyncSessionNotFoundException`, `PaymentMethodNotFoundException`,
  `TransactionOrderNotFoundException`, `BatchNotFoundException`,
  `BankConnectionNotFoundException`, `InvalidCursorException`,
  `InvalidCountException`, `InvalidLimitException`, `InvalidQueryException`,
  `MissingDateRangeException`, `SyncInProgressException`,
  `SyncRateLimitExceededException`, `BankUnderMaintenanceException`,
  `ValidationException`, `BatchValidationException`,
  `InvalidOrderStateException`, `BankSubmissionException`,
  `InternalServerException`, plus a shared `NotFoundException` base.
- `Transport.RequestRawAsync` for file-download endpoints (CSV / JSON export);
  preserves retries, rate-limit awareness, idempotency, and error mapping.
- xUnit + WireMock.Net coverage per resource: success, typed-error mapping,
  pagination, idempotency, and the 415 case.

### Changed

- `Internal.Json.DefaultOptions` no longer applies a global
  `PropertyNamingPolicy`; per-property `[JsonPropertyName]` markers keep the
  wire format stable and allow PascalCase model surfaces.

## 0.1.0 - 2026-04-28

Initial bootstrap.

### Added

- `Tesote.Sdk.csproj` targeting `net8.0` with `Nullable`, `ImplicitUsings`,
  and `TreatWarningsAsErrors` enabled. Documentation XML generated.
- `Tesote.Sdk.Transport` on `System.Net.Http.HttpClient` with:
  - Bearer token injection
  - Async-first API (every public method returns `Task<T>` and accepts a
    `CancellationToken`)
  - Retries (3 attempts, exponential backoff with jitter, 250ms base, 8s cap)
  - Rate-limit header capture via `LastRateLimit`
  - Idempotency-Key auto-generation (UUIDv4) for mutating methods
  - Request-id propagation into thrown exceptions
  - Bearer-token redaction to `Bearer ****<last4>`
  - Opt-in TTL response cache via pluggable `ICacheBackend` interface
  - `IDisposable` + `IAsyncDisposable`
- Full typed error hierarchy under `Tesote.Sdk.Errors`:
  `TesoteException`, `ApiException`, `UnauthorizedException`,
  `ApiKeyRevokedException`, `WorkspaceSuspendedException`,
  `AccountDisabledException`, `HistorySyncForbiddenException`,
  `MutationDuringPaginationException`, `UnprocessableContentException`,
  `InvalidDateRangeException`, `RateLimitExceededException`,
  `ServiceUnavailableException`, `TransportException`, `NetworkException`,
  `TesoteTimeoutException`, `TlsException`, `ConfigException`,
  `EndpointRemovedException`.
- Versioned clients: `V1Client`, `V2Client`, each exposing `Accounts.ListAsync`
  and `Accounts.GetAsync` wired end-to-end. Other resources stub with
  `NotImplementedException`.
- xUnit + WireMock.Net test suite covering transport behaviors and error
  dispatch.

### Removed

- v3 client surface (will return as a separate release once the upstream OpenAPI is finalized).

### Notes

- **Zero runtime dependencies.** HTTP, JSON, retries, caching, and
  concurrency primitives all use the .NET standard library
  (`System.Net.Http`, `System.Text.Json`, `System.Buffers`).
- No `dynamic`, no Newtonsoft.Json, no third-party HTTP client.
