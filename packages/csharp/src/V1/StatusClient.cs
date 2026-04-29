using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V1;

/// <summary>v1 status + whoami endpoints (cross-version: <c>/status</c> and <c>/whoami</c>).</summary>
public sealed class StatusClient
{
    private readonly Transport _transport;

    internal StatusClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>Public status check. Auth not required server-side, but the SDK sends the bearer anyway.</summary>
    public async Task<StatusResponse> GetAsync(CancellationToken ct = default)
    {
        var node = await _transport.RequestAsync(RequestOptions.Get("/status"), ct).ConfigureAwait(false);
        return Json.Deserialize<StatusResponse>(node);
    }

    /// <summary>Identifies the bearer-token holder (workspace or user).</summary>
    public async Task<WhoamiResponse> WhoamiAsync(CancellationToken ct = default)
    {
        var node = await _transport.RequestAsync(RequestOptions.Get("/whoami"), ct).ConfigureAwait(false);
        return Json.Deserialize<WhoamiResponse>(node);
    }
}
