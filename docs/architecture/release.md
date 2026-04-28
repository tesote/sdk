# Release

Per-language, version-file-driven, automated. All jobs on **Blacksmith 2vcpu** (`runs-on: blacksmith-2vcpu-ubuntu-2204`).

## How a release happens

1. Bump the version source file for the target language (table below).
2. Update `<lang>/CHANGELOG.md`.
3. Merge to `main`. The workflow detects the version change, publishes to the package manager, creates the tag, opens a GitHub Release.

No manual `git tag`, no `gh release create` by hand, no `npm publish`. The workflow does it. Idempotent — if the tag already exists, the release job is a no-op.

## Version source per language

| Language   | Version source file                              | How the workflow reads it                                             |
|------------|--------------------------------------------------|-----------------------------------------------------------------------|
| TypeScript | `packages/ts/package.json` (`version`)           | `node -p "require('./package.json').version"`                         |
| Python     | `packages/python/pyproject.toml` (`project.version`) | `python -c "import tomllib, pathlib; print(tomllib.loads(...).['project']['version'])"` |
| Ruby       | `packages/ruby/lib/tesote_sdk/version.rb`        | `ruby -r ./lib/tesote_sdk/version.rb -e 'print TesoteSdk::VERSION'`   |
| Java       | `packages/java/build.gradle.kts` (`version`)     | `grep -E '^version = ' build.gradle.kts \| sed -E 's/.*"([^"]+)".*/\1/'` |
| PHP        | `packages/php/VERSION`                           | `tr -d '[:space:]' < VERSION`                                         |
| Go         | `packages/go/version.go` (`const Version`)       | `grep -oE 'Version = "[^"]+"' version.go \| sed -E 's/.*"([^"]+)".*/\1/'` |

Every language: bumping that file is what triggers a release. No tags pushed by humans.

## Tag scheme

One namespace per language, semver: `ts-v1.4.2`, `python-v0.9.0`, `ruby-v2.0.0`, `java-v1.1.0`, `php-v0.5.3`, `go-v3.0.0`. The release job creates the tag; `proxy.golang.org` and Packagist consume them. npm, PyPI, RubyGems, Maven Central use direct publish from the workflow.

Pre-release tags allowed: `ts-v1.5.0-rc.1`, `python-v1.0.0-beta.2`. Pre-releases publish to the registry's pre-release channel (npm `--tag next`, PyPI `pre-release` flag).

## Workflow files

```
.github/workflows/
├── ts.yml  python.yml  ruby.yml  java.yml  php.yml  go.yml   ← two jobs each: test → release
└── parity-check.yml                                          ← cross-language method/error parity (on PR)
```

### Two jobs per language

| Job       | Trigger                                                          | Purpose                                                                                      |
|-----------|------------------------------------------------------------------|----------------------------------------------------------------------------------------------|
| `test`    | push/PR matching workflow-level `paths:` filter                  | matrix across floor + latest LTS + current stable; lint, typecheck, unit + integration replay |
| `release` | `needs: test`, `if: push to main`                                | reads version source, checks if `<lang>-v<version>` tag exists, publishes + tags + GH Release |

Path filtering happens **at the workflow level** — a PR touching only `packages/python/` doesn't even queue TS, Ruby, Java, PHP, Go. No separate `detect` job needed.

The release job is **idempotent**: if the tag already exists (no version bump in this push), it short-circuits at the gate step.

## Trusted publishers (OIDC)

npm, RubyGems, PyPI publish via GitHub Actions OIDC trusted publishers. Workflows declare:

```yaml
permissions:
  id-token: write
```

No `NPM_TOKEN`, no `RUBYGEMS_API_KEY`, no `PYPI_API_TOKEN` — registries verify the workflow identity directly.

**One-time setup per registry** (not codified — done once in each registry's UI):

| Registry | Where                                                       |
|----------|-------------------------------------------------------------|
| npm      | package settings → trusted publishers → add repo + workflow |
| RubyGems | gem settings → OIDC                                         |
| PyPI     | project settings → publishing → add trusted publisher       |

Maven Central: Sonatype Central Portal user token (NOT legacy OSSRH) — token-based, no OIDC. Packagist: GitHub webhook indexes new tags — no token in the workflow. Go: tag push via proxy — no token, no OIDC.

## Secrets

Stored in repo secrets (settings → secrets → actions):

| Secret                                                  | Used by         |
|---------------------------------------------------------|-----------------|
| `MAVEN_CENTRAL_USERNAME`, `MAVEN_CENTRAL_PASSWORD`      | java-release (Sonatype Central Portal user token, NOT legacy OSSRH) |
| `MAVEN_GPG_KEY`, `MAVEN_GPG_PASSPHRASE`                 | java-release (artifact signing) |
| `TESOTE_STAGING_API_KEY`                                | smoke tests     |

npm / RubyGems / PyPI: OIDC — no secrets. Go: tag push — no secret. PHP / Packagist: webhook on tag push — no secret.

Never echo secrets, never write to logs, never check into the repo.

## Versioning rules

Strict semver:

- **Major** — back-compat break (per [versioning.md](versioning.md)). Vanishingly rare. New API version (`V3Client`) is **minor**, not major — purely additive.
- **Minor** — additive: new resources, methods, optional params, error subclasses.
- **Patch** — bug fixes, internal refactors, doc-only changes.

Patch is per-language only. Minor and major land across all six in lockstep, gated by `parity-check.yml`.

Changelog entries in `<lang>/CHANGELOG.md`. Release workflow uses commits since the prior `<lang>-v*` tag for the GitHub Release notes.
