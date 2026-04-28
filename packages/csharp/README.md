# tesote-sdk (C# / .NET)

Official C# SDK for the [equipo.tesote.com](https://equipo.tesote.com) API.

- Package: `Tesote.Sdk` (NuGet, `0.1.0`)
- Min .NET: **net8.0** (LTS)
- Tested on: net8.0, net9.0

## Install

```bash
dotnet add package Tesote.Sdk
```

Or in `.csproj`:

```xml
<ItemGroup>
  <PackageReference Include="Tesote.Sdk" Version="0.1.0" />
</ItemGroup>
```

## Usage

```csharp
using Tesote.Sdk;
using Tesote.Sdk.V2;

await using var client = new V2Client(new ClientOptions
{
    ApiKey = Environment.GetEnvironmentVariable("TESOTE_API_KEY"),
});

var accounts = await client.Accounts.ListAsync();
var one = await client.Accounts.GetAsync("acct_42");
```

Versioned clients live side-by-side; consumers pick a version:

```csharp
using Tesote.Sdk.V1;
using Tesote.Sdk.V2;
```

## Configuration

```csharp
using Tesote.Sdk;
using Tesote.Sdk.Internal;
using Tesote.Sdk.V2;

var client = new V2Client(new ClientOptions
{
    ApiKey = apiKey,
    BaseUrl = "https://equipo-staging.tesote.com/api",
    UserAgent = "MyApp/1.0",
    RequestTimeout = TimeSpan.FromSeconds(30),
    RetryPolicy = new RetryPolicy(3,
        TimeSpan.FromMilliseconds(250),
        TimeSpan.FromSeconds(8),
        RetryOnNetwork: true),
    CacheBackend = new InMemoryCacheBackend(),
    Logger = ev => Console.WriteLine(ev),
});
```

If `ApiKey` is omitted the SDK falls back to the `TESOTE_SDK_API_KEY`
environment variable; if still missing, construction throws
`ConfigException`. `BaseUrl` follows the same pattern via
`TESOTE_SDK_API_URL`.

## Errors

Every error is a typed subclass of `TesoteException`. Pattern:

```csharp
using Tesote.Sdk.Errors;

try
{
    await client.Accounts.ListAsync();
}
catch (RateLimitExceededException ex)
{
    await Task.Delay(TimeSpan.FromSeconds(ex.RetryAfter ?? 1));
}
catch (UnauthorizedException)
{
    // rotate key
}
catch (TesoteException ex)
{
    // catch-all last resort; full context on ex.ToString()
}
```

Every error carries: `ErrorCode`, `HttpStatus`, `RequestId`, `ErrorId`,
`RetryAfter`, `ResponseBody`, `RequestSummary`, `Attempts`. The bearer
token is never included; it's redacted to `Bearer ****<last4>`.

See [`docs/architecture/errors.md`](../../docs/architecture/errors.md) for
the full taxonomy.

## Polling model

The platform is poll-based for v1/v2 — there are no server push
notifications.

## Runtime dependencies

**Zero.** HTTP, JSON, retries, caching, and concurrency primitives all use
the .NET standard library (`System.Net.Http`, `System.Text.Json`,
`System.Buffers`). Customers don't get any transitive deps from this SDK.

## Development

```bash
cd packages/csharp
dotnet restore
dotnet build --configuration Release
dotnet test --configuration Release
```

## Release

Pushing a bump of the `<Version>` in `src/Tesote.Sdk.csproj` to `main`
triggers `.github/workflows/csharp.yml` which packs and pushes via
`dotnet nuget push` to nuget.org. See
[`docs/architecture/release.md`](../../docs/architecture/release.md).
