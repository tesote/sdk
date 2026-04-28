# Versioning

Two independent version axes. Don't conflate them.

| Axis        | What it tracks                  | How it changes |
|-------------|---------------------------------|----------------|
| **API version** (v1, v2, v3) | URL prefix on the server (`/api/v3/...`); a coherent set of resources and semantics | Only the platform team bumps this; SDKs follow when a new version ships |
| **SDK version** (semver per language) | Public surface of the SDK package | Bumped on every release; per language, independently |

## API versions ship side-by-side

Every SDK exports all currently-supported API versions as named clients. The consumer chooses:

```ts
import { V1Client, V2Client, V3Client } from '@tesote/sdk'
const accounts = await new V3Client({ apiKey }).accounts.list()
```

```python
from tesote_sdk import V1Client, V2Client, V3Client
```

```ruby
TesoteSdk::V3::Client.new(api_key: ...)
```

```go
import "github.com/tesote/sdk/go/v3"   // major-version subpath per Go module rules
```

```java
import com.tesote.sdk.v3.V3Client;
```

```php
use Tesote\Sdk\V3\Client as V3Client;
```

A consumer can mix versions in one process — call `V1Client.transactions.list()` for legacy code paths and `V3Client.webhooks.create()` for new ones, sharing nothing but the auth token.

## What's in each version

See [resources.md](resources.md) for the full per-resource inventory. Summary:

- **v1** — read-only: accounts, transactions
- **v2** — v1 + sync sessions, transaction orders, batches, payment methods, bulk + search
- **v3** — v2 + categories, counterparties, legal entities, connections, webhooks, reports, balance history, workspace, MCP

Within each version, the SDK matches the API's endpoint surface 1:1 — no convenience methods that span versions.

## Back-compat policy

**Removing or renaming a public symbol in any shipped API version is forbidden, in every language.** This includes:

- Removing a resource client (`V1Client.accounts` must exist forever).
- Renaming a method, field, or enum value on a returned model.
- Changing a method signature in a non-additive way (adding a required param = breaking; adding an optional param = OK).
- Tightening accepted input types.

Allowed without a major bump:

- Adding new versioned clients (`V4Client` lands → minor bump).
- Adding new resources, methods, or optional params to existing clients.
- Adding new error subclasses (callers catching the parent class still work).
- Internal refactors (transport, serialization) that preserve the public surface.

## Deprecation

When the platform marks an endpoint deprecated:

1. SDK keeps the method.
2. SDK adds a runtime deprecation warning (language-idiomatic — `warnings.warn` in Python, `console.warn` in TS, etc.).
3. README + the doc page on `tesote.com/docs/sdk` flag it.
4. The method is **not** removed even when the upstream endpoint is removed — it then throws a typed `EndpointRemovedError` pointing at the replacement.

## Spec snapshots

`spec/` vendors a frozen copy of each version's OpenAPI doc:

```
spec/
├── v1.openapi.yaml
├── v2.openapi.yaml
└── v3.openapi.yaml      ← TODO: derive from v3 controllers; upstream lacks one
```

Codegen reads from `spec/`, not from the live API. Bumping a snapshot is a deliberate PR with a changelog entry per affected language.
