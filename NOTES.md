# sdk

Public monorepo holding the official client SDKs for the equipo.tesote.com API (repo name
`tesote-sdk`). Seven language implementations live side by side, each independently versioned,
tested and released, with no cross-language code sharing — every SDK is written idiomatically for its
own ecosystem. Runtime dependencies are deliberately zero: stdlib only for HTTP, JSON, retries and
caching. Clients are versioned side by side (`V1Client`, `V2Client`) and back-compat is permanent.

- **Stack:** TypeScript (`@tesote.com/sdk`, npm), Python (`tesote-sdk`, PyPI), Ruby (RubyGems),
  Java (`com.tesote:sdk`, Maven Central), PHP (`tesote/sdk`, Packagist), Go (`github.com/tesote/sdk/go`),
  C#/.NET (`Tesote.Sdk`, NuGet). CI on Blacksmith runners; releases are triggered by bumping a
  language's version file, published via OIDC trusted publishers where available.
- **Key commands:** each language has its own toolchain and its own workflow at
  `.github/workflows/<lang>.yml` (a `test` job plus a `release` job gated by a `paths:` filter).
  `bin/check_parity.py` enforces cross-language surface parity against `spec/parity.yaml`.
- **Layout:**
  - `packages/{ts,python,ruby,java,php,csharp}/` — per-language SDK sources
  - `go/` — the Go SDK at repo root (Go module path requirement): transport, cache, backoff, typed
    errors, `v1/` and `v2/` clients
  - `spec/` — vendored OpenAPI snapshots (`v1.openapi.yaml`, `v2.openapi.yaml`), `error_codes.txt`,
    `parity.yaml`
  - `docs/architecture/` — deep dives on versioning, transport, errors, resources, auth, testing, release
  - `bin/check_parity.py` — parity checker used by CI
- **Design rules worth knowing:** transport owns pagination, retry with jitter, rate-limit awareness,
  caching and idempotency keys; resource clients stay thin and deserialize into typed model objects;
  one error class per `error_code`, never a single collapsed `ApiError`.
- **State as of 2026-07-21:** branch `main`, working tree clean (no uncommitted work).
  Last commit before this note: `7133716` (2026-04-28).
  Note: CLAUDE.md still describes the repo as "greenfield", but all seven language packages plus the
  Go SDK are now implemented — treat that line as stale.
