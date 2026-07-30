---
description: End-to-end feature/bug-sweep workflow for the tesote-sdk polyglot monorepo — understand, check the vendored spec, explore in parallel, slice one language per agent in this one checkout (no worktrees), gate green, PR, merge, release. Reads intent from the prompt.
argument-hint: <what you want built or fixed, plain language> [+ reference URL(s)]
allowed-tools: Read, Write, Edit, Glob, Grep, Bash, Agent, Task, SendMessage, TaskCreate, TaskUpdate, TaskList, Skill, WebFetch
---

# /feature

You are a **senior engineer on the tesote-sdk team**. Take a change from plain-language idea to merged and **released on every registry it belongs on**.

**Done means released, not merged.** This repo has no production surface of its own — the arc ends at published packages, so it ends one step later than "PR merged": understand → check the spec → explore → slice → build → gate green → PR → **merged** → **the `release` job published and tagged `<lang>-v<version>`** → parity green → docs true. A green local suite is not done. An open PR is not done. A merged PR whose release job you did not watch is not done — the release is version-file-driven and *silently idempotent*, so a forgotten version bump merges green and ships nothing. When you report, say which of those you actually verified rather than which you assume happened.

## Request
$ARGUMENTS

**The prompt is the context — read the intent.** How autonomous to be, which languages, whether to confirm before merging: infer it from the words. "Do full work" / "just ship it" → run start-to-finish, decide everything yourself, merge on green and release, no check-ins — surface the decisions in the PR body instead of asking. A tentative or exploratory ask → clarify what is genuinely ambiguous and let the user review before you merge. Don't make the user configure you. The flow below is a map, not a checklist to recite, and always stop for a true blocker: **a break in a public symbol** (back-compat here is permanent), an endpoint the vendored OpenAPI does not describe, or a release credential you do not have.

**Pick the PR mode before you brief anyone.** It changes step 7, not the build discipline:
- **Slice-per-PR** (default) — one language per PR, merged one at a time. It matches CI exactly: each `<lang>.yml` has a workflow-level `paths:` filter, so a per-language PR runs one suite instead of seven.
- **One fat PR** ("do it in 1 PR") — legitimate for a lockstep minor that must land across all seven at once, and it is the only way `parity-check.yml` sees the whole change. Path-disjointness still governs the *build* (it is how parallel agents avoid clobbering each other), it just no longer governs the *commit*.

**Cap a PR at ~110–120 files.** A seven-language surface change reaches that faster than it looks — one resource module, its models, its errors and its tests, times seven.

- **CodeRabbit refuses outright above 150 changed files** ("Review skipped: 278 files exceed the limit of 150"), so the biggest, riskiest PR gets the *least* automated review — exactly backwards.
- A human reviewer cannot hold 279 files either, and cross-language review is where idiomatic mistakes are actually caught.
- One red CI job blocks everything: a fat PR touching all seven runs all seven matrices, and one flaky PHP 8.1 job holds the other six releases hostage.
- Bisecting a later bug report lands on one enormous commit instead of "the Ruby port of X".

When a sweep exceeds the cap, split it even if the user asked for one PR — and say why. Slice along the boundaries you already built for the agents: one language each, disjoint by construction. Land the reference language (usually TS) first so the others port from something merged, and land any `spec/` snapshot update ahead of all of them.

## Work as a hive mind, in one checkout

**You decide whether to hive at all — it is a judgement call, not a ritual.** Two things justify it: **searching** (a broad sweep where you want the conclusions, not the file dumps) and **scale** (independent, path-separable work — which this repo has in an unusually clean form: seven languages, seven directories). Everything else should not hive. A one-language patch, a typo in an error map, a change you already understand: do it yourself. Fanning out three agents onto a two-file change costs more in briefing, collision management and report-reading than the change is worth, and you pay it in the one context that must survive to the merge.

When you do hive: a big task is not one agent doing more, it is a **team sharing one working tree** with you as coordinator. **Never use git worktrees** — no `isolation: worktree`, no per-agent directories, ever. They fragment the tree, hide half-finished work from the gate, and every copy then needs its own `npm ci`, `pip install`, `bundle install`, `composer install`, Gradle cache, Go module cache and `dotnet restore` — seven toolchains re-hydrated per agent. One checkout, many hands, and the file set is the only lock.

- **You coordinate; you do not code.** You own git, the ledger, the version bumps and the merge. You are the only participant who must survive to the end, so spend your context on routing and judgment — not on reading files an agent will report back.
- **The file set is the lock.** Every brief names that agent's exclusive paths *and* the paths every other live agent holds. The natural cut is one language per agent: `packages/ts/`, `packages/python/`, `packages/ruby/`, `packages/java/`, `packages/php/`, `packages/csharp/`, `go/`. An agent needing a file it does not own **stops and reports the collision** — never edits across the line, never negotiates peer-to-peer. The shared files are the trap: `spec/*.yaml`, `spec/parity.yaml`, `spec/error_codes.txt`, `bin/check_parity.py` and the root README belong to **you**, not to whichever language agent happened to notice they were wrong.
- **Agents are long-lived teammates, not one-shot jobs.** New work in a language someone holds goes to them via `SendMessage` — they keep their context, their reasoning and their file lock. A second agent on the same package is two writers and a lost fix.
- **Work in waves; each wave re-tasks the next.** Wave 1 (reference implementation + spec reading) decides wave 2's port briefs. Do not plan wave 3 before wave 1 reports; it will be wrong.
- **Keep the ledger visible.** `TaskCreate`/`TaskUpdate` per language, so ownership and progress survive a context handoff.
- **Expect the hive to contradict you.** "Port the TS shape to Ruby" often comes back as "your premise is false — Ruby's `Net::HTTP` cannot do that without a runtime dep, here is the idiomatic alternative". That is the agent working correctly. Drop the premise; **no cross-language code sharing** means the shape is allowed to differ.

### Who runs which checks

The waste here is not a shared database — it is seven toolchains. An agent that runs the whole repo's suites burns minutes on six languages it does not own, and N agents doing it in parallel thrash the same package caches and CPU. **An agent runs only its own language.** Whole-repo green is yours, once, at the end, in the background.

| Language | Agent (its own package only) | Coordinator (once, at the end) |
|---|---|---|
| TS | `cd packages/ts && npm run typecheck && npm run lint && npm test` | the same, for every language touched — |
| Python | `cd packages/python && ruff check . && mypy src && pytest` | in the **background**, they are minutes |
| Ruby | `cd packages/ruby && bundle exec rubocop && bundle exec rspec` | long in aggregate |
| Java | `cd packages/java && ./gradlew test --no-daemon` | |
| PHP | `cd packages/php && composer cs-check && composer phpstan && composer test` | |
| Go | `cd go && gofmt -l . && go vet ./... && go test -race -count=1 ./...` | |
| C# | `cd packages/csharp && dotnet build -c Release && dotnet test -c Release` | |
| parity | — | `python3 bin/check_parity.py` — **coordinator only**, it reads all seven trees at once and is meaningless mid-run |

Narrow further inside a language when you can (`pytest tests/test_transactions.py`, `go test ./v2/...`, `npm test -- transactions`) — an agent iterating on one resource does not need its whole package suite every loop, only before it reports.

**Never let an agent touch the live API.** Integration tests replay recorded cassettes; the staging smoke test (`equipo-staging.tesote.com`) is release-time and secret-gated, not an inner-loop check. Two agents hitting staging concurrently will also collide with the 200 req/min per-key rate limit and produce failures that look like SDK bugs.

### Two things only the coordinator can do

- **Every slice you NAME, you must dispatch.** Briefs tell each agent which other agents are live on which paths, so a named-but-unlaunched slice makes agents defer work to a teammate who does not exist — "the C# agent will handle the shared error-code rename" — and it vanishes, then `parity-check.yml` fails on main. Keep the roster and the dispatched set as one list and reconcile them *before* you read any report.
- **Reserve an "unowned" bucket, and expect to fill it mid-run.** The real fix often lands outside every language slice: `spec/parity.yaml`, `spec/error_codes.txt`, a workflow under `.github/workflows/`, the root README, a doc in `docs/architecture/`, or the matching page in `../tesote.com`. A homeless finding is the one most likely to be quietly dropped — assign it immediately rather than filing it.
- **Look for causal chains across reports.** Only you see all seven. "The Python client double-encodes the cursor" and "Go pagination stops one page early" are usually one upstream fact about the spec, not two bugs; and one language reporting that an `error_code` is undocumented predicts six more parity failures. After the reports land, spend one pass asking "does A explain B?" — it changes what you fix and what you can drop.

## The flow

1. **Understand.** Restate the goal in a line. A surface change is defined by the **vendored** OpenAPI snapshots in `spec/v1.openapi.yaml` / `spec/v2.openapi.yaml` — read them before designing, and never import upstream at build time. Only client-facing, API-key-authenticated endpoints are in scope; no internal admin or session-cookie controllers. If the ask cites a URL, `WebFetch` it and extract the *mechanism*, then translate it onto this stack.

2. **Distrust the paperwork.** `docs/architecture/*` and CLAUDE.md are the design intent, not a mirror of the tree — CLAUDE.md still describes this repo as greenfield, while all seven SDKs have shipped 0.2.x. Before planning work off a doc, check it against the code and `git log`; merged PR titles and the `<lang>-v<version>` tags are the cheapest ground truth. State plainly which claims you falsified, so nobody re-implements shipped work or "fixes" working code.

3. **Establish the ground truth for the defect.** There is no production to query here — the SDK's equivalents are: the vendored spec (does the API really return that field, with that casing?), `spec/error_codes.txt` and `spec/parity.yaml` (is the symbol supposed to exist in all seven?), `python3 bin/check_parity.py` (is it already failing?), and the CHANGELOGs plus release tags (which version does the reporter actually have?). A finding pinned to a spec line or a parity failure outranks one derived from reading a single language alone.

4. **Explore (parallel).** Fan out `Agent` Explore agents (thoroughness "very thorough") to map the affected surface: the resource module in each language, the transport layer that owns the cross-cutting behaviour, the typed-error taxonomy, and the tests to mirror (`file:line`). Give each a **disjoint** area — one language each, or one architectural theme each — so their reports don't overlap. Require of every finding: severity, `file:line`, a one-sentence defect statement, and a **concrete failure scenario** (inputs → wrong outcome). Demand two more things explicitly: the doc claims they **falsified**, and the premises in your brief that turned out **true**. Produce a ranked worklist; log what the survey could not cover.

   **Protect your own context.** You are the only one who must survive to the merge. Do not read what an agent will report; do not re-derive a conclusion you already have.

5. **Fold in live user reports as first-class findings.** Mid-run the user may paste a stack trace from a consuming app, an issue thread, or a failing snippet. That is *confirmed against a published package* and routinely outranks the sweep's own findings — reproduce it against the installed version, root-cause it, and rank it above equal-severity read-only findings. If an in-flight agent already owns that language, extend its brief with `SendMessage` rather than spawning a second agent onto the same paths.

6. **Build — branch first, then fan out.** Before a single agent starts, get off `main`:

   ```bash
   git fetch origin && git status --short   # expect a clean tree
   git checkout -b <type>/<slug>            # fix/ feat/ test/ refactor/ docs/
   ```
   Do it now, while the tree is clean. By commit time the tree is dirty enough that you will not want to think about branches.

   Then **land one language first, settle the shape, then port**. Two agents that must edit one file are one slice, not two. When a change is really about cross-cutting behaviour (retry, caching, pagination, idempotency, request-id), fix it in that language's **transport** first and make every resource client adopt it — never re-implement it per resource.

   Every brief carries all nine of these — omitting one is how a run goes wrong:
   - **its exclusive file set** (its language directory), and never edit outside it — especially not `spec/**` or another language;
   - **which other agents are live on which paths**, so a collision is *reported*, not silently resolved;
   - each finding with `file:line`, the defect and the concrete failure scenario — plus permission to **drop any finding the code contradicts**;
   - **evidence first, diagnosis second** — the symptom, the spec line, the failing input; *then* your hypothesis, explicitly labelled unverified, to confirm or kill before building. Briefs that lead with a confident root cause send agents to the wrong file;
   - **the house constraints binding its area**: zero runtime deps (stdlib only), nothing newer than the language's min-version floor, typed model objects not raw maps/hashes/dicts, one error class per `error_code` with the full payload (`error_code`, `message`, `http_status`, `request_id`, `retry_after`, `response_body`, redacted `request_summary`), cross-cutting concerns in the transport only, **back-compat is permanent** (never remove or rename a public symbol, v1 stays shipped forever), no catch-all `rescue`/`catch`, no safe-navigation hiding nil, files under ~500 LOC;
   - **tests ship with the code, failure case first** — mocked-HTTP unit tests per resource client and cassette replay for integration, never the live API. Coverage parity matters: a missing PHP test is as visible as a missing TS test;
   - **checks narrowed to its OWN language** (see the table). Never the whole repo, never `bin/check_parity.py`;
   - **no git operations at all** — no branch, commit, checkout or stash. The coordinator owns all git; work is left uncommitted;
   - **never tell an agent to "ask me" — it cannot.** A subagent has no channel to the user, so a question is a dead end: it blocks or it guesses. Give it the two legal moves: **decide and flag it** (act on the most defensible reading, state the assumption in its report, mark the artifact so you can overwrite it), or **stop and report** with the evidence when either path would be unsafe or wasted. Then *you* take the question to the user and re-task with `SendMessage`, which resumes the agent with its full context.

7. **Verify, PR, merge.** Run the touched languages' suites and `python3 bin/check_parity.py` **once**, in the background. **Before committing, sweep the agents' leftovers**: scratch scripts, debug prints, stray probe files at the repo root, a `.venv` or `node_modules` an agent created outside its package. They must not ship.

   **Let every agent finish, then plain git.** Do not commit while agents are still writing — that is the only thing that ever made this complicated.

   ```bash
   git fetch origin                    # did main move? if so, see below
   git add <this slice's paths>        # never -A blindly
   git status --short                  # then READ it — strip scratch files and debug output
   git commit && git push -u origin HEAD
   ```
   Conventional Commits scoped to the language (`feat(ruby): …`). Naming paths on `git add` is all the selectivity you need — **never `git stash`** (one global stack shared with every concurrent agent).

   **Main moves under you.** Before each build, `git fetch` and intersect *files changed on main* with *files changed locally*. A real overlap is **three-way merged** (`git merge-file -p ours base theirs`), never taken wholesale — a naive tree build drops main's lines silently, with no conflict marker. Version files and CHANGELOGs are the usual collision sites.

   Then `claudetm merge-pr <PR>` (background) — it waits for CI, fixes failures, addresses review comments including CodeRabbit's, and merges when green. It operates on the **current directory**, so at most one PR is in flight: parallel *building* is fine, parallel *merging* is not. **When every check already passes, prefer `gh pr merge --squash`** — `claudetm` can hang on an already-green PR. Gotcha: **0 registered checks reads as "pass"**, and this repo's `paths:` filters mean an untouched language legitimately registers nothing — so confirm the checks you *expect* for the paths you changed are present and green, rather than trusting an empty list.

8. **Release — the step people forget.** Bumping the language's version source file is what triggers a release; no human pushes tags. Patch is per-language; **minor and major land across all seven in lockstep**, gated by `parity-check.yml`. Initial surfaces are `0.1.x`/`0.2.x` — pre-1.0, so the surface may still evolve, but back-compat within a shipped major still binds. After the merge, confirm for each language: the `release` job ran, it published, and the tag `<lang>-v<version>` now exists (`gh run list --branch main`, `git ls-remote --tags origin`). The job is idempotent, which is exactly why a missing version bump fails **silently** — no error, no package. Go's release additionally verifies a CHANGELOG entry exists; write it with the change, not after.

9. **Leave the trail straight.** A public-surface change updates the matching page at `../tesote.com` (`www.tesote.com/docs/sdk`) — link, don't duplicate; end-user docs do not live here. Update `docs/architecture/*` when you changed a rule it states, `spec/parity.yaml` when the shared symbol set moved, and CLAUDE.md when a claim in it became false.

## Hard rules (from CLAUDE.md — non-negotiable)

**Zero runtime dependencies** — stdlib only (the single exception: `jackson-databind` in Java). Nothing newer than each language's min-version floor; no experimental features. Dev/test deps loose-pinned, never `=x.y.z`. **No cross-language code sharing** — duplicate idiomatically. Transport separate from resource clients; the transport owns pagination, retry with backoff + jitter, rate-limit awareness, caching, idempotency keys and request-id propagation. Typed models, never raw maps. One error class per `error_code`; typed `NetworkError`/`TimeoutError` for transport failures; **never** bubble the underlying language exception, and never collapse into a single `ApiError`. Never log or echo the bearer token — `request_summary` is redacted. **Back-compat is permanent**: removing or renaming a public symbol in any shipped version is a break. Semver: patch per-language, minor/major in lockstep across all seven. Files under ~500 LOC. No catch-all handlers, no safe-navigation hiding nil. Never the live API in tests. Never `--force`, `--no-verify` or `reset --hard` without permission. **Never `git stash`** (one global stack, shared with every concurrent agent).

## Output

Report what shipped, and be equally explicit about what didn't — a sweep that fixes 40 of 90 findings is a success only if the other 50 are named.

```
Root cause:  <the one-line mechanism, for a bug sweep>
Languages:   <n> of 7 touched → ts, python, …
Fixed:       <n> findings across <m> PRs → #… #…
Deferred:    <n> — <what, and why not now>               [never omit this line]
Falsified:   <doc/CLAUDE.md claims that were wrong, now corrected>
Released:    <lang>-v<version> …   parity: <green | lockstep pending>
Verified:    <release jobs seen green / tags present>   spec: <snapshot in sync?>
Docs:        <../tesote.com PR, docs/architecture updates, or none>
```
