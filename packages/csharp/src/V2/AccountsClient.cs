using System;
using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V2;

/// <summary>
/// v2 Accounts resource. Same shape as v1 plus a sync trigger; transaction listing,
/// export, and sync live on dedicated client classes.
/// </summary>
public sealed class AccountsClient
{
    private const string BasePath = "/v2/accounts";

    private readonly Transport _transport;

    internal AccountsClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>List accounts with page-based pagination.</summary>
    public async Task<AccountListResponse> ListAsync(
        int? page = null,
        int? perPage = null,
        string? include = null,
        string? sort = null,
        CancellationToken ct = default)
    {
        var query = new QueryBuilder()
            .Add("page", page)
            .Add("per_page", perPage)
            .Add("include", include)
            .Add("sort", sort)
            .BuildOrNull();

        var opts = RequestOptions.Get(BasePath);
        opts.Query = query;
        opts.CacheTtl = TimeSpan.FromMinutes(1);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<AccountListResponse>(node);
    }

    /// <summary>Fetch a single account by id.</summary>
    public async Task<Account> GetAsync(string accountId, CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        var opts = RequestOptions.Get(BasePath + "/" + Uri.EscapeDataString(accountId));
        opts.CacheTtl = TimeSpan.FromMinutes(5);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<Account>(node);
    }

    /// <summary>Trigger a bank sync for an account. Returns the started sync session metadata.</summary>
    public async Task<SyncStartResponse> SyncAsync(
        string accountId,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        var path = BasePath + "/" + Uri.EscapeDataString(accountId) + "/sync";
        var opts = Requests.Json("POST", path, body: null, idempotencyKey: idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<SyncStartResponse>(node);
    }
}
