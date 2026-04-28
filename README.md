<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset=".github/assets/logo-white.svg">
    <img src=".github/assets/logo.svg" alt="Tesote" height="72">
  </picture>
</p>

<h1 align="center">tesote-sdk</h1>

<p align="center">
  <a href="https://www.npmjs.com/package/@tesote/sdk"><img alt="npm" src="https://img.shields.io/npm/v/@tesote/sdk?label=npm&color=cb3837"></a>
  <a href="https://pypi.org/project/tesote-sdk/"><img alt="PyPI" src="https://img.shields.io/pypi/v/tesote-sdk?label=pypi&color=3776ab"></a>
  <a href="https://rubygems.org/gems/tesote-sdk"><img alt="Gem" src="https://img.shields.io/gem/v/tesote-sdk?label=rubygems&color=cc0000"></a>
  <a href="https://central.sonatype.com/artifact/com.tesote/sdk"><img alt="Maven Central" src="https://img.shields.io/maven-central/v/com.tesote/sdk?label=maven&color=c71a36"></a>
  <a href="https://packagist.org/packages/tesote/sdk"><img alt="Packagist" src="https://img.shields.io/packagist/v/tesote/sdk?label=packagist&color=f28d1a"></a>
  <a href="https://pkg.go.dev/github.com/tesote/sdk/go"><img alt="Go" src="https://img.shields.io/github/v/tag/tesote/sdk?filter=go-*&label=go&color=00add8"></a>
  <br/>
  <a href="https://github.com/tesote/sdk/actions/workflows/parity-check.yml"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/tesote/sdk/parity-check.yml?label=CI"></a>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/github/license/tesote/sdk?color=blue"></a>
</p>

Official client SDKs for the [equipo.tesote.com](https://equipo.tesote.com) API. One repo, six languages, identical surface.

| Language | Package | Install |
|----------|---------|---------|
| TypeScript | [`@tesote/sdk`](https://www.npmjs.com/package/@tesote/sdk) | `npm i @tesote/sdk` |
| Python | [`tesote-sdk`](https://pypi.org/project/tesote-sdk/) | `pip install tesote-sdk` |
| Ruby | [`tesote-sdk`](https://rubygems.org/gems/tesote-sdk) | `gem install tesote-sdk` |
| Java | `com.tesote:sdk` | Maven Central |
| PHP | [`tesote/sdk`](https://packagist.org/packages/tesote/sdk) | `composer require tesote/sdk` |
| Go | `github.com/tesote/sdk/go` | `go get github.com/tesote/sdk/go` |

Full docs: **https://www.tesote.com/docs/sdk**

---

## Quick start

```ts
import { V2Client } from '@tesote/sdk'

const tesote = new V2Client({ apiKey: process.env.TESOTE_API_KEY! })
const accounts = await tesote.accounts.list()
```

```python
from tesote_sdk import V2Client

tesote = V2Client(api_key=os.environ["TESOTE_API_KEY"])
for account in tesote.accounts.list_all():
    print(account.id, account.balance)
```

```ruby
require 'tesote_sdk'

tesote = TesoteSdk::V2::Client.new(api_key: ENV.fetch('TESOTE_API_KEY'))
tesote.accounts.list.each { |a| puts a.id }
```

```go
import tesote "github.com/tesote/sdk/go/v2"

c := tesote.New(tesote.Config{APIKey: os.Getenv("TESOTE_API_KEY")})
accounts, _ := c.Accounts.List(ctx, nil)
```

```php
use Tesote\Sdk\V2\Client;

$tesote = new Client(['apiKey' => getenv('TESOTE_API_KEY')]);
$accounts = $tesote->accounts->list();
```

```java
import com.tesote.sdk.v2.V2Client;

var tesote = V2Client.builder().apiKey(System.getenv("TESOTE_API_KEY")).build();
var accounts = tesote.accounts().list();
```

---

## What you get

- **Versioned clients side-by-side** — `V1Client`, `V2Client` from the same import. Pick per call site, mix in one process. Old versions never get removed.
- **Transport-level reliability** — automatic retries with backoff + jitter, rate-limit-aware throttling, opt-in response caching, idempotency keys for mutations.
- **Typed errors with full context** — one class per `error_code`, every error carries `request_id`, `http_status`, `retry_after`, `response_body`. Catch the narrow type you care about; ignore the rest.
- **Cursor pagination** — `list()` for one page, `listAll()` for an iterator. Mutation-mid-iteration surfaces a typed `MutationDuringPaginationError`, not silent corruption.

---

## API versions

| Version | Adds |
|---------|------|
| **v1** | Accounts, transactions (read-only) |
| **v2** | + sync sessions, transaction orders, batches, payment methods, bulk + search |

Both ship from every SDK. Back-compat is permanent.

---

## Auth

```
Authorization: Bearer <api_key>
```

Get a key from your Tesote workspace settings. The SDK never persists it; it lives on the client instance only.

---

## Errors

```ts
import { RateLimitExceededError, WorkspaceSuspendedError } from '@tesote/sdk'

try {
  await tesote.transactions.bulk(items)
} catch (e) {
  if (e instanceof RateLimitExceededError) {
    console.log(`retry in ${e.retryAfter}s; req ${e.requestId}`)
  } else if (e instanceof WorkspaceSuspendedError) {
    // contact support
  } else {
    throw e
  }
}
```

Full error taxonomy: [`docs/architecture/errors.md`](docs/architecture/errors.md).

---

## Development

This repo is a multi-language monorepo. Each language is independently testable and releasable.

| Task | Command (per-language dir) |
|------|---------------------------|
| Test | `bun test` · `pytest` · `bundle exec rspec` · `./gradlew test` · `composer test` · `go test ./...` |
| Lint | language-native (`biome`, `ruff`, `rubocop`, `spotless`, `phpstan`, `golangci-lint`) |
| Replay-record | `bin/record-cassettes` (per-language; needs staging key) |

CI runs on **Blacksmith 2vcpu** runners. Releases are tag-driven per language: `ts-v1.4.2`, `python-v0.9.0`, etc. See [`docs/architecture/release.md`](docs/architecture/release.md).

---

## Architecture

| Doc | Topic |
|-----|-------|
| [versioning.md](docs/architecture/versioning.md) | v1/v2 coexistence, back-compat policy |
| [transport.md](docs/architecture/transport.md)   | retries, caching, rate-limits, idempotency, pagination |
| [errors.md](docs/architecture/errors.md)         | typed-error taxonomy, "good error" definition |
| [resources.md](docs/architecture/resources.md)   | endpoint inventory by version |
| [auth.md](docs/architecture/auth.md)             | bearer token, key-type rules |
| [testing.md](docs/architecture/testing.md)       | unit / replay / smoke layers, cross-language parity |
| [release.md](docs/architecture/release.md)       | Blacksmith CI + per-language tag releases |

Start here: [`docs/architecture/README.md`](docs/architecture/README.md).

---

## Contributing

Issues and PRs welcome. Read [`CLAUDE.md`](CLAUDE.md) and the architecture docs first — public-API changes need to land in **all six languages** in the same PR.

## License

MIT.
