using System.Threading.Tasks;
using Tesote.Sdk.Errors;
using WireMock.RequestBuilders;
using WireMock.ResponseBuilders;
using WireMock.Server;
using Xunit;

namespace Tesote.Sdk.Tests;

public sealed class AccountsTests : System.IDisposable
{
    private readonly WireMockServer _server;

    public AccountsTests()
    {
        _server = WireMockServer.Start();
    }

    public void Dispose()
    {
        _server.Stop();
        _server.Dispose();
    }

    private const string AccountJson =
        "{\"id\":\"acc_1\",\"name\":\"Checking\"," +
        "\"data\":{\"masked_account_number\":\"****1234\",\"currency\":\"VES\"," +
        "\"transactions_data_current_as_of\":null,\"balance_data_current_as_of\":null," +
        "\"custom_user_provided_identifier\":null,\"balance_cents\":\"1000\"," +
        "\"available_balance_cents\":\"950\"}," +
        "\"bank\":{\"name\":\"Banesco\"}," +
        "\"legal_entity\":{\"id\":null,\"legal_name\":null}," +
        "\"tesote_created_at\":\"2026-01-01T00:00:00Z\"," +
        "\"tesote_updated_at\":\"2026-04-01T00:00:00Z\"}";

    [Fact]
    public async Task V1ListReturnsTypedAccounts()
    {
        _server
            .Given(Request.Create().WithPath("/api/v1/accounts").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithBody("{\"total\":1,\"accounts\":[" + AccountJson + "]," +
                          "\"pagination\":{\"current_page\":1,\"per_page\":50," +
                          "\"total_pages\":1,\"total_count\":1}}"));

        using var client = TestHelpers.NewV1(_server.Url + "/api");
        var result = await client.Accounts.ListAsync(page: 1, perPage: 50);
        Assert.Equal(1, result.Total);
        Assert.Equal("acc_1", result.Accounts[0].Id);
        Assert.Equal("Banesco", result.Accounts[0].Bank.Name);
        Assert.Equal("1000", result.Accounts[0].Data.BalanceCents);
    }

    [Fact]
    public async Task V2GetReturnsTypedAccount()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200).WithBody(AccountJson));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var account = await client.Accounts.GetAsync("acc_1");
        Assert.Equal("acc_1", account.Id);
        Assert.Equal("Checking", account.Name);
    }

    [Fact]
    public async Task V2GetMaps404ToAccountNotFound()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/missing").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(404)
                .WithBody("{\"error\":\"missing\",\"error_code\":\"ACCOUNT_NOT_FOUND\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<AccountNotFoundException>(() => client.Accounts.GetAsync("missing"));
    }

    [Fact]
    public async Task V2WorkspaceSuspendedMapsTo403()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(403)
                .WithBody("{\"error\":\"suspended\",\"error_code\":\"WORKSPACE_SUSPENDED\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<WorkspaceSuspendedException>(() => client.Accounts.ListAsync());
    }

    [Fact]
    public async Task V2SyncReturnsSessionAndCarriesIdempotencyKey()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/sync").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(202)
                .WithBody("{\"message\":\"Sync started\",\"sync_session_id\":\"ss_1\"," +
                          "\"status\":\"pending\",\"started_at\":\"2026-04-01T00:00:00Z\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var resp = await client.Accounts.SyncAsync("acc_1", idempotencyKey: "abc-123");
        Assert.Equal("ss_1", resp.SyncSessionId);
        Assert.Equal("pending", resp.Status);

        var entry = Assert.Single(_server.LogEntries);
        Assert.Equal("abc-123", System.Linq.Enumerable.First(entry.RequestMessage.Headers!["Idempotency-Key"]));
    }

    [Fact]
    public async Task V2SyncMaps409ToSyncInProgress()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/sync").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(409)
                .WithBody("{\"error\":\"in progress\",\"error_code\":\"SYNC_IN_PROGRESS\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<SyncInProgressException>(() => client.Accounts.SyncAsync("acc_1"));
    }
}
