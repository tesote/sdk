# Changelog

All notable changes to the `Tesote.Sdk` NuGet package are documented here.

This project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
