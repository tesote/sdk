using System.Linq;
using System.Threading.Tasks;
using Tesote.Sdk.Errors;
using WireMock.RequestBuilders;
using WireMock.ResponseBuilders;
using WireMock.Server;
using Xunit;

namespace Tesote.Sdk.Tests;

public sealed class SyncSessionsTests : System.IDisposable
{
    private readonly WireMockServer _server;

    public SyncSessionsTests()
    {
        _server = WireMockServer.Start();
    }

    public void Dispose()
    {
        _server.Stop();
        _server.Dispose();
    }

    [Fact]
    public async Task ListReturnsTypedSessionsWithPagination()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/sync_sessions").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithBody("{\"sync_sessions\":[" +
                          "{\"id\":\"ss_1\",\"status\":\"completed\"," +
                          "\"started_at\":\"2026-04-01T00:00:00Z\",\"completed_at\":\"2026-04-01T00:01:00Z\"," +
                          "\"transactions_synced\":10,\"accounts_count\":1,\"error\":null,\"performance\":null}]," +
                          "\"limit\":50,\"offset\":0,\"has_more\":false}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var result = await client.SyncSessions.ListAsync("acc_1", limit: 50, offset: 0, status: "completed");
        Assert.Single(result.SyncSessions);
        Assert.Equal("ss_1", result.SyncSessions[0].Id);
        Assert.False(result.HasMore);

        var entry = Assert.Single(_server.LogEntries);
        Assert.Contains("status=completed", entry.RequestMessage.Url);
    }

    [Fact]
    public async Task GetMaps404ToSyncSessionNotFound()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/sync_sessions/missing").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(404)
                .WithBody("{\"error\":\"not found\",\"error_code\":\"SYNC_SESSION_NOT_FOUND\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<SyncSessionNotFoundException>(
            () => client.SyncSessions.GetAsync("acc_1", "missing"));
    }
}
