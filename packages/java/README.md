# tesote-sdk (Java)

Official Java SDK for the [equipo.tesote.com](https://equipo.tesote.com) API.

- Coordinates: `com.tesote:sdk:0.1.0`
- Min Java: **17** (LTS)
- Tested on: 17, 21

## Install

Gradle (Kotlin DSL):

```kotlin
dependencies {
    implementation("com.tesote:sdk:0.1.0")
}
```

Maven:

```xml
<dependency>
    <groupId>com.tesote</groupId>
    <artifactId>sdk</artifactId>
    <version>0.1.0</version>
</dependency>
```

## Usage

```java
import com.tesote.sdk.v2.V2Client;
import com.fasterxml.jackson.databind.JsonNode;

V2Client client = V2Client.builder()
    .apiKey(System.getenv("TESOTE_API_KEY"))
    .build();

JsonNode accounts = client.accounts().list();
JsonNode one = client.accounts().get("acct_42");
```

Versioned clients live side-by-side; consumers pick a version:

```java
import com.tesote.sdk.v1.V1Client;
import com.tesote.sdk.v2.V2Client;
```

## Configuration

```java
V2Client client = V2Client.builder()
    .apiKey(apiKey)
    .baseUrl("https://equipo-staging.tesote.com/api")
    .userAgent("MyApp/1.0")
    .requestTimeout(java.time.Duration.ofSeconds(30))
    .retryPolicy(new com.tesote.sdk.Transport.RetryPolicy(
        3, java.time.Duration.ofMillis(250),
        java.time.Duration.ofSeconds(8), true))
    .httpClientBuilder(java.net.http.HttpClient.newBuilder()
        .connectTimeout(java.time.Duration.ofSeconds(5)))
    .cacheBackend(new com.tesote.sdk.internal.InMemoryCacheBackend())
    .logger(event -> System.out.println(event))
    .build();
```

## Errors

Every error is a typed subclass of `TesoteException` (unchecked). Pattern:

```java
import com.tesote.sdk.errors.*;

try {
    client.accounts().list();
} catch (RateLimitExceededException e) {
    Thread.sleep(e.retryAfter() * 1000L);
} catch (UnauthorizedException e) {
    // rotate key
} catch (TesoteException e) {
    // catch-all last resort; full context on e.getMessage()
}
```

Every error carries: `errorCode`, `httpStatus`, `requestId`, `errorId`,
`retryAfter`, `responseBody`, `requestSummary`, `attempts`. The bearer
token is never included; it's redacted to `Bearer ****<last4>`.

See [`docs/architecture/errors.md`](../../docs/architecture/errors.md) for the
full taxonomy.

## Polling model

The platform is poll-based for v1/v2 — there are no server push notifications.

## Runtime dependencies

The SDK ships with **one** runtime dependency:

- `com.fasterxml.jackson.core:jackson-databind` — JSON parsing.

`jakarta.json` was the preferred choice (truly stdlib-adjacent), but its
pull-style API made the dynamic error envelope and response shapes (`JsonNode`
trees with optional fields) substantially more code for callers and the
transport. Jackson's `JsonNode` lets resource clients return raw trees today
and migrate to typed records in subsequent releases without breaking the
public surface.

HTTP, retries, caching, and concurrency primitives use the JDK standard
library (`java.net.http.HttpClient`, `java.util.concurrent`).

## Development

```bash
cd packages/java
./gradlew check
```

Wrapper is committed; no system Gradle install required.

## Release

Tag-driven from the monorepo: pushing `java-v0.1.0` to `main` triggers
`.github/workflows/java.yml` which signs and publishes via the Sonatype
Central Portal (nmcp plugin). See
[`docs/architecture/release.md`](../../docs/architecture/release.md).
