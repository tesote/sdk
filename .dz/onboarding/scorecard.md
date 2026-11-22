# AI-first scorecard

> This card is the repository's own — established by the developerz.ai maintainer agent for the maintainers of this repository; each onboarding run rewrites the file whole, the latest run wins, and the same card is stored on the platform.

- Run: `run_a2e90f2bf87c40bdb02e13008f0680a9`
- Computed: 2026-11-22T17:19:59.478Z
- Proof: **unprovable** — `bin/check` did not run green at head on this box. Nothing in the repository’s own code was changed to force it; the excerpt below is the record.
- Score: 10 of 15 items ok

## `policy` — ok

- .dz/maintainer/maintainer.yml established
- .dz/maintainer/reviewer.md established, still the unfilled template

## `brain` — ok

- CLAUDE.md already present, left alone

## `roster` — ok

- .claude/agents/transport.md established
- .claude/agents/api-parity.md established
- .claude/agents/release.md established

## `skills` — ok

- .claude/skills/verify-change/SKILL.md established, still the unfilled template
- .claude/skills/parity-check/SKILL.md established

## `scripts` — unprovable

- bin/setup established
- bin/dev established, still the unfilled template
- bin/check established
- bin/stop established
- down: ./bin/stop (bin)
- ./bin/setup at head: red, exit 1 (0.0s)

## `wiring` — ok

- .mcp.json established

## `commands` — ok

- .claude/commands already present, left alone

## `access` — gap

- landed 0 of 10 asked
- still missing: TESOTE_API_KEY, TESOTE_SDK_API_KEY, TESOTE_SDK_API_URL, TESOTE_STAGING_API_KEY, MAVEN_CENTRAL_USERNAME, MAVEN_CENTRAL_PASSWORD, MAVEN_GPG_KEY, MAVEN_GPG_PASSPHRASE, NUGET_API_KEY, GITHUB_TOKEN

## `layout` — gap

- partial house layout
- no apps/ at the repository root

## `smoke_test` — unprovable

- packages/ts/test/transport.test.ts in the repository
- ./bin/check at head: red, exit 1 (0.0s)

## `env_example` — ok

- .env.example established

## `agents_md` — ok

- AGENTS.md established

## `per_workspace_brains` — ok

- packages/ts/CLAUDE.md established
- packages/python/CLAUDE.md established
- packages/ruby/CLAUDE.md established
- packages/java/CLAUDE.md established
- packages/php/CLAUDE.md established
- go/CLAUDE.md established
- packages/csharp/CLAUDE.md established

## `docs_readme` — ok

- docs/README.md established

## `gate_proved` — unprovable

- ./bin/check at head: red, exit 1 (0.0s)
