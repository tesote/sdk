using System.Threading.Tasks;
using Tesote.Sdk.Errors;
using Tesote.Sdk.V2;
using WireMock.RequestBuilders;
using WireMock.ResponseBuilders;
using WireMock.Server;
using Xunit;

namespace Tesote.Sdk.Tests;

public sealed class BatchesTests : System.IDisposable
{
    private readonly WireMockServer _server;

    public BatchesTests()
    {
        _server = WireMockServer.Start();
    }

    public void Dispose()
    {
        _server.Stop();
        _server.Dispose();
    }

    [Fact]
    public async Task CreateReturnsOrdersAndBatchId()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/batches").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(201)
                .WithBody("{\"batch_id\":\"b_1\",\"orders\":[],\"errors\":[]}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var orders = new[]
        {
            new BatchesClient.BatchOrderInput
            {
                DestinationPaymentMethodId = "pm_d",
                Amount = "10.00",
                Currency = "VES",
                Description = "fee",
            },
        };
        var resp = await client.Batches.CreateAsync("acc_1", orders);
        Assert.Equal("b_1", resp.BatchId);
        Assert.Empty(resp.Errors);
    }

    [Fact]
    public async Task CreateMapsBatchValidationError()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/batches").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(400)
                .WithBody("{\"error\":\"bad batch\",\"error_code\":\"BATCH_VALIDATION_ERROR\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var ex = await Assert.ThrowsAsync<BatchValidationException>(() =>
            client.Batches.CreateAsync("acc_1", System.Array.Empty<BatchesClient.BatchOrderInput>()));
        // why: subclass of ValidationException so callers can catch the parent.
        Assert.IsAssignableFrom<ValidationException>(ex);
    }

    [Fact]
    public async Task GetReturnsSummary()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/batches/b_1").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithBody("{\"batch_id\":\"b_1\",\"total_orders\":3,\"total_amount_cents\":3000," +
                          "\"amount_currency\":\"VES\"," +
                          "\"statuses\":{\"draft\":3,\"pending_approval\":0,\"approved\":0," +
                          "\"processing\":0,\"completed\":0,\"failed\":0,\"cancelled\":0}," +
                          "\"batch_status\":\"draft\",\"created_at\":\"2026-04-01T00:00:00Z\"," +
                          "\"orders\":[]}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var summary = await client.Batches.GetAsync("acc_1", "b_1");
        Assert.Equal("b_1", summary.BatchId);
        Assert.Equal(3, summary.Statuses.Draft);
    }

    [Fact]
    public async Task GetMaps404ToBatchNotFound()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/batches/missing").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(404)
                .WithBody("{\"error\":\"missing\",\"error_code\":\"BATCH_NOT_FOUND\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<BatchNotFoundException>(() => client.Batches.GetAsync("acc_1", "missing"));
    }

    [Fact]
    public async Task ApproveReturnsCounts()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/batches/b_1/approve").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(200).WithBody("{\"approved\":5,\"failed\":0}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var resp = await client.Batches.ApproveAsync("acc_1", "b_1");
        Assert.Equal(5, resp.Approved);
    }

    [Fact]
    public async Task SubmitReturnsCounts()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/batches/b_1/submit").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(200).WithBody("{\"enqueued\":4,\"failed\":1}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var resp = await client.Batches.SubmitAsync("acc_1", "b_1", token: "otp");
        Assert.Equal(4, resp.Enqueued);
        Assert.Equal(1, resp.Failed);
    }

    [Fact]
    public async Task CancelReturnsCounts()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/batches/b_1/cancel").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithBody("{\"cancelled\":3,\"skipped\":1,\"errors\":[]}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var resp = await client.Batches.CancelAsync("acc_1", "b_1");
        Assert.Equal(3, resp.Cancelled);
        Assert.Equal(1, resp.Skipped);
    }
}
