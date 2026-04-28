# Testing

Three layers in every language. Cross-language **parity** is a release-blocker — missing test in PHP is as visible as missing test in TS.

## Layer 1 — Unit (mocked HTTP)

Where: per-resource files (`ts/test/v3/accounts.test.ts`, `ruby/spec/v3/accounts_spec.rb`, etc.).

Asserts:
- Request shape (method, path, query, headers, body) for every public method.
- Response deserialization (typed return shape).
- Error mapping (every `error_code` → typed exception, all required fields populated).
- Transport behavior: retry counts, idempotency-key generation, cache-key generation, rate-limit-header parsing, request-id propagation into errors.

HTTP mocked at the language's standard layer (`nock` for TS, `responses`/`respx` for Python, `WebMock` for Ruby, `okhttp.mockwebserver` for Java, `Guzzle MockHandler` for PHP, `httptest` for Go). Never touch the real network.

Coverage: 90% line, 100% on the transport.

## Layer 2 — Replay (recorded cassettes)

Where: `<lang>/test/replay/` with cassettes in `<lang>/test/replay/cassettes/`.

One representative call per resource recorded against staging (`equipo-staging.tesote.com`) and replayed in CI. Cassettes scrubbed: bearer tokens redacted, real account IDs replaced with synthetic ones, timestamps frozen.

Asserts: SDK round-trips real API responses correctly. Catches drift the moment platform changes a serializer.

Re-record only via `bin/record-cassettes` (per-language script) with a developer's staging key. Recording in CI is forbidden.

## Layer 3 — Smoke (live, gated)

Where: `<lang>/test/smoke/`.

One end-to-end happy path per major resource against staging. Run nightly + on release-tag pipelines. Skipped on PRs from forks (no secret access).

Gating:
- `TESOTE_STAGING_API_KEY` secret must be present (PRs from forks → skip cleanly, not fail).
- Dedicated isolated test workspace; never against production.

## Cross-language parity

`parity-check` CI job verifies:

- Every method in the canonical method list (extracted from `spec/`) exists in every language's client.
- Every `error_code` in the canonical list maps to a typed exception in every language.
- Public method names follow each language's idiomatic casing but resolve to the same canonical name.

Source of canonical lists: `spec/parity.yaml` (hand-maintained for now; codegen later).

## Never test

- Platform behavior — assert the SDK's behavior given a response, not that the API returns a specific response.
- Network primitives — don't test that `fetch` works.
- Codegen output verbatim — test *behavior* of generated code, not formatting.

## Local test commands

| Language | Run all                         | Single file |
|----------|---------------------------------|-------------|
| TS       | `cd ts && bun test`             | `bun test test/v3/accounts.test.ts` |
| Python   | `cd python && pytest`           | `pytest tests/v3/test_accounts.py::test_list` |
| Ruby     | `cd ruby && bundle exec rspec`  | `bundle exec rspec spec/v3/accounts_spec.rb:42` |
| Java     | `cd java && ./gradlew test`     | `./gradlew test --tests V3AccountsTest.list` |
| PHP      | `cd php && composer test`       | `vendor/bin/phpunit tests/V3/AccountsTest.php` |
| Go       | `cd go && go test ./...`        | `go test ./v3 -run TestAccountsList` |
