# Auth

Single auth scheme: bearer API key.

```
Authorization: Bearer <api_key>
```

No OAuth, no HMAC signing, no session cookies. The SDK accepts the key at client construction:

```ts
new V3Client({ apiKey: process.env.TESOTE_API_KEY })
```

If `apiKey` is missing or empty, raise `ConfigError` synchronously at construction — never let a half-built client make a request.

## Key types

The platform tags API keys with an `access_type` (`general`, `odoo`, `sap`). Specialized keys require matching `User-Agent` substrings (`TesoteOdooConnector`, `TesoteSapConnector`). General keys must **not** carry those substrings.

The SDK ships only a `general` configuration: its default User-Agent is `tesote-sdk-<lang>/<sdk_version> (<runtime>)`. Consumers integrating from inside Odoo or SAP custom connectors set the User-Agent themselves via the `userAgent` option.

> Note: the matching check is currently disabled server-side (see `ApiKeyAccessValidator`), but the SDK still respects the contract so it works the day enforcement is re-enabled.

## Storage

- The SDK never persists the API key.
- The key is held in memory only on the client instance.
- The key is never logged. Redact to `Bearer <last4>` in any log/error output.

## Rotating keys

Construct a new client with the new key. Existing clients keep using their old key until garbage-collected. Don't add a "rotateKey" method — it encourages mutable state in long-lived clients.

## Multi-workspace

One API key is scoped to one workspace server-side. Multi-workspace consumers instantiate one client per workspace. The SDK does not provide a "workspace switcher" — that's caller orchestration.
