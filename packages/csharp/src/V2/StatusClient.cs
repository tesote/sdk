using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V2;

/// <summary>v2 status + whoami endpoints.</summary>
public sealed class StatusClient
{
    private readonly Transport _transport;

    internal StatusClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>Public status check.</summary>
    public async Task<StatusResponse> GetAsync(CancellationToken ct = default)
    {
        var node = await _transport.RequestAsync(RequestOptions.Get("/v2/status"), ct).ConfigureAwait(false);
        return Json.Deserialize<StatusResponse>(node);
    }

    /// <summary>Identifies the bearer-token holder.</summary>
    public async Task<WhoamiResponse> WhoamiAsync(CancellationToken ct = default)
    {
        var node = await _transport.RequestAsync(RequestOptions.Get("/v2/whoami"), ct).ConfigureAwait(false);
        return Json.Deserialize<WhoamiResponse>(node);
    }
}
