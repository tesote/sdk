# Changelog

All notable changes to the `com.tesote:sdk` artifact are documented here.

This project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 0.1.0 - 2026-04-28

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
