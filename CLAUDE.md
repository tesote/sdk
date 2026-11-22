# CLAUDE.md

Every agent run on this repository loads this file as its project brain. `tesote-sdk` is the populated public monorepo for official [equipo.tesote.com](https://equipo.tesote.com) API clients. Read this file before changing a package.

## Stack

| SDK | Location | Package | Runtime floor | Registry |
|---|---|---|---|---|
| TypeScript | `packages/ts/` | `@tesote.com/sdk` | Node 18 | npm |
| Python | `packages/python/` | `tesote-sdk` | Python 3.9 | PyPI |
| Ruby | `packages/ruby/` | `tesote-sdk` | Ruby 3.0 | RubyGems |
| Java | `packages/java/` | `com.tesote:sdk` | Java 17 | Maven Central |
| PHP | `packages/php/` | `tesote/sdk` | PHP 8.1 | Packagist |
| Go | `go/` | `github.com/tesote/sdk/go` | Go 1.21 | proxy.golang.org |
| C# / .NET | `packages/csharp/` | `Tesote.Sdk` | .NET 8 | NuGet |

Use no language feature newer than its floor. Keep runtime dependencies at zero: use each language's standard HTTP and JSON facilities. Java may use `jackson-databind` when `jakarta.json` is insufficient. Keep development dependencies loosely pinned and GitHub Actions on current major releases.

Each SDK is independently implemented, tested, versioned, and released. Do not share implementation across languages. Patch releases are per language. Minor and major releases move all seven SDKs together and require parity validation.

## Commands

| Purpose | Command | Scope |
|---|---|---|
| Set up local dependencies | `./bin/setup` | Go, Python, Ruby, TypeScript |
| Aggregate check | `./bin/check` | Go, Python, Ruby, TypeScript only |
| Validate cross-language parity | `python3 bin/check_parity.py` | All SDK contract inputs |

`./bin/check` is not the complete CI gate. It currently omits Java, PHP, and C#. Run the affected package gate below. CI workflows in `.github/workflows/` are the authoritative package gates and release definitions.

| SDK | Package gate |
|---|---|
| TypeScript | `cd packages/ts && npm ci && npm run typecheck && npm run lint && npm test` |
| Python | `cd packages/python && pip install -e .[test] && ruff check . && mypy src && pytest` |
| Ruby | `cd packages/ruby && bundle install && bundle exec rubocop && bundle exec rspec` |
| Java | `cd packages/java && ./gradlew test --no-daemon` |
| PHP | `cd packages/php && composer install --no-interaction --no-progress --prefer-dist && composer cs-check && composer phpstan -- --no-progress && composer test -- --colors=never` |
| Go | `cd go && test -z "$(gofmt -l .)" && go vet ./... && go test -race -count=1 ./...` |
| C# / .NET | `cd packages/csharp && dotnet restore && dotnet build --configuration Release --no-restore && dotnet test --configuration Release --no-build` |

## Structure

| Path | Purpose |
|---|---|
| `packages/` | TypeScript, Python, Ruby, Java, PHP, and C# SDKs |
| `go/` | Go module and SDK |
| `spec/` | Vendored v1 and v2 OpenAPI snapshots, error codes, and parity inputs |
| `docs/architecture/` | Versioning, transport, errors, resources, auth, testing, and release decisions |
| `.github/workflows/` | Per-language, path-filtered test and release workflows |
| `bin/` | Shared setup, checks, parity validation, and cassette tooling |
| `README.md` | Public SDK overview and consumer quick starts |

End-user documentation and API reference live at `https://www.tesote.com/docs/sdk`. Link there from package READMEs instead of duplicating public reference material here.

## Non-negotiables

### API and models

- Vendor API snapshots into `spec/`. Never import an upstream OpenAPI document during a build.
- Ship `V1Client` and `V2Client` side by side. v1 remains supported after v2 ships. Removing or renaming a public symbol is breaking.
- Expose only client-facing, API-key-authenticated endpoints. Never add internal admin or session-cookie endpoints.
- Keep one resource module per resource. Resource clients marshal parameters, call transport, and deserialize to typed language-native models. Never return raw maps, hashes, or dictionaries.
- Preserve wire casing on the wire and use the public model casing idiomatic to each language.

### Transport and errors

- Keep a swappable transport separate from thin resource clients. Transport alone owns pagination, retry with exponential backoff and jitter, rate-limit awareness, ETag and cache-control caching, opt-in TTL caching, mutation idempotency keys, and request-id propagation.
- Send only `Authorization: Bearer <api_key>` for authentication. POST, PUT, and PATCH requests require `Content-Type: application/json`.
- Use typed errors, one class per API `error_code`. Never collapse API errors into a generic `ApiError` or expose underlying transport exceptions.
- Every API error includes `error_code`, `message`, `http_status`, `request_id`, `retry_after`, `response_body`, and a redacted method-plus-path request summary. Never expose bearer credentials.
- Respect `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`, `Retry-After`, and `X-Request-Id`. Map exhausted retry budgets to `RateLimitExceededError`.
- The API is poll-based. Do not add push assumptions.

### Testing and releases

- Test resource clients with mocked HTTP and fully cover transport behavior. Integration tests replay recorded cassettes. CI never calls the live API.
- Keep smoke tests against staging behind a secret so fork pull requests skip them. Maintain equivalent coverage across all SDKs.
- Each language workflow has a path-filtered test job and a release job that runs only on `main` after tests. Version-source changes trigger idempotent language tags and releases.
- Do not manually publish packages or create release tags. Follow [release.md](docs/architecture/release.md).

### Style

- Prefer small, single-purpose modules. Keep files near 500 lines or split them.
- Lead with the rule. Use fragments and tables for operational documentation.
- Do not use catch-all exception handlers or nil-hiding safe-navigation syntax. Make absence explicit.
- Read [versioning.md](docs/architecture/versioning.md), [transport.md](docs/architecture/transport.md), [errors.md](docs/architecture/errors.md), [resources.md](docs/architecture/resources.md), [auth.md](docs/architecture/auth.md), [testing.md](docs/architecture/testing.md), and [release.md](docs/architecture/release.md) before changing their concern.
