using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Tesote.Sdk.Models;

/// <summary>Counterparty metadata embedded in a <see cref="Transaction"/>.</summary>
public sealed record TransactionCounterparty(
    [property: JsonPropertyName("name")] string Name);

/// <summary>Category attached to a <see cref="Transaction"/>.</summary>
public sealed record TransactionCategory(
    [property: JsonPropertyName("name")] string Name,
    [property: JsonPropertyName("external_category_code")] string? ExternalCategoryCode,
    [property: JsonPropertyName("created_at")] string CreatedAt,
    [property: JsonPropertyName("updated_at")] string UpdatedAt);

/// <summary>Wire-side transaction-data envelope.</summary>
public sealed record TransactionData(
    [property: JsonPropertyName("amount_cents")] long AmountCents,
    [property: JsonPropertyName("currency")] string Currency,
    [property: JsonPropertyName("description")] string Description,
    [property: JsonPropertyName("transaction_date")] string TransactionDate,
    [property: JsonPropertyName("created_at")] string? CreatedAt,
    [property: JsonPropertyName("created_at_date")] string? CreatedAtDate,
    [property: JsonPropertyName("note")] string? Note,
    [property: JsonPropertyName("external_service_id")] string? ExternalServiceId,
    [property: JsonPropertyName("running_balance_cents")] long? RunningBalanceCents);

/// <summary>v1 transaction (also returned from /v2/transactions/{id}).</summary>
public sealed record Transaction(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("status")] string Status,
    [property: JsonPropertyName("data")] TransactionData Data,
    [property: JsonPropertyName("tesote_imported_at")] string TesoteImportedAt,
    [property: JsonPropertyName("tesote_updated_at")] string TesoteUpdatedAt,
    [property: JsonPropertyName("transaction_categories")] IReadOnlyList<TransactionCategory> TransactionCategories,
    [property: JsonPropertyName("counterparty")] TransactionCounterparty? Counterparty);

/// <summary>Cursor-based pagination envelope used by transaction listings.</summary>
public sealed record CursorPagination(
    [property: JsonPropertyName("has_more")] bool HasMore,
    [property: JsonPropertyName("per_page")] int PerPage,
    [property: JsonPropertyName("after_id")] string? AfterId,
    [property: JsonPropertyName("before_id")] string? BeforeId);

/// <summary>Response envelope for GET /v1|v2/accounts/{id}/transactions.</summary>
public sealed record TransactionListResponse(
    [property: JsonPropertyName("total")] int Total,
    [property: JsonPropertyName("transactions")] IReadOnlyList<Transaction> Transactions,
    [property: JsonPropertyName("pagination")] CursorPagination Pagination);
