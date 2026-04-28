# Errors

A developer who catches an SDK error has everything to debug without re-running the request.

## API error envelope

```json
{
  "error": "Human-readable message",
  "error_code": "MACHINE_CODE",
  "error_id": "uuid (server-side log correlation)",
  "retry_after": 60
}
```

SDK parses this into a typed exception. `error_id`, `retry_after`, and any vendor fields preserved as attributes on the exception.

## Required fields on every error class

| Field             | Source                                          | Purpose |
|-------------------|-------------------------------------------------|---------|
| `errorCode`       | `error_code` from envelope                      | Programmatic dispatch |
| `message`         | `error` from envelope (or synthesized)          | Human display |
| `httpStatus`      | HTTP response status                            | Tier triage (4xx vs 5xx) |
| `requestId`       | `X-Request-Id` response header                  | Support ticket correlation |
| `errorId`         | `error_id` from envelope                        | Server log correlation |
| `retryAfter`      | `Retry-After` header or envelope                | Backoff hint |
| `responseBody`    | Raw bytes/string                                | Unexpected-shape debugging |
| `requestSummary`  | `{ method, path, query (redacted), bodyShape }` | Reproduce without secrets |
| `attempts`        | Retry count when raised                         | Distinguish transient vs persistent |

Bearer token **never** in `requestSummary` — redact to `Bearer <last4>` when serializing.

## Class hierarchy

```
TesoteError                       (base; catch-all only as last resort)
├── ApiError                      (server-returned, typed below)
│   ├── UnauthorizedError                  (401, UNAUTHORIZED)
│   ├── ApiKeyRevokedError                 (401, API_KEY_REVOKED)
│   ├── WorkspaceSuspendedError            (403, WORKSPACE_SUSPENDED)
│   ├── AccountDisabledError               (403, ACCOUNT_DISABLED)
│   ├── HistorySyncForbiddenError          (403, HISTORY_SYNC_FORBIDDEN)
│   ├── MutationDuringPaginationError      (409, MUTATION_CONFLICT)
│   ├── UnprocessableContentError          (422, UNPROCESSABLE_CONTENT)
│   ├── InvalidDateRangeError              (422, INVALID_DATE_RANGE)
│   ├── RateLimitExceededError             (429, RATE_LIMIT_EXCEEDED)
│   └── ServiceUnavailableError            (503, pause mode)
├── TransportError               (no usable HTTP response)
│   ├── NetworkError             (DNS, connection refused, reset)
│   ├── TimeoutError             (connect or read timeout)
│   └── TlsError                 (certificate / handshake failures)
├── ConfigError                  (bad SDK config; raised at construction)
└── EndpointRemovedError         (calling a method whose upstream endpoint is gone in this version)
```

## Naming across languages

| Language   | Convention |
|------------|------------|
| TypeScript | `RateLimitExceededError extends TesoteError` |
| Python     | `RateLimitExceededError(TesoteError)` |
| Ruby       | `TesoteSdk::RateLimitExceededError < TesoteSdk::Error` |
| Java       | `RateLimitExceededException extends TesoteException` |
| PHP        | `RateLimitExceededException extends TesoteException` |
| Go         | sentinel + typed: `ErrRateLimitExceeded` and `*RateLimitExceededError` implementing `error`; use `errors.As` |

Class names mirror so docs and stack traces stay searchable.

## "Good error" in practice

Bad: `"422 Unprocessable Entity"`.

Good:

```
RateLimitExceededError: 429 Too Many Requests
  error_code: RATE_LIMIT_EXCEEDED
  request_id: 7f3d2c11-...
  retry_after: 42s
  attempts: 4
  request: POST /api/v3/transactions/bulk?account_id=acct_... (body: 47 items)
  response: { "error": "Rate limit exceeded", "error_code": "RATE_LIMIT_EXCEEDED", "retry_after": 42 }
```

First line is greppable. Rest is everything to file a support ticket or reproduce the call.

## Error-handling rules for SDK code

- Never catch the language's base exception (`Exception`, `Throwable`, `error`). Catch the narrowest type.
- Never swallow an error to "make a method nicer." Failed request → caller hears about it.
- Transport-level retries are the only place errors are caught and re-raised; everywhere else, let typed errors propagate.
- Wrapping a lower-level exception (e.g. HTTP-library `ConnectionError`): preserve as `cause` / `__cause__` / `Unwrap()`. Never lose the chain.
