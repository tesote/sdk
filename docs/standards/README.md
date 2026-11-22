# Standards

> established by the developerz.ai maintainer agent for the maintainers of this repository from what onboarding read and ran in this repo, then owned by
> its maintainers. Correct it in place; setup never overwrites a page that exists.

- Preserve v1 and v2 symbols and implement public API parity across SDKs. Enforced by `python3 bin/check_parity.py` and `.github/workflows/parity-check.yml`. `docs/architecture/versioning.md:50-64`, `spec/parity.yaml:1-8`
- Keep resource clients thin and transport-owned cross-cutting behavior centralized. Checked by per-language mocked transport tests. [`CLAUDE.md`'s Non-negotiables](../../CLAUDE.md#non-negotiables), `docs/architecture/testing.md:5-17`
- Use typed models and typed errors; redact bearer data. Checked by language test suites. [`CLAUDE.md`'s Non-negotiables](../../CLAUDE.md#non-negotiables), `docs/architecture/errors.md:18-32`
- Run package-native gates before merge. Commands are defined in `.github/workflows/ts.yml:48-51`, `.github/workflows/python.yml:40-49`, `.github/workflows/ruby.yml:41-44`, `.github/workflows/java.yml:36-38`, `.github/workflows/php.yml:46-53`, `.github/workflows/go.yml:47-57`, and `.github/workflows/csharp.yml:48-53`.
- Do not manually publish or create tags. Enforced by release jobs in `.github/workflows` and policy in `docs/architecture/release.md:5-11`.
