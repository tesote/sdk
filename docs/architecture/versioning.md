# Versioning

Two independent axes. Don't conflate.

| Axis                                  | Tracks                                                                            | Changes                                              |
|---------------------------------------|-----------------------------------------------------------------------------------|------------------------------------------------------|
| **API version** (v1, v2)              | URL prefix on the server (`/api/v2/...`); coherent set of resources and semantics | Platform team bumps; SDKs follow on new ship          |
| **SDK version** (semver per language) | Public surface of the SDK package                                                 | Every release; per language, independently           |

## API versions ship side-by-side

Every SDK exports all currently-supported API versions as named clients. Consumer chooses:

```ts
import { V1Client, V2Client } from '@tesote.com/sdk'
const accounts = await new V2Client({ apiKey }).accounts.list()
```

```python
from tesote_sdk import V1Client, V2Client
```

```ruby
TesoteSdk::V2::Client.new(api_key: ...)
```

```go
import "github.com/tesote/sdk/go/v2"   // major-version subpath per Go module rules
```

```java
import com.tesote.sdk.v2.V2Client;
```

```php
use Tesote\Sdk\V2\Client as V2Client;
```

Mix versions in one process — `V1Client.transactions.list()` for legacy, `V2Client.batches.create()` for new, sharing only the auth token.

## What's in each version

Full per-resource inventory: [resources.md](resources.md).

- **v1** — read-only: accounts, transactions
- **v2** — v1 + sync sessions, transaction orders, batches, payment methods, bulk + search

Within each version, SDK matches the API endpoint surface 1:1 — no convenience methods spanning versions.

## Back-compat policy

**Removing or renaming a public symbol in any shipped API version is forbidden, in every language.** Includes:

- Removing a resource client (`V1Client.accounts` exists forever).
- Renaming a method, field, or enum value on a returned model.
- Non-additive signature changes (new required param = breaking; new optional param = OK).
- Tightening accepted input types.

Allowed without major bump:

- New versioned clients (`V3Client` lands → minor bump).
- New resources, methods, optional params on existing clients.
- New error subclasses (callers catching the parent still work).
- Internal refactors preserving public surface.

## Deprecation

Platform marks an endpoint deprecated:

1. SDK keeps the method.
2. Add a runtime deprecation warning (idiomatic — `warnings.warn` in Python, `console.warn` in TS, etc.).
3. README + `tesote.com/docs/sdk` page flag it.
4. Method is **not** removed when the upstream endpoint is removed — it throws a typed `EndpointRemovedError` pointing at the replacement.

## Spec snapshots

`spec/` vendors a frozen copy of each version's OpenAPI doc:

```
spec/
├── v1.openapi.yaml
└── v2.openapi.yaml
```

Codegen reads from `spec/`, not the live API. Bumping a snapshot is a deliberate PR with a per-language changelog entry.
