using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Tesote.Sdk.Models;

/// <summary>Flattened Plaid-shaped transaction returned from POST /v2/.../transactions/sync.</summary>
public sealed record SyncTransaction(
    [property: JsonPropertyName("transaction_id")] string TransactionId,
    [property: JsonPropertyName("account_id")] string AccountId,
    [property: JsonPropertyName("amount")] decimal Amount,
    [property: JsonPropertyName("iso_currency_code")] string IsoCurrencyCode,
    [property: JsonPropertyName("unofficial_currency_code")] string? UnofficialCurrencyCode,
    [property: JsonPropertyName("date")] string Date,
    [property: JsonPropertyName("datetime")] string? Datetime,
    [property: JsonPropertyName("name")] string Name,
    [property: JsonPropertyName("merchant_name")] string? MerchantName,
    [property: JsonPropertyName("pending")] bool Pending,
    [property: JsonPropertyName("category")] IReadOnlyList<string> Category,
    [property: JsonPropertyName("running_balance_cents")] long? RunningBalanceCents);

/// <summary>Identifiers for a transaction removed since the last sync.</summary>
public sealed record SyncRemoved(
    [property: JsonPropertyName("transaction_id")] string TransactionId,
    [property: JsonPropertyName("account_id")] string AccountId);

/// <summary>Response envelope for POST /v2/.../transactions/sync.</summary>
public sealed record SyncResult(
    [property: JsonPropertyName("added")] IReadOnlyList<SyncTransaction> Added,
    [property: JsonPropertyName("modified")] IReadOnlyList<SyncTransaction> Modified,
    [property: JsonPropertyName("removed")] IReadOnlyList<SyncRemoved> Removed,
    [property: JsonPropertyName("next_cursor")] string? NextCursor,
    [property: JsonPropertyName("has_more")] bool HasMore);
