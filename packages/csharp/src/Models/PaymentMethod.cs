using System.Collections.Generic;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace Tesote.Sdk.Models;

/// <summary>Counterparty associated with a beneficiary <see cref="PaymentMethod"/>.</summary>
public sealed record PaymentMethodCounterparty(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("name")] string Name);

/// <summary>Tesote-side account associated with a source <see cref="PaymentMethod"/>.</summary>
public sealed record PaymentMethodTesoteAccount(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("name")] string Name);

/// <summary>
/// Type-specific payload for a <see cref="PaymentMethod"/>. Common fields are typed; the
/// <see cref="Extra"/> bag captures method-type-specific properties (e.g. crypto wallet address).
/// </summary>
/// <remarks>
/// Declared as a class (not a positional record) because <see cref="JsonExtensionDataAttribute"/>
/// cannot bind to a constructor parameter.
/// </remarks>
public sealed class PaymentMethodDetails
{
    /// <summary>Bank routing/clabe code, when applicable.</summary>
    [JsonPropertyName("bank_code")] public string? BankCode { get; init; }

    /// <summary>Beneficiary account number, when applicable.</summary>
    [JsonPropertyName("account_number")] public string? AccountNumber { get; init; }

    /// <summary>Account holder name, when applicable.</summary>
    [JsonPropertyName("holder_name")] public string? HolderName { get; init; }

    /// <summary>Identification document type, when applicable.</summary>
    [JsonPropertyName("identification_type")] public string? IdentificationType { get; init; }

    /// <summary>Identification document number, when applicable.</summary>
    [JsonPropertyName("identification_number")] public string? IdentificationNumber { get; init; }

    /// <summary>Open-ended bag for type-specific extras the SDK doesn't model.</summary>
    [JsonExtensionData] public IDictionary<string, JsonElement>? Extra { get; init; }
}

/// <summary>Stored beneficiary or source payment method.</summary>
public sealed record PaymentMethod(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("method_type")] string MethodType,
    [property: JsonPropertyName("currency")] string Currency,
    [property: JsonPropertyName("label")] string? Label,
    [property: JsonPropertyName("details")] PaymentMethodDetails Details,
    [property: JsonPropertyName("verified")] bool Verified,
    [property: JsonPropertyName("verified_at")] string? VerifiedAt,
    [property: JsonPropertyName("last_used_at")] string? LastUsedAt,
    [property: JsonPropertyName("counterparty")] PaymentMethodCounterparty? Counterparty,
    [property: JsonPropertyName("tesote_account")] PaymentMethodTesoteAccount? TesoteAccount,
    [property: JsonPropertyName("created_at")] string CreatedAt,
    [property: JsonPropertyName("updated_at")] string UpdatedAt);

/// <summary>Response envelope for GET /v2/payment_methods.</summary>
public sealed record PaymentMethodListResponse(
    [property: JsonPropertyName("items")] IReadOnlyList<PaymentMethod> Items,
    [property: JsonPropertyName("has_more")] bool HasMore,
    [property: JsonPropertyName("limit")] int Limit,
    [property: JsonPropertyName("offset")] int Offset);
