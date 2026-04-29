using System;
using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V1;

/// <summary>v1 Accounts resource: list, get, and per-account transaction listing.</summary>
public sealed class AccountsClient
{
    private const string BasePath = "/v1/accounts";

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

    /// <summary>List transactions for a single account, with optional date and cursor filters.</summary>
    public async Task<TransactionListResponse> ListTransactionsAsync(
        string accountId,
        string? startDate = null,
        string? endDate = null,
        string? scope = null,
        int? page = null,
        int? perPage = null,
        string? transactionsAfterId = null,
        string? transactionsBeforeId = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        var query = new QueryBuilder()
            .Add("start_date", startDate)
            .Add("end_date", endDate)
            .Add("scope", scope)
            .Add("page", page)
            .Add("per_page", perPage)
            .Add("transactions_after_id", transactionsAfterId)
            .Add("transactions_before_id", transactionsBeforeId)
            .BuildOrNull();

        var opts = RequestOptions.Get(BasePath + "/" + Uri.EscapeDataString(accountId) + "/transactions");
        opts.Query = query;
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<TransactionListResponse>(node);
    }
}
