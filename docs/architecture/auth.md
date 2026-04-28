# Auth

Single scheme: bearer API key.

```
Authorization: Bearer <api_key>
```

No OAuth, no HMAC signing, no session cookies. SDK accepts the key at client construction:

```ts
new V2Client({ apiKey: '...' })             // explicit
new V2Client()                              // falls back to env TESOTE_SDK_API_KEY
```

Resolution order at construction: explicit `apiKey` arg → `TESOTE_SDK_API_KEY` env var → raise `ConfigError` synchronously. Never let a half-built client make a request.

Per-language env-var read:
- TS: `process.env.TESOTE_SDK_API_KEY`
- Python: `os.environ.get('TESOTE_SDK_API_KEY')`
- Ruby: `ENV['TESOTE_SDK_API_KEY']`
- Java: `System.getenv("TESOTE_SDK_API_KEY")`
- PHP: `getenv('TESOTE_SDK_API_KEY')`
- Go: `os.Getenv("TESOTE_SDK_API_KEY")`

## Base URL

Same pattern as the API key — explicit arg overrides env, env overrides default.

```ts
new V2Client({ apiKey: '...', baseUrl: 'https://equipo-staging.tesote.com/api' })  // explicit
new V2Client()  // baseUrl from TESOTE_SDK_API_URL, then default https://equipo.tesote.com/api
```

Env var: `TESOTE_SDK_API_URL`. Default: `https://equipo.tesote.com/api`. Trailing slash stripped at construction. Invalid URL → `ConfigError`.

## Key types

Platform tags API keys with `access_type` (`general`, `odoo`, `sap`). Specialized keys require matching `User-Agent` substrings (`TesoteOdooConnector`, `TesoteSapConnector`). General keys must **not** carry those substrings.

SDK ships only `general` configuration: default User-Agent `tesote-sdk-<lang>/<sdk_version> (<runtime>)`. Consumers integrating from inside Odoo or SAP custom connectors set the User-Agent themselves via the `userAgent` option.

> Server-side matching is currently disabled (see `ApiKeyAccessValidator`); SDK still respects the contract for the day enforcement returns.

## Storage

- SDK never persists the API key.
- Held in memory only on the client instance.
- Never logged. Redact to `Bearer <last4>` in any log/error output.

## Rotating keys

Construct a new client with the new key. Existing clients keep using their old key until garbage-collected. No `rotateKey` method — that encourages mutable state in long-lived clients.

## Multi-workspace

One API key = one workspace, server-side. Multi-workspace consumers instantiate one client per workspace. No "workspace switcher" — caller orchestration.
