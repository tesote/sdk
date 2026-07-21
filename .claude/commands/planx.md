---
description: Write a concise, self-contained execution plan to docs/plans/<YYYY>/<MM>/<DD>/<1NN>-<slug>/ for another AI to implement
argument-hint: [what you want done]
allowed-tools: Write, Read, Glob, Grep, Task, Bash
---

# /planx

Produce a concise plan another AI can execute with zero extra context. Plan only — no implementation, no code execution, no edits outside the plan dir.

## Goal
$ARGUMENTS

## Steps

1. **Resolve path.** Run `date +%Y`, `date +%m`, `date +%d`. Dir = `docs/plans/<YYYY>/<MM>/<DD>/`. `Glob docs/plans/<YYYY>/<MM>/<DD>/1*` → next number = highest existing `1NN-*` + 1, else `101`. Slug = kebab-case title, max 5 words. Final plan dir: `docs/plans/<YYYY>/<MM>/<DD>/<1NN>-<slug>/`.

2. **Explore.** `Task` (subagent_type=Explore, thoroughness="very thorough"): the upstream OpenAPI contract (`../<platform>/engines/tesote_api/docs/openapi{,_v2}.yaml`), the affected resource module + transport in each language, the typed-error taxonomy, test layers, and the design docs in `docs/architecture/`. Skip only for trivial asks.

   **No git worktrees.** Explore agents — and any executor agents later — run directly in this checkout. Never `isolation: worktree`, never a per-agent worktree dir. Parallel agents share this one tree; give each a disjoint path (one language each).

3. **Write the plan as multiple files** in the plan dir — never one big `plan.md`. `overview.md` index plus one `<NN>-<aspect>.md` per separable area. Here the natural split is **one slice per language** plus a shared contract slice (e.g. `01-contract.md`, `02-typescript.md`, `03-python.md`, … `09-docs.md`).

   **`overview.md`** — the map. Sections:

```markdown
# <Title>

## Goal
1-2 sentences: what + why.

## Context
- Stack facts the executor needs (7 languages, zero runtime deps, min-version floors, versioned V1Client/V2Client, transport owns cross-cutting concerns).
- Reference patterns: `packages/ts/src/resources/<thing>.ts:12` — follow this for Z.

## Plan files (execute in order)
1. [`01-<aspect>.md`](01-<aspect>.md) — one line: what it covers.

## Done when
- Verifiable acceptance criteria spanning the whole change.

## Risks / open questions
```

   **Each `<NN>-<aspect>.md`** — one slice: `## Files to change` (`path:line` + why), `## Steps` (ordered, concrete), `## Tests` (unit with mocked HTTP + cassette replay; the language's lint + test command), `## Done when`.

4. **Write a `status.yml`** in the plan dir — the live tracker. New plans start `not_started` / `0%`. `created_by` + `owner` from `git config user.name`; leave `worked_by: ""` for the executor to claim.

```yaml
plan: <1NN>-<slug>
title: <human title from overview.md>
status: not_started        # not_started | in_progress | blocked | complete | superseded
created_by: <git config user.name>
worked_by: ""
owner: <git config user.name>
percent: 0
current_focus: ""
slices:
  - file: 01-<aspect>.md
    status: not_started
    percent: 0
evidence: []
notes: ""
last_updated: <YYYY-MM-DD>
```

## Rules
- Compact English. Fragments over sentences. `file:line` / symbol refs over prose. Tables for structured data.
- Reference-only: point at code, don't paste it.
- No checkboxes. Multiple files always. Self-contained.
- Respect `CLAUDE.md`: zero runtime deps, min-version floor per language, typed models not raw maps, one error class per `error_code`, transport owns pagination/retry/caching/idempotency, permanent back-compat, cross-language coverage parity, files <500 LOC.
- Note the release trigger in any slice that changes a public surface: bump the language's version source file; minor/major land across all seven in lockstep.
- **No git worktrees** — plans execute in this checkout; never instruct an executor to create one.

## Output
```
✓ docs/plans/<YYYY>/<MM>/<DD>/<1NN>-<slug>/overview.md
  + 01-<aspect>.md, … + status.yml
Next: run an executor on overview.md.
```
