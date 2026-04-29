using System.Linq;
using System.Threading.Tasks;
using Tesote.Sdk.Errors;
using Tesote.Sdk.V2;
using WireMock.RequestBuilders;
using WireMock.ResponseBuilders;
using WireMock.Server;
using Xunit;

namespace Tesote.Sdk.Tests;

public sealed class TransactionsTests : System.IDisposable
{
    private readonly WireMockServer _server;

    public TransactionsTests()
    {
        _server = WireMockServer.Start();
    }

    public void Dispose()
    {
        _server.Stop();
        _server.Dispose();
    }

    private const string TxJson =
        "{\"id\":\"tx_1\",\"status\":\"posted\"," +
        "\"data\":{\"amount_cents\":1000,\"currency\":\"VES\",\"description\":\"hi\"," +
        "\"transaction_date\":\"2026-04-01\",\"created_at\":null,\"created_at_date\":null," +
        "\"note\":null,\"external_service_id\":null,\"running_balance_cents\":null}," +
        "\"tesote_imported_at\":\"2026-04-01T00:00:00Z\"," +
        "\"tesote_updated_at\":\"2026-04-01T00:00:00Z\"," +
        "\"transaction_categories\":[]," +
        "\"counterparty\":{\"name\":\"Acme\"}}";

    [Fact]
    public async Task V1ListTransactionsReturnsCursorPagination()
    {
        _server
            .Given(Request.Create().WithPath("/api/v1/accounts/acc_1/transactions").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithBody("{\"total\":1,\"transactions\":[" + TxJson + "]," +
                          "\"pagination\":{\"has_more\":false,\"per_page\":50,\"after_id\":\"tx_1\",\"before_id\":\"tx_1\"}}"));

        using var client = TestHelpers.NewV1(_server.Url + "/api");
        var result = await client.Accounts.ListTransactionsAsync("acc_1", perPage: 50);
        Assert.Equal(1, result.Total);
        Assert.Equal("tx_1", result.Transactions[0].Id);
        Assert.Equal(1000, result.Transactions[0].Data.AmountCents);
        Assert.False(result.Pagination.HasMore);
    }

    [Fact]
    public async Task V2ListMaps422ToInvalidDateRange()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transactions").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(422)
                .WithBody("{\"error\":\"bad date range\",\"error_code\":\"INVALID_DATE_RANGE\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var ex = await Assert.ThrowsAsync<InvalidDateRangeException>(
            () => client.Transactions.ListAsync("acc_1"));
        Assert.IsAssignableFrom<UnprocessableContentException>(ex);
    }

    [Fact]
    public async Task V2GetByIdReturnsV1Schema()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/transactions/tx_1").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200).WithBody(TxJson));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var tx = await client.Transactions.GetAsync("tx_1");
        Assert.Equal("tx_1", tx.Id);
        Assert.Equal("posted", tx.Status);
    }

    [Fact]
    public async Task V2GetByIdMapsTransactionNotFound()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/transactions/missing").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(404)
                .WithBody("{\"error\":\"missing\",\"error_code\":\"TRANSACTION_NOT_FOUND\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<TransactionNotFoundException>(() => client.Transactions.GetAsync("missing"));
    }

    [Fact]
    public async Task V2SyncSendsBodyAndReturnsResult()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transactions/sync").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithBody("{\"added\":[],\"modified\":[],\"removed\":[]," +
                          "\"next_cursor\":\"c2\",\"has_more\":false}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var req = new TransactionsClient.SyncRequest
        {
            Count = 100,
            Cursor = "now",
            Options = new TransactionsClient.SyncOptions { IncludeRunningBalance = true },
        };
        var result = await client.Transactions.SyncAsync("acc_1", req);
        Assert.Equal("c2", result.NextCursor);
        Assert.False(result.HasMore);

        var entry = Assert.Single(_server.LogEntries);
        var body = entry.RequestMessage.Body!;
        Assert.Contains("\"count\":100", body);
        Assert.Contains("\"cursor\":\"now\"", body);
        Assert.Contains("\"include_running_balance\":true", body);
    }

    [Fact]
    public async Task V2SyncMaps422ToInvalidCount()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transactions/sync").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(422)
                .WithBody("{\"error\":\"bad count\",\"error_code\":\"INVALID_COUNT\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<InvalidCountException>(() =>
            client.Transactions.SyncAsync("acc_1", new TransactionsClient.SyncRequest { Count = 9999 }));
    }

    [Fact]
    public async Task V2SyncMaps403HistorySyncForbidden()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transactions/sync").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(403)
                .WithBody("{\"error\":\"too old\",\"error_code\":\"HISTORY_SYNC_FORBIDDEN\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<HistorySyncForbiddenException>(() =>
            client.Transactions.SyncAsync("acc_1", new TransactionsClient.SyncRequest()));
    }

    [Fact]
    public async Task V2BulkValidatesBodyAndDeserializes()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/transactions/bulk").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithBody("{\"bulk_results\":[{\"account_id\":\"acc_1\",\"transactions\":[" + TxJson + "]," +
                          "\"pagination\":{\"has_more\":false,\"per_page\":50,\"after_id\":\"tx_1\",\"before_id\":\"tx_1\"}}]}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var result = await client.Transactions.BulkAsync(new[] { "acc_1", "acc_2" });
        Assert.Single(result.BulkResults);
        Assert.Equal("acc_1", result.BulkResults[0].AccountId);
    }

    [Fact]
    public async Task V2SearchSendsQueryAndDeserializes()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/transactions/search").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithBody("{\"transactions\":[" + TxJson + "],\"total\":1}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var result = await client.Transactions.SearchAsync("hi", limit: 10);
        Assert.Equal(1, result.Total);
        var entry = Assert.Single(_server.LogEntries);
        Assert.Contains("q=hi", entry.RequestMessage.Url);
        Assert.Contains("limit=10", entry.RequestMessage.Url);
    }

    [Fact]
    public async Task V2ExportReturnsRawCsvBytes()
    {
        const string csv = "Transaction ID,Date\ntx_1,2026-04-01\n";
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transactions/export").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithHeader("Content-Type", "text/csv; charset=utf-8")
                .WithHeader("Content-Disposition", "attachment; filename=\"transactions_acc_1_2026-04-01.csv\"")
                .WithBody(csv));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var raw = await client.Transactions.ExportAsync("acc_1", format: "csv");
        Assert.Contains("text/csv", raw.ContentType);
        Assert.Contains("transactions_acc_1", raw.ContentDisposition);
        Assert.Equal(csv, System.Text.Encoding.UTF8.GetString(raw.Body));
    }

    [Fact]
    public async Task UnsupportedMediaType415ReturnedAsApiException()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/transactions/bulk").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(415)
                .WithBody("{\"error\":\"need json\",\"error_code\":\"UNSUPPORTED_MEDIA_TYPE\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var ex = await Assert.ThrowsAsync<ApiException>(() =>
            client.Transactions.BulkAsync(new[] { "acc_1" }));
        Assert.Equal(415, ex.HttpStatus);

        // why: SDK still sends Content-Type: application/json on POST bodies; we verify here.
        var entry = Assert.Single(_server.LogEntries);
        var contentType = entry.RequestMessage.Headers!["Content-Type"].First();
        Assert.Contains("application/json", contentType);
    }
}
