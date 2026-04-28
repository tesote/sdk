# tesote-sdk-go

Official Go SDK for the [equipo.tesote.com](https://equipo.tesote.com) API.

Zero runtime dependencies. Stdlib only (`net/http`, `encoding/json`, `crypto/rand`).

## Install

```sh
go get github.com/tesote/sdk/go@latest
```

Min Go: **1.21**.

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    tesote "github.com/tesote/sdk/go"
    v2 "github.com/tesote/sdk/go/v2"
)

func main() {
    tc, err := tesote.NewClient(tesote.Options{
        APIKey: os.Getenv("TESOTE_API_KEY"),
    })
    if err != nil {
        log.Fatal(err)
    }
    client := v2.New(tc)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    page, err := client.Accounts.List(ctx, v2.ListOptions{PageSize: 50})
    if err != nil {
        log.Fatal(err)
    }
    for _, a := range page.Data {
        fmt.Println(a.ID, a.Name, a.Currency)
    }
}
```

## Versioned clients

`v1`, `v2` ship side by side. Pick the version your code targets; mix in one process if needed.

```go
import (
    v1 "github.com/tesote/sdk/go/v1"
    v2 "github.com/tesote/sdk/go/v2"
)
```

See `docs/architecture/resources.md` in the parent repo for the per-version endpoint inventory.

## Errors

Every error implements `error` and either `errors.Is` (sentinel match) or `errors.As` (typed extraction):

```go
_, err := client.Accounts.List(ctx, v2.ListOptions{})
if errors.Is(err, tesote.ErrRateLimitExceeded) {
    var rl *tesote.RateLimitExceededError
    errors.As(err, &rl)
    log.Printf("retry after %ds (request_id=%s)", rl.RetryAfter, rl.RequestID)
}
```

Sentinels: `ErrUnauthorized`, `ErrAPIKeyRevoked`, `ErrWorkspaceSuspended`, `ErrAccountDisabled`, `ErrHistorySyncForbidden`, `ErrMutationDuringPagination`, `ErrUnprocessableContent`, `ErrInvalidDateRange`, `ErrRateLimitExceeded`, `ErrServiceUnavailable`, `ErrNetwork`, `ErrTimeout`, `ErrTLS`, `ErrConfig`, `ErrEndpointRemoved`.

Typed errors carry: `ErrorCode`, `Message`, `HTTPStatus`, `RequestID`, `ErrorID`, `RetryAfter`, `ResponseBody`, `RequestSummary` (with the bearer redacted to `Bearer <last4>`), `Attempts`.

## Rate-limit headers

```go
rl, ok := tc.LastRateLimit()
if ok && rl.Remaining < 10 {
    // back off
}
```

## Polling

The v1/v2 API is poll-based. Schedule your own ticker; the SDK does not run background goroutines.

## Caching

Pass a `tesote.CacheBackend` (e.g. the built-in `tesote.NewLRUCache`) on the `Options` struct, then opt in per call:

```go
client.Accounts.List(ctx, v2.ListOptions{CacheTTL: 30 * time.Second})
```

## Major-version subpath rule

Go modules require a `/vN` suffix in the import path for every major version `>= 2`. This is **module versioning**, not API versioning.

| SDK release        | Import path                       |
|--------------------|-----------------------------------|
| v0.x.y / v1.x.y    | `github.com/tesote/sdk/go`        |
| v2.x.y             | `github.com/tesote/sdk/go/v2`     |

This is independent of the API version sub-packages (`v1/`, `v2/`) you import from any module-version of this SDK.

## Docs

End-user docs: [www.tesote.com/docs/go](https://www.tesote.com/docs/go).
