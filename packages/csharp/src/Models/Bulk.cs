using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Tesote.Sdk.Models;

/// <summary>Per-account result of a POST /v2/transactions/bulk request.</summary>
public sealed record BulkResult(
    [property: JsonPropertyName("account_id")] string AccountId,
    [property: JsonPropertyName("transactions")] IReadOnlyList<Transaction> Transactions,
    [property: JsonPropertyName("pagination")] CursorPagination Pagination);

/// <summary>Response envelope for POST /v2/transactions/bulk.</summary>
public sealed record BulkResponse(
    [property: JsonPropertyName("bulk_results")] IReadOnlyList<BulkResult> BulkResults);

/// <summary>Response envelope for GET /v2/transactions/search.</summary>
public sealed record SearchResult(
    [property: JsonPropertyName("transactions")] IReadOnlyList<Transaction> Transactions,
    [property: JsonPropertyName("total")] int Total);
