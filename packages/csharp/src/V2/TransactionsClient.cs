using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V2;

/// <summary>
/// v2 Transactions resource. Per-account list/export/sync, plus the cross-account
/// lookup, bulk fetch, and search endpoints.
/// </summary>
public sealed class TransactionsClient
{
    private readonly Transport _transport;

    internal TransactionsClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>Filter parameters accepted by the v2 transactions list/search endpoints.</summary>
    public sealed class ListFilters
    {
        /// <summary>ISO8601 date — earliest tesote-imported date to include.</summary>
        public string? StartDate { get; set; }
        /// <summary>ISO8601 date — latest tesote-imported date to include.</summary>
        public string? EndDate { get; set; }
        /// <summary>Server-side scope name (per docs).</summary>
        public string? Scope { get; set; }
        /// <summary>1-indexed page number.</summary>
        public int? Page { get; set; }
        /// <summary>Page size, default 50, max 100.</summary>
        public int? PerPage { get; set; }
        /// <summary>Cursor — fetch transactions after this id.</summary>
        public string? TransactionsAfterId { get; set; }
        /// <summary>Cursor — fetch transactions before this id.</summary>
        public string? TransactionsBeforeId { get; set; }
        /// <summary>ISO8601 — transaction-date floor.</summary>
        public string? TransactionDateAfter { get; set; }
        /// <summary>ISO8601 — transaction-date ceiling.</summary>
        public string? TransactionDateBefore { get; set; }
        /// <summary>ISO8601 — server-side <c>created_at</c> floor.</summary>
        public string? CreatedAfter { get; set; }
        /// <summary>ISO8601 — server-side <c>updated_at</c> floor.</summary>
        public string? UpdatedAfter { get; set; }
        /// <summary>Numeric — minimum amount (signed).</summary>
        public decimal? AmountMin { get; set; }
        /// <summary>Numeric — maximum amount.</summary>
        public decimal? AmountMax { get; set; }
        /// <summary>Numeric — exact amount match.</summary>
        public decimal? Amount { get; set; }
        /// <summary>Status filter (e.g. <c>posted</c>, <c>pending</c>).</summary>
        public string? Status { get; set; }
        /// <summary>Filter by category id.</summary>
        public string? CategoryId { get; set; }
        /// <summary>Filter by counterparty id.</summary>
        public string? CounterpartyId { get; set; }
        /// <summary>Search string (description / counterparty name).</summary>
        public string? Q { get; set; }
        /// <summary>Type filter (per docs).</summary>
        public string? Type { get; set; }
        /// <summary>Reference-code filter.</summary>
        public string? ReferenceCode { get; set; }

        internal IReadOnlyDictionary<string, string>? Build()
        {
            return new QueryBuilder()
                .Add("start_date", StartDate)
                .Add("end_date", EndDate)
                .Add("scope", Scope)
                .Add("page", Page)
                .Add("per_page", PerPage)
                .Add("transactions_after_id", TransactionsAfterId)
                .Add("transactions_before_id", TransactionsBeforeId)
                .Add("transaction_date_after", TransactionDateAfter)
                .Add("transaction_date_before", TransactionDateBefore)
                .Add("created_after", CreatedAfter)
                .Add("updated_after", UpdatedAfter)
                .Add("amount_min", AmountMin)
                .Add("amount_max", AmountMax)
                .Add("amount", Amount)
                .Add("status", Status)
                .Add("category_id", CategoryId)
                .Add("counterparty_id", CounterpartyId)
                .Add("q", Q)
                .Add("type", Type)
                .Add("reference_code", ReferenceCode)
                .BuildOrNull();
        }
    }

    /// <summary>List transactions for an account with the full v2 filter surface.</summary>
    public async Task<TransactionListResponse> ListForAccountAsync(
        string accountId,
        ListFilters? filters = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId) + "/transactions";
        var opts = RequestOptions.Get(path);
        opts.Query = (filters ?? new ListFilters()).Build();
        opts.CacheTtl = TimeSpan.FromMinutes(1);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<TransactionListResponse>(node);
    }

    /// <summary>Export transactions as CSV or JSON file. Returns raw bytes plus content-type.</summary>
    public Task<RawResponse> ExportAsync(
        string accountId,
        string format = "csv",
        ListFilters? filters = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentException.ThrowIfNullOrEmpty(format);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId) + "/transactions/export";
        var query = (filters ?? new ListFilters()).Build();
        var qb = new QueryBuilder().Add("format", format);
        if (query is not null)
        {
            foreach (var kv in query)
            {
                qb.Add(kv.Key, kv.Value);
            }
        }
        var opts = RequestOptions.Get(path);
        opts.Query = qb.BuildOrNull();
        return _transport.RequestRawAsync(opts, ct);
    }

    /// <summary>Sync request body for POST /v2/.../transactions/sync.</summary>
    public sealed class SyncRequest
    {
        /// <summary>Number of transactions to fetch (1..1000).</summary>
        public int? Count { get; set; }
        /// <summary>Opaque cursor or the literal <c>"now"</c>.</summary>
        public string? Cursor { get; set; }
        /// <summary>Optional flags accepted by the API.</summary>
        public SyncOptions? Options { get; set; }
    }

    /// <summary>Optional flag bag for <see cref="SyncRequest"/>.</summary>
    public sealed class SyncOptions
    {
        /// <summary>Include per-transaction running balances when the workspace allows it.</summary>
        public bool? IncludeRunningBalance { get; set; }
    }

    /// <summary>Cursor-based incremental sync for a single account.</summary>
    public async Task<SyncResult> SyncAsync(
        string accountId,
        SyncRequest request,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentNullException.ThrowIfNull(request);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId) + "/transactions/sync";
        var body = BuildSyncBody(request);
        var opts = Requests.Json("POST", path, body, idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<SyncResult>(node);
    }

    /// <summary>Legacy non-nested sync route kept for backwards compatibility.</summary>
    public async Task<SyncResult> SyncLegacyAsync(
        SyncRequest request,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentNullException.ThrowIfNull(request);
        var body = BuildSyncBody(request);
        var opts = Requests.Json("POST", "/v2/transactions/sync", body, idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<SyncResult>(node);
    }

    /// <summary>Fetch a single transaction by id (v1-shape payload).</summary>
    public async Task<Transaction> GetAsync(string transactionId, CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(transactionId);
        var opts = RequestOptions.Get("/v2/transactions/" + Uri.EscapeDataString(transactionId));
        opts.CacheTtl = TimeSpan.FromMinutes(5);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<Transaction>(node);
    }

    /// <summary>Bulk-fetch transactions for up to 100 accounts in one call.</summary>
    public async Task<BulkResponse> BulkAsync(
        IReadOnlyList<string> accountIds,
        int? page = null,
        int? perPage = null,
        int? limit = null,
        int? offset = null,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentNullException.ThrowIfNull(accountIds);
        var body = new Dictionary<string, object?>
        {
            ["account_ids"] = accountIds,
            ["page"] = page,
            ["per_page"] = perPage,
            ["limit"] = limit,
            ["offset"] = offset,
        };
        var opts = Requests.Json("POST", "/v2/transactions/bulk", body, idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<BulkResponse>(node);
    }

    /// <summary>Search transactions across accounts using a substring match plus the standard filter set.</summary>
    public async Task<SearchResult> SearchAsync(
        string query,
        string? accountId = null,
        int? limit = null,
        int? offset = null,
        ListFilters? filters = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(query);
        var qb = new QueryBuilder()
            .Add("q", query)
            .Add("account_id", accountId)
            .Add("limit", limit)
            .Add("offset", offset);
        var extra = (filters ?? new ListFilters()).Build();
        if (extra is not null)
        {
            foreach (var kv in extra)
            {
                if (!kv.Key.Equals("q", StringComparison.Ordinal))
                {
                    qb.Add(kv.Key, kv.Value);
                }
            }
        }
        var opts = RequestOptions.Get("/v2/transactions/search");
        opts.Query = qb.BuildOrNull();
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<SearchResult>(node);
    }

    private static Dictionary<string, object?> BuildSyncBody(SyncRequest request)
    {
        var body = new Dictionary<string, object?>();
        if (request.Count is not null)
        {
            body["count"] = request.Count.Value;
        }
        if (request.Cursor is not null)
        {
            body["cursor"] = request.Cursor;
        }
        if (request.Options is not null)
        {
            var inner = new Dictionary<string, object?>();
            if (request.Options.IncludeRunningBalance is not null)
            {
                inner["include_running_balance"] = request.Options.IncludeRunningBalance.Value;
            }
            body["options"] = inner;
        }
        return body;
    }
}
