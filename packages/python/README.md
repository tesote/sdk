# tesote-sdk

Official Python SDK for the [equipo.tesote.com](https://equipo.tesote.com) API.

- Zero runtime dependencies. Uses only the Python standard library.
- Min runtime: **Python 3.9**.
- Versioned clients side-by-side: `V1Client`, `V2Client`.

## Install

```bash
pip install tesote-sdk
```

## Usage

```python
from tesote_sdk import V2Client

client = V2Client(api_key="sk_live_...")

for account in client.accounts.list():
    print(account["id"])
```

Mix versions in the same process when you need to:

```python
from tesote_sdk import V1Client, V2Client

legacy = V1Client(api_key="sk_live_...")
new = V2Client(api_key="sk_live_...")
```

## Auth

Single scheme: bearer API key. Pass it at construction. The SDK never persists the key, never logs it, and redacts it to `Bearer <last4>` in logs and error summaries.

`V*Client(api_key="")` raises `ConfigError` synchronously.

## Errors

Every error inherits from `TesoteError` and carries: `error_code`, `message`, `http_status`, `request_id`, `error_id`, `retry_after`, `response_body`, `request_summary`, `attempts`. `__cause__` is preserved when wrapping a lower-level exception.

Catch the narrowest type:

```python
from tesote_sdk import RateLimitExceededError, V2Client

try:
    V2Client(api_key=key).accounts.list()
except RateLimitExceededError as e:
    print(f"slow down for {e.retry_after}s (req {e.request_id})")
```

Full hierarchy: see `docs/architecture/errors.md` in the monorepo.

## Transport features

Configured at the client; resource modules never reimplement them:

| Concern | Default |
|---|---|
| Retries | 3 attempts, exp backoff + jitter, base 250ms, cap 8s, retry on 429/502/503/504 + network errors |
| Timeouts | connect 5s, read 30s |
| Idempotency | auto `Idempotency-Key` UUIDv4 on POST/PUT/PATCH/DELETE |
| Rate limits | `client.last_rate_limit` after every request |
| Cache | opt-in TTL LRU via `cache_ttl=` per call; pluggable `CacheBackend` |
| Logging | optional callback, `Authorization` always redacted |

## Polling model (v1, v2)

The platform is poll-based. Use `accounts.sync(...)` (v2) to trigger a refresh, then poll `accounts.get(...)` until the data settles.

## Versioning

- API versions (`v1`, `v2`) ship side-by-side and never get removed.
- SDK semver is independent. Tag releases as `python-vX.Y.Z`.

See `CHANGELOG.md`.

## Development

```bash
cd packages/python
pip install -e .[test]
ruff check .
mypy src
pytest
```

## License

MIT.
