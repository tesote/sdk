using System.Linq;
using System.Threading.Tasks;
using Tesote.Sdk.Errors;
using Tesote.Sdk.Models;
using Tesote.Sdk.V2;
using WireMock.RequestBuilders;
using WireMock.ResponseBuilders;
using WireMock.Server;
using Xunit;

namespace Tesote.Sdk.Tests;

public sealed class TransactionOrdersTests : System.IDisposable
{
    private readonly WireMockServer _server;

    public TransactionOrdersTests()
    {
        _server = WireMockServer.Start();
    }

    public void Dispose()
    {
        _server.Stop();
        _server.Dispose();
    }

    private const string OrderJson =
        "{\"id\":\"to_1\",\"status\":\"draft\",\"amount\":100.50,\"currency\":\"VES\"," +
        "\"description\":\"pay\",\"reference\":null,\"external_reference\":null," +
        "\"idempotency_key\":null,\"batch_id\":null,\"scheduled_for\":null," +
        "\"approved_at\":null,\"submitted_at\":null,\"completed_at\":null," +
        "\"failed_at\":null,\"cancelled_at\":null," +
        "\"source_account\":{\"id\":\"acc_1\",\"name\":\"Checking\",\"payment_method_id\":\"pm_s\"}," +
        "\"destination\":{\"payment_method_id\":\"pm_d\",\"counterparty_id\":\"cp_1\",\"counterparty_name\":\"Dest\"}," +
        "\"fee\":null,\"execution_strategy\":null,\"tesote_transaction\":null,\"latest_attempt\":null," +
        "\"metadata\":null,\"created_at\":\"2026-04-01T00:00:00Z\",\"updated_at\":\"2026-04-01T00:00:00Z\"}";

    [Fact]
    public async Task ListReturnsOffsetPagination()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transaction_orders").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithBody("{\"items\":[" + OrderJson + "],\"has_more\":false,\"limit\":50,\"offset\":0}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var result = await client.TransactionOrders.ListAsync("acc_1", limit: 50);
        Assert.Single(result.Items);
        Assert.Equal("to_1", result.Items[0].Id);
        Assert.Equal("draft", result.Items[0].Status);
    }

    [Fact]
    public async Task GetMaps404ToTransactionOrderNotFound()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transaction_orders/missing").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(404)
                .WithBody("{\"error\":\"missing\",\"error_code\":\"TRANSACTION_ORDER_NOT_FOUND\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<TransactionOrderNotFoundException>(
            () => client.TransactionOrders.GetAsync("acc_1", "missing"));
    }

    [Fact]
    public async Task CreateSendsTransactionOrderEnvelopeAndIdempotencyKey()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transaction_orders").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(201).WithBody(OrderJson));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var req = new TransactionOrdersClient.CreateRequest
        {
            DestinationPaymentMethodId = "pm_d",
            Amount = "100.50",
            Currency = "VES",
            Description = "pay",
        };
        var order = await client.TransactionOrders.CreateAsync("acc_1", req, idempotencyKey: "tk-1");
        Assert.Equal("to_1", order.Id);

        var entry = Assert.Single(_server.LogEntries);
        Assert.Contains("transaction_order", entry.RequestMessage.Body);
        Assert.Equal("tk-1", entry.RequestMessage.Headers!["Idempotency-Key"].First());
    }

    [Fact]
    public async Task CreateMaps400ToValidationException()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transaction_orders").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(400)
                .WithBody("{\"error\":\"bad amount\",\"error_code\":\"VALIDATION_ERROR\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<ValidationException>(() =>
            client.TransactionOrders.CreateAsync("acc_1",
                new TransactionOrdersClient.CreateRequest()));
    }

    [Fact]
    public async Task SubmitMaps409ToInvalidOrderState()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transaction_orders/to_1/submit").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(409)
                .WithBody("{\"error\":\"bad state\",\"error_code\":\"INVALID_ORDER_STATE\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<InvalidOrderStateException>(
            () => client.TransactionOrders.SubmitAsync("acc_1", "to_1", token: "otp-123"));
    }

    [Fact]
    public async Task CancelReturnsCancelledOrder()
    {
        var cancelled = OrderJson.Replace("\"status\":\"draft\"", "\"status\":\"cancelled\"")
            .Replace("\"cancelled_at\":null", "\"cancelled_at\":\"2026-04-01T00:00:00Z\"");
        _server
            .Given(Request.Create().WithPath("/api/v2/accounts/acc_1/transaction_orders/to_1/cancel").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(200).WithBody(cancelled));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var order = await client.TransactionOrders.CancelAsync("acc_1", "to_1");
        Assert.Equal("cancelled", order.Status);
        Assert.NotNull(order.CancelledAt);
    }
}
