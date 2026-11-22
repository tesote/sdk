# tesote/sdk — packages/java

> Workspace brain — established by the developerz.ai maintainer agent for the maintainers of this repository from this workspace, then owned by
> its maintainers. Correct it in place; setup never overwrites a brain that exists.

Java SDK library, independently tested and published to Maven Central.

The gate is the repo root’s: `bin/check` (see the root `CLAUDE.md`).

## Structure

- `src/main/java/com/tesote/sdk` contains implementation and `src/test/java/com/tesote/sdk` contains JUnit tests. `packages/java/build.gradle.kts:1-4`, `packages/java/src/test/java/com/tesote/sdk/TransportTest.java:1-14`
- The Gradle wrapper is the package build entrypoint. `.github/workflows/java.yml:36-38`

## Conventions

- Compile against Java 17 with UTF-8, full linting, and warnings as errors. `packages/java/build.gradle.kts:15-21`, `packages/java/build.gradle.kts:51-55`
- Use JUnit Platform and MockWebServer for tests. `packages/java/build.gradle.kts:38-49`, `packages/java/src/test/java/com/tesote/sdk/TransportTest.java:27-39`
- Run `./gradlew test --no-daemon`. `.github/workflows/java.yml:31-38`
- Keep Jackson as the documented Java runtime dependency exception. `packages/java/build.gradle.kts:27-40`
