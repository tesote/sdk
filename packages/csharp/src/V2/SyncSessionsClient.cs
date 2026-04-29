using System;
using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V2;

/// <summary>v2 sync-sessions resource — list and lookup per account.</summary>
public sealed class SyncSessionsClient
{
    private readonly Transport _transport;

    internal SyncSessionsClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>List sync sessions for an account, ordered by creation desc.</summary>
    public async Task<SyncSessionListResponse> ListAsync(
        string accountId,
        int? limit = null,
        int? offset = null,
        string? status = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId) + "/sync_sessions";
        var query = new QueryBuilder()
            .Add("limit", limit)
            .Add("offset", offset)
            .Add("status", status)
            .BuildOrNull();
        var opts = RequestOptions.Get(path);
        opts.Query = query;
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<SyncSessionListResponse>(node);
    }

    /// <summary>Fetch a single sync session by id.</summary>
    public async Task<SyncSession> GetAsync(string accountId, string sessionId, CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentException.ThrowIfNullOrEmpty(sessionId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId)
            + "/sync_sessions/" + Uri.EscapeDataString(sessionId);
        var opts = RequestOptions.Get(path);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<SyncSession>(node);
    }
}
