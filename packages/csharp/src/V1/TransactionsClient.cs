using System;
using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V1;

/// <summary>v1 Transactions resource: per-account listing and single-id lookup.</summary>
public sealed class TransactionsClient
{
    private const string BasePath = "/v1/transactions";
    private const string AccountsBasePath = "/v1/accounts";

    private readonly Transport _transport;

    internal TransactionsClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>List transactions for a single account, with optional date and cursor filters.</summary>
    public async Task<TransactionListResponse> ListForAccountAsync(
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

        var opts = RequestOptions.Get(AccountsBasePath + "/" + Uri.EscapeDataString(accountId) + "/transactions");
        opts.Query = query;
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<TransactionListResponse>(node);
    }

    /// <summary>Fetch a single transaction by id.</summary>
    public async Task<Transaction> GetAsync(string transactionId, CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(transactionId);
        var opts = RequestOptions.Get(BasePath + "/" + Uri.EscapeDataString(transactionId));
        opts.CacheTtl = TimeSpan.FromMinutes(5);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<Transaction>(node);
    }
}
