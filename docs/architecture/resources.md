# Resources

Endpoint inventory the SDK must expose, by API version. **Only client-facing, API-key-authenticated endpoints.** Anything requiring session-cookie auth, admin scope, or internal-only access is excluded.

## v1 — `/api/v1`

Read-only foundation.

| Resource     | Methods                 | Endpoint |
|--------------|-------------------------|----------|
| Status       | `status`, `whoami`      | `GET /status`, `GET /whoami` |
| Accounts     | `list`, `get`           | `GET /accounts`, `GET /accounts/:id` |
| Transactions | `listForAccount`, `get` | `GET /accounts/:id/transactions`, `GET /transactions/:id` |

## v2 — `/api/v2`

Adds writes for payments + sync orchestration. v1 surface still works at `/v1`.

| Resource           | Methods                                                     | Notes |
|--------------------|-------------------------------------------------------------|-------|
| Accounts           | `list`, `get`, `sync`                                       | `POST /accounts/:id/sync` triggers a sync |
| Transactions       | `listForAccount`, `get`, `export`, `sync`, `bulk`, `search` | `bulk` and `search` non-nested |
| Sync sessions      | `list`, `get`                                               | scoped under an account |
| Transaction orders | `list`, `get`, `create`, `submit`, `cancel`                 | scoped under an account |
| Batches            | `create`, `get`, `approve`, `submit`, `cancel`              | non-nested under v2 |
| Payment methods    | `list`, `get`, `create`, `update`, `delete`                 | beneficiaries |
| Status             | `status`, `whoami`                                          | |

## Out of scope (do not expose)

- `ApplicationController` cookie-auth endpoints (marketing/dashboard app).
- Active Admin routes.
- Internal `/sms_forwarder/*`, `/integration_scanning/*` routes.
- Anything under `app/controllers/` outside `engines/tesote_api/`.
