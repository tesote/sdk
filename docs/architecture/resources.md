# Resources

Endpoint inventory the SDK must expose, by API version. **Only client-facing, API-key-authenticated endpoints.** Anything requiring session-cookie auth, admin scope, or internal-only access is excluded.

## v1 — `/api/v1`

Read-only foundation.

| Resource     | Methods       | Endpoint |
|--------------|---------------|----------|
| Status       | `status`, `whoami` | `GET /status`, `GET /whoami` |
| Accounts     | `list`, `get` | `GET /accounts`, `GET /accounts/:id` |
| Transactions | `listForAccount`, `get` | `GET /accounts/:id/transactions`, `GET /transactions/:id` |

## v2 — `/api/v2`

Adds writes for payments + sync orchestration. v1 surface still works at `/v1`.

| Resource           | Methods                                   | Notes |
|--------------------|-------------------------------------------|-------|
| Accounts           | `list`, `get`, `sync`                     | `POST /accounts/:id/sync` triggers a sync |
| Transactions       | `listForAccount`, `get`, `export`, `sync`, `bulk`, `search` | `bulk` and `search` are non-nested |
| Sync sessions      | `list`, `get`                             | scoped under an account |
| Transaction orders | `list`, `get`, `create`, `submit`, `cancel` | scoped under an account |
| Batches            | `create`, `get`, `approve`, `submit`, `cancel` | non-nested under v2 |
| Payment methods    | `list`, `get`, `create`, `update`, `delete` | beneficiaries |
| Status             | `status`, `whoami`                        | |

## v3 — `/api/v3`

Adds reporting, configuration, and webhook delivery. **No upstream OpenAPI doc yet** — derive from controllers under `engines/tesote_api/app/controllers/tesote_api/v3/`.

| Resource           | Methods                                   | Notes |
|--------------------|-------------------------------------------|-------|
| (everything from v2) | same shape                              | re-implemented under `/v3` controllers |
| Balance history    | `listForAccount`                          | `GET /accounts/:id/balance_history` |
| Categories         | `list`, `get`, `create`, `update`, `delete` | |
| Counterparties     | `list`, `get`, `create`, `update`, `delete` | |
| Legal entities     | `list`, `get` (read-only)                 | |
| Connections        | `list`, `get`, `status`                   | bank connections; `status` is a member route |
| Webhooks           | `list`, `get`, `create`, `update`, `delete` | + signature-verification helper (see below) |
| Reports            | `cashFlow`                                | `GET /reports/cash_flow` |
| Workspace          | `get`                                     | read-only `GET /workspace` |
| MCP                | `handle`                                  | `POST /mcp` — pass-through; SDK exposes raw call, not a parsed model |
| Status             | `status`, `whoami`                        | |

## Webhook signature verification (v3)

The SDK ships a stateless helper, not a server:

```ts
import { verifyWebhookSignature } from '@tesote/sdk/v3'
verifyWebhookSignature({ body, signatureHeader, secret })  // throws on mismatch
```

Signature scheme follows the platform's webhook spec — confirm before implementing. Helper lives next to the v3 client; do not pull in HTTP-server dependencies.

## Out of scope (do not expose)

- `ApplicationController` cookie-auth endpoints (the marketing/dashboard app).
- Active Admin routes.
- Internal `/sms_forwarder/*`, `/integration_scanning/*` routes.
- Anything under `app/controllers/` outside `engines/tesote_api/`.
