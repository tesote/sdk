using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json.Nodes;
using System.Threading.Tasks;
using Tesote.Sdk;
using Tesote.Sdk.Errors;
using Tesote.Sdk.Internal;
using WireMock.RequestBuilders;
using WireMock.ResponseBuilders;
using WireMock.Server;
using Xunit;

namespace Tesote.Sdk.Tests;

public sealed class TransportTests : IDisposable
{
    private readonly WireMockServer _server;

    public TransportTests()
    {
        _server = WireMockServer.Start();
    }

    public void Dispose()
    {
        _server.Stop();
        _server.Dispose();
    }

    private Transport NewTransport(RetryPolicy? policy = null, ICacheBackend? cache = null)
    {
        return new Transport(new ClientOptions
        {
            ApiKey = "sk_test_abcd1234",
            BaseUrl = _server.Url + "/api",
            RequestTimeout = TimeSpan.FromSeconds(2),
            RetryPolicy = policy,
            CacheBackend = cache,
        });
    }

    [Fact]
    public void MissingApiKeyThrowsConfigException()
    {
        Environment.SetEnvironmentVariable("TESOTE_SDK_API_KEY", null);
        var ex = Assert.Throws<ConfigException>(() => new Transport(new ClientOptions { ApiKey = "" }));
        Assert.Equal("CONFIG", ex.ErrorCode);
    }

    [Fact]
    public async Task SuccessfulGetReturnsParsedJsonAndInjectsBearer()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .RespondWith(Response.Create()
                .WithStatusCode(200)
                .WithHeader("Content-Type", "application/json")
                .WithHeader("X-Request-Id", "req_123")
                .WithBody("{\"data\":[{\"id\":\"acct_1\"}]}"));

        using var t = NewTransport();
        var result = await t.RequestAsync(RequestOptions.Get("/v3/accounts"));

        Assert.NotNull(result);
        Assert.Equal("acct_1", result!["data"]![0]!["id"]!.GetValue<string>());

        var entries = _server.LogEntries.ToList();
        Assert.Single(entries);
        var headers = entries[0].RequestMessage.Headers!;
        Assert.Equal("Bearer sk_test_abcd1234", headers["Authorization"].First());
        Assert.Equal("application/json", headers["Accept"].First());
        Assert.StartsWith("tesote-sdk-csharp/", headers["User-Agent"].First());
    }

    [Fact]
    public async Task RateLimitHeadersCapturedIntoSnapshot()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .RespondWith(Response.Create()
                .WithStatusCode(200)
                .WithHeader("X-RateLimit-Limit", "200")
                .WithHeader("X-RateLimit-Remaining", "199")
                .WithHeader("X-RateLimit-Reset", "1700000000")
                .WithBody("{}"));

        using var t = NewTransport();
        await t.RequestAsync(RequestOptions.Get("/v3/accounts"));

        var snap = t.LastRateLimit;
        Assert.Equal(200, snap.Limit);
        Assert.Equal(199, snap.Remaining);
        Assert.NotNull(snap.ResetAt);
        Assert.Equal(1700000000L, snap.ResetAt!.Value.ToUnixTimeSeconds());
    }

    [Fact]
    public async Task RetriesOn503WithBackoffThenSucceeds()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .InScenario("retry")
            .WillSetStateTo("after-1")
            .RespondWith(Response.Create().WithStatusCode(503).WithBody("{}"));

        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .InScenario("retry")
            .WhenStateIs("after-1")
            .WillSetStateTo("after-2")
            .RespondWith(Response.Create().WithStatusCode(503).WithBody("{}"));

        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .InScenario("retry")
            .WhenStateIs("after-2")
            .RespondWith(Response.Create().WithStatusCode(200).WithBody("{\"ok\":true}"));

        using var t = NewTransport(new RetryPolicy(3, TimeSpan.FromMilliseconds(1), TimeSpan.FromMilliseconds(5), true));
        var result = await t.RequestAsync(RequestOptions.Get("/v3/accounts"));

        Assert.True(result!["ok"]!.GetValue<bool>());
        Assert.Equal(3, _server.LogEntries.Count());
    }

    [Fact]
    public async Task RetryExhaustionThrowsRateLimitWithAttempts()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .RespondWith(Response.Create()
                .WithStatusCode(429)
                .WithHeader("Retry-After", "0")
                .WithBody("{\"error\":\"rate limited\",\"error_code\":\"RATE_LIMIT_EXCEEDED\"}"));

        using var t = NewTransport(new RetryPolicy(3, TimeSpan.FromMilliseconds(1), TimeSpan.FromMilliseconds(5), true));
        var ex = await Assert.ThrowsAsync<RateLimitExceededException>(
            () => t.RequestAsync(RequestOptions.Get("/v3/accounts")));
        Assert.Equal(3, ex.Attempts);
        Assert.Equal("RATE_LIMIT_EXCEEDED", ex.ErrorCode);
    }

    [Fact]
    public async Task DoesNotRetryOn4xxOtherThan429()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .RespondWith(Response.Create()
                .WithStatusCode(401)
                .WithBody("{\"error\":\"bad key\",\"error_code\":\"UNAUTHORIZED\"}"));

        using var t = NewTransport();
        var ex = await Assert.ThrowsAsync<UnauthorizedException>(
            () => t.RequestAsync(RequestOptions.Get("/v3/accounts")));
        Assert.Single(_server.LogEntries);
        Assert.Equal(1, ex.Attempts);
    }

    [Fact]
    public async Task RequestIdAttachedToThrownException()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .RespondWith(Response.Create()
                .WithStatusCode(401)
                .WithHeader("X-Request-Id", "req_abc")
                .WithBody("{\"error\":\"bad key\",\"error_code\":\"UNAUTHORIZED\",\"error_id\":\"err_1\"}"));

        using var t = NewTransport();
        var ex = await Assert.ThrowsAsync<UnauthorizedException>(
            () => t.RequestAsync(RequestOptions.Get("/v3/accounts")));
        Assert.Equal("req_abc", ex.RequestId);
        Assert.Equal("err_1", ex.ErrorId);
    }

    [Fact]
    public async Task IdempotencyKeyAutoGeneratedOnPost()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts/acct_1/sync").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(200).WithBody("{}"));

        using var t = NewTransport();
        var opts = new RequestOptions
        {
            Method = "POST",
            Path = "/v3/accounts/acct_1/sync",
            Body = System.Text.Encoding.UTF8.GetBytes("{}"),
            BodyShape = "0 bytes",
        };
        await t.RequestAsync(opts);

        var entry = _server.LogEntries.Single();
        var key = entry.RequestMessage.Headers!["Idempotency-Key"].First();
        Assert.False(string.IsNullOrEmpty(key));
        // Throws if not a UUID:
        _ = Guid.Parse(key);
    }

    [Fact]
    public async Task IdempotencyKeyHonoredWhenProvided()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts/acct_1/sync").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(200).WithBody("{}"));

        using var t = NewTransport();
        var opts = new RequestOptions
        {
            Method = "POST",
            Path = "/v3/accounts/acct_1/sync",
            Body = System.Text.Encoding.UTF8.GetBytes("{}"),
            IdempotencyKey = "my-key",
        };
        await t.RequestAsync(opts);

        Assert.Equal("my-key", _server.LogEntries.Single().RequestMessage.Headers!["Idempotency-Key"].First());
    }

    [Fact]
    public async Task GetDoesNotCarryIdempotencyKey()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200).WithBody("{}"));

        using var t = NewTransport();
        await t.RequestAsync(RequestOptions.Get("/v3/accounts"));

        var headers = _server.LogEntries.Single().RequestMessage.Headers!;
        Assert.False(headers.ContainsKey("Idempotency-Key"));
    }

    [Fact]
    public async Task CacheHitAvoidsSecondRequest()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .RespondWith(Response.Create()
                .WithStatusCode(200)
                .WithHeader("Content-Type", "application/json")
                .WithBody("{\"hits\":1}"));

        var cache = new InMemoryCacheBackend();
        using var t = NewTransport(cache: cache);

        var opts = RequestOptions.Get("/v3/accounts");
        opts.CacheTtl = TimeSpan.FromSeconds(30);

        var first = await t.RequestAsync(opts);
        var second = await t.RequestAsync(opts);

        Assert.Single(_server.LogEntries);
        Assert.Equal(first!.ToJsonString(), second!.ToJsonString());
    }

    [Fact]
    public async Task RequestSummaryRedactsBearerToken()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .RespondWith(Response.Create()
                .WithStatusCode(401)
                .WithBody("{\"error\":\"x\",\"error_code\":\"UNAUTHORIZED\"}"));

        using var t = NewTransport();
        var ex = await Assert.ThrowsAsync<UnauthorizedException>(
            () => t.RequestAsync(RequestOptions.Get("/v3/accounts")));
        var redacted = ex.RequestSummary!.RedactedAuthorization;
        Assert.StartsWith("Bearer ****", redacted);
        Assert.Contains("1234", redacted);
        // why: original key must never appear verbatim.
        Assert.DoesNotContain("abcd1234", redacted);
    }

    [Fact]
    public async Task ApiExceptionMessageIsHumanGreppable()
    {
        _server
            .Given(Request.Create().WithPath("/api/v3/accounts").UsingGet())
            .RespondWith(Response.Create()
                .WithStatusCode(429)
                .WithHeader("Retry-After", "42")
                .WithHeader("X-Request-Id", "req_xyz")
                .WithBody("{\"error\":\"Rate limit exceeded\",\"error_code\":\"RATE_LIMIT_EXCEEDED\"}"));

        using var t = NewTransport(new RetryPolicy(1, TimeSpan.FromMilliseconds(1), TimeSpan.FromMilliseconds(2), true));
        var ex = await Assert.ThrowsAsync<RateLimitExceededException>(
            () => t.RequestAsync(RequestOptions.Get("/v3/accounts")));
        var msg = ex.ToString();
        Assert.Contains("RateLimitExceededException", msg);
        Assert.Contains("429", msg);
        Assert.Contains("RATE_LIMIT_EXCEEDED", msg);
        Assert.Contains("req_xyz", msg);
        Assert.Contains("retry_after: 42", msg);
    }

    [Fact]
    public async Task PutAndPatchAndDeleteAlsoGetIdempotencyKey()
    {
        _server
            .Given(Request.Create().WithPath("/api/x").UsingAnyMethod())
            .RespondWith(Response.Create().WithStatusCode(200).WithBody("{}"));

        using var t = NewTransport();
        foreach (var method in new[] { "PUT", "PATCH", "DELETE" })
        {
            var opts = new RequestOptions { Method = method, Path = "/x", Body = System.Text.Encoding.UTF8.GetBytes("{}") };
            await t.RequestAsync(opts);
        }
        var keys = _server.LogEntries
            .Select(e => e.RequestMessage.Headers!["Idempotency-Key"].First())
            .ToList();
        Assert.Equal(3, keys.Count);
        foreach (var k in keys)
        {
            _ = Guid.Parse(k);
        }
    }
}
