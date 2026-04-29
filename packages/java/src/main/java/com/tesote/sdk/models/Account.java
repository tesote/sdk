package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * Account record (v1 and v2 share the same shape).
 *
 * <p>Balance fields inside {@link AccountData} only appear when the
 * workspace has {@code display_balances_in_api} enabled.
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record Account(
        @JsonProperty("id") String id,
        @JsonProperty("name") String name,
        @JsonProperty("data") AccountData data,
        @JsonProperty("bank") Bank bank,
        @JsonProperty("legal_entity") LegalEntity legalEntity,
        @JsonProperty("tesote_created_at") String tesoteCreatedAt,
        @JsonProperty("tesote_updated_at") String tesoteUpdatedAt
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record AccountData(
            @JsonProperty("masked_account_number") String maskedAccountNumber,
            @JsonProperty("currency") String currency,
            @JsonProperty("transactions_data_current_as_of") String transactionsDataCurrentAsOf,
            @JsonProperty("balance_data_current_as_of") String balanceDataCurrentAsOf,
            @JsonProperty("custom_user_provided_identifier") String customUserProvidedIdentifier,
            @JsonProperty("balance_cents") String balanceCents,
            @JsonProperty("available_balance_cents") String availableBalanceCents
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Bank(@JsonProperty("name") String name) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record LegalEntity(
            @JsonProperty("id") String id,
            @JsonProperty("legal_name") String legalName
    ) {}
}
