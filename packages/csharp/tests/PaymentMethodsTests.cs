using System.Linq;
using System.Threading.Tasks;
using Tesote.Sdk.Errors;
using Tesote.Sdk.V2;
using WireMock.RequestBuilders;
using WireMock.ResponseBuilders;
using WireMock.Server;
using Xunit;

namespace Tesote.Sdk.Tests;

public sealed class PaymentMethodsTests : System.IDisposable
{
    private readonly WireMockServer _server;

    public PaymentMethodsTests()
    {
        _server = WireMockServer.Start();
    }

    public void Dispose()
    {
        _server.Stop();
        _server.Dispose();
    }

    private const string PmJson =
        "{\"id\":\"pm_1\",\"method_type\":\"bank_account\",\"currency\":\"VES\"," +
        "\"label\":null,\"details\":{\"bank_code\":\"0001\",\"account_number\":\"1234\"," +
        "\"holder_name\":\"Acme\",\"identification_type\":null,\"identification_number\":null}," +
        "\"verified\":true,\"verified_at\":\"2026-04-01T00:00:00Z\"," +
        "\"last_used_at\":null,\"counterparty\":{\"id\":\"cp_1\",\"name\":\"Dest\"}," +
        "\"tesote_account\":null,\"created_at\":\"2026-04-01T00:00:00Z\"," +
        "\"updated_at\":\"2026-04-01T00:00:00Z\"}";

    [Fact]
    public async Task ListSendsFiltersAndDeserializes()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/payment_methods").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(200)
                .WithBody("{\"items\":[" + PmJson + "],\"has_more\":false,\"limit\":50,\"offset\":0}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var result = await client.PaymentMethods.ListAsync(
            limit: 50, methodType: "bank_account", verified: true);
        Assert.Single(result.Items);
        Assert.Equal("pm_1", result.Items[0].Id);

        var entry = Assert.Single(_server.LogEntries);
        Assert.Contains("method_type=bank_account", entry.RequestMessage.Url);
        Assert.Contains("verified=true", entry.RequestMessage.Url);
    }

    [Fact]
    public async Task GetMaps404ToPaymentMethodNotFound()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/payment_methods/missing").UsingGet())
            .RespondWith(Response.Create().WithStatusCode(404)
                .WithBody("{\"error\":\"missing\",\"error_code\":\"PAYMENT_METHOD_NOT_FOUND\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<PaymentMethodNotFoundException>(
            () => client.PaymentMethods.GetAsync("missing"));
    }

    [Fact]
    public async Task CreateSendsPaymentMethodEnvelopeAndDeserializes()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/payment_methods").UsingPost())
            .RespondWith(Response.Create().WithStatusCode(201).WithBody(PmJson));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var req = new PaymentMethodsClient.WriteRequest
        {
            MethodType = "bank_account",
            Currency = "VES",
            Counterparty = new PaymentMethodsClient.CounterpartyInput { Name = "Dest" },
            Details = new PaymentMethodsClient.DetailsInput
            {
                BankCode = "0001",
                AccountNumber = "1234",
                HolderName = "Acme",
            },
        };
        var pm = await client.PaymentMethods.CreateAsync(req);
        Assert.Equal("pm_1", pm.Id);

        var entry = Assert.Single(_server.LogEntries);
        Assert.Contains("payment_method", entry.RequestMessage.Body);
        Assert.Contains("bank_account", entry.RequestMessage.Body);
    }

    [Fact]
    public async Task UpdateUsesPatch()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/payment_methods/pm_1").UsingPatch())
            .RespondWith(Response.Create().WithStatusCode(200).WithBody(PmJson));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var pm = await client.PaymentMethods.UpdateAsync("pm_1",
            new PaymentMethodsClient.WriteRequest { Label = "renamed" });
        Assert.Equal("pm_1", pm.Id);
        Assert.Equal("PATCH", _server.LogEntries.Single().RequestMessage.Method);
    }

    [Fact]
    public async Task DeleteIssuesDeleteAndAcceptsNoContent()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/payment_methods/pm_1").UsingDelete())
            .RespondWith(Response.Create().WithStatusCode(204));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await client.PaymentMethods.DeleteAsync("pm_1");
        Assert.Equal("DELETE", _server.LogEntries.Single().RequestMessage.Method);
    }

    [Fact]
    public async Task DeleteMaps409ToValidationException()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/payment_methods/pm_1").UsingDelete())
            .RespondWith(Response.Create().WithStatusCode(409)
                .WithBody("{\"error\":\"in use\",\"error_code\":\"VALIDATION_ERROR\"}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        await Assert.ThrowsAsync<ValidationException>(() => client.PaymentMethods.DeleteAsync("pm_1"));
    }
}
