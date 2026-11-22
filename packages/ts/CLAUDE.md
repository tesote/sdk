# tesote/sdk — packages/ts

> Workspace brain — established by the developerz.ai maintainer agent for the maintainers of this repository from this workspace, then owned by
> its maintainers. Correct it in place; setup never overwrites a brain that exists.

TypeScript SDK package, independently tested and published to npm.

The gate is the repo root’s: `bin/check` (see the root `CLAUDE.md`).

## Structure

- `src` holds the SDK implementation; `test` holds Vitest coverage, both included by TypeScript configuration. `packages/ts/tsconfig.json:2-27`
- `package.json` defines build, typecheck, Biome lint, and Vitest scripts. `packages/ts/package.json:31-38`

## Conventions

- Use strict TypeScript settings. `packages/ts/tsconfig.json:7-17`
- Follow Biome: two spaces, 100 columns, single quotes, semicolons, trailing commas, no explicit `any`. `packages/ts/biome.json:8-34`
- Run `npm ci`, `npm run typecheck`, `npm run lint`, and `npm test`. `.github/workflows/ts.yml:42-51`
- Test transport behavior using mocked `fetch`. `packages/ts/test/transport.test.ts:15-27`
