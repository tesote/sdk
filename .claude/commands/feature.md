---
description: End-to-end feature workflow for the tesote-sdk polyglot monorepo — understand, explore, build across the seven languages, test, PR, release. Reads intent from the prompt.
argument-hint: <what you want built, plain language> [+ reference URL(s)]
allowed-tools: Read, Write, Edit, Glob, Grep, Bash, Task, Skill, WebFetch
---

# /feature

You are a **senior engineer on the tesote-sdk team**. Take a change from plain-language idea to merged and released.

## Request
$ARGUMENTS

**The prompt is the context — read the intent.** How autonomous to be, which languages, whether to confirm before merging: infer it from the words. The flow below is a map, not a checklist to recite. Stop for a true blocker (a public-surface break, a missing upstream OpenAPI contract, a release credential you don't have).

## No git worktrees

**Work directly in this checkout.** Never `isolation: worktree`, never a per-agent worktree dir. If you fan out parallel `Task` agents (one per language is the natural split), they all share **this one tree** — so coordinate: each agent owns a disjoint path (`packages/ts/`, `packages/python/`, `packages/ruby/`, `packages/java/`, `packages/php/`, `go/`, `packages/csharp/`), no concurrent `git checkout` or branch switch, and each stages only its own paths.

## The flow

1. **Understand.** Restate the goal in a line. A surface change is defined by the upstream OpenAPI (`../<platform>/engines/tesote_api/docs/openapi{,_v2}.yaml`) — read it before designing. Only client-facing, API-key-authenticated endpoints are in scope.

2. **Explore.** Fan out `Task` Explore agents (very thorough) to map the affected resource module in every language, the transport layer, the typed-error taxonomy, and the tests to mirror (`file:line`). `docs/architecture/` is the design source of truth (`versioning.md`, `transport.md`, `errors.md`, `resources.md`, `testing.md`, `release.md`).

3. **Build — one language first, then port.** Land the change in one language (usually TS) with its tests, settle the shape, then port idiomatically to the other six. **No cross-language code sharing** — duplicate idiomatically. Rules that never bend: zero runtime deps (stdlib only), respect the min-version floor, typed models not raw maps, one error class per `error_code`, cross-cutting concerns (pagination, retry, caching, idempotency, request-id) live in the transport, never in a resource client. **Back-compat is permanent** — never remove or rename a public symbol.

4. **Verify.** Unit tests with mocked HTTP per resource client; integration via recorded cassette replay (never the live API). Cross-language coverage parity — a missing PHP test is as visible as a missing TS test. Run the affected language's lint + tests locally before pushing.

5. **PR + merge.** Branch `<type>/<slug>`, Conventional Commits scoped to the language (`feat(ruby): …`). One concern per PR. `gh pr create` with Summary + Test plan. CI runs only the affected language (workflow `paths:` filter); wait for green, address review comments, then merge.

6. **Release.** Bumping the language's version source file is what triggers a release — no human pushes tags. Patch is per-language; **minor and major land across all seven in lockstep**, gated by `parity-check.yml`. Confirm the `release` job published and tagged `<lang>-v<version>`.

7. **Docs.** Public-surface changes update the matching page at `../tesote.com` (`www.tesote.com/docs/sdk`) in the same PR — don't duplicate end-user docs here.

## Output

```
Languages: <n> of 7 touched → <ts, python, …>
PRs:       #… (merged)
Released:  <lang>-v<version> …   parity: <lockstep | patch-only>
Docs:      <../tesote.com PR or n/a>
```
