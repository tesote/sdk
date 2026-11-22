# tesote/sdk — packages/csharp

> Workspace brain — established by the developerz.ai maintainer agent for the maintainers of this repository from this workspace, then owned by
> its maintainers. Correct it in place; setup never overwrites a brain that exists.

C# .NET SDK solution, independently tested and published to NuGet.

The gate is the repo root’s: `bin/check` (see the root `CLAUDE.md`).

## Structure

- `src/Tesote.Sdk.csproj` is the net8.0 library and `tests/Tesote.Sdk.Tests.csproj` references it. `packages/csharp/Tesote.Sdk.sln:5-8`, `packages/csharp/tests/Tesote.Sdk.Tests.csproj:22-24`
- Source separates errors, models, transport, and v1/v2 clients. `packages/csharp` file inventory.

## Conventions

- Target net8.0 with nullable references, generated docs, and warnings as errors. `packages/csharp/src/Tesote.Sdk.csproj:3-14`
- Use four-space indentation and sort System imports first. `packages/csharp/.editorconfig:3-23`
- Run `dotnet restore`, `dotnet build --configuration Release --no-restore`, and `dotnet test --configuration Release --no-build --logger "console;verbosity=normal"`. `.github/workflows/csharp.yml:44-53`
- Use xUnit and WireMock.Net through the test project. `packages/csharp/tests/Tesote.Sdk.Tests.csproj:15-24`
