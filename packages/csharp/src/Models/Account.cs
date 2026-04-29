using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Tesote.Sdk.Models;

/// <summary>Bank metadata embedded in an <see cref="Account"/>.</summary>
public sealed record AccountBank(
    [property: JsonPropertyName("name")] string Name);

/// <summary>Legal-entity attribution for an <see cref="Account"/>.</summary>
public sealed record AccountLegalEntity(
    [property: JsonPropertyName("id")] string? Id,
    [property: JsonPropertyName("legal_name")] string? LegalName);

/// <summary>Wire-side account-data envelope. Snake_case is preserved.</summary>
public sealed record AccountData(
    [property: JsonPropertyName("masked_account_number")] string MaskedAccountNumber,
    [property: JsonPropertyName("currency")] string Currency,
    [property: JsonPropertyName("transactions_data_current_as_of")] string? TransactionsDataCurrentAsOf,
    [property: JsonPropertyName("balance_data_current_as_of")] string? BalanceDataCurrentAsOf,
    [property: JsonPropertyName("custom_user_provided_identifier")] string? CustomUserProvidedIdentifier,
    [property: JsonPropertyName("balance_cents")] string? BalanceCents,
    [property: JsonPropertyName("available_balance_cents")] string? AvailableBalanceCents);

/// <summary>Account model. Identical between v1 and v2.</summary>
public sealed record Account(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("name")] string Name,
    [property: JsonPropertyName("data")] AccountData Data,
    [property: JsonPropertyName("bank")] AccountBank Bank,
    [property: JsonPropertyName("legal_entity")] AccountLegalEntity LegalEntity,
    [property: JsonPropertyName("tesote_created_at")] string TesoteCreatedAt,
    [property: JsonPropertyName("tesote_updated_at")] string TesoteUpdatedAt);

/// <summary>Page-based pagination envelope (v1 / v2 accounts list).</summary>
public sealed record PageBasedPagination(
    [property: JsonPropertyName("current_page")] int CurrentPage,
    [property: JsonPropertyName("per_page")] int PerPage,
    [property: JsonPropertyName("total_pages")] int TotalPages,
    [property: JsonPropertyName("total_count")] int TotalCount);

/// <summary>Response envelope for GET /v1/accounts and /v2/accounts.</summary>
public sealed record AccountListResponse(
    [property: JsonPropertyName("total")] int Total,
    [property: JsonPropertyName("accounts")] IReadOnlyList<Account> Accounts,
    [property: JsonPropertyName("pagination")] PageBasedPagination Pagination);
