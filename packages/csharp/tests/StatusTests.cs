using System.Threading.Tasks;
using WireMock.RequestBuilders;
using WireMock.ResponseBuilders;
using WireMock.Server;
using Xunit;

namespace Tesote.Sdk.Tests;

public sealed class StatusTests : System.IDisposable
{
    private readonly WireMockServer _server;

    public StatusTests()
    {
        _server = WireMockServer.Start();
    }

    public void Dispose()
    {
        _server.Stop();
        _server.Dispose();
    }

    [Fact]
    public async Task V1StatusReturnsTypedResponse()
    {
        _server
            .Given(Request.Create().WithPath("/api/status").UsingGet())
            .RespondWith(Response.Create()
                .WithStatusCode(200)
                .WithHeader("Content-Type", "application/json")
                .WithBody("{\"status\":\"ok\",\"authenticated\":false}"));

        using var client = TestHelpers.NewV1(_server.Url + "/api");
        var result = await client.Status.GetAsync();
        Assert.Equal("ok", result.Status);
        Assert.False(result.Authenticated);
    }

    [Fact]
    public async Task V2WhoamiReturnsClientStub()
    {
        _server
            .Given(Request.Create().WithPath("/api/v2/whoami").UsingGet())
            .RespondWith(Response.Create()
                .WithStatusCode(200)
                .WithHeader("Content-Type", "application/json")
                .WithBody("{\"client\":{\"id\":\"c_1\",\"name\":\"Acme\",\"type\":\"workspace\"}}"));

        using var client = TestHelpers.NewV2(_server.Url + "/api");
        var result = await client.Status.WhoamiAsync();
        Assert.Equal("c_1", result.Client.Id);
        Assert.Equal("workspace", result.Client.Type);
    }
}
