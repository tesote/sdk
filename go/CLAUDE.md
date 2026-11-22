# tesote/sdk — go

> Workspace brain — established by the developerz.ai maintainer agent for the maintainers of this repository from this workspace, then owned by
> its maintainers. Correct it in place; setup never overwrites a brain that exists.

Go SDK module for github.com/tesote/sdk/go, independently tested and released by Go submodule tag.

The gate is the repo root’s: `bin/check` (see the root `CLAUDE.md`).

## Structure

- `v1` and `v2` contain versioned resource clients; root package contains transport, errors, options, and models. `go/go.mod:1-3`, `.github/workflows/go.yml:3-7`
- Go transport tests use `httptest`. `go/transport_test.go:3-17`

## Conventions

- Keep module compatibility at Go 1.21. `go/go.mod:1-3`
- Run gofmt verification, `go vet ./...`, `go test -race -count=1 ./...`, and golangci-lint. `.github/workflows/go.yml:47-57`
- Test HTTP behavior with `httptest` servers and deterministic backoff helpers. `go/transport_test.go:15-32`
- Preserve the `go/vX.Y.Z` release tag scheme. `docs/architecture/release.md:21-29`
