package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

/**
 * Transaction record. Used for /v1 reads and /v2/transactions/{id}.
 * The flat {@link SyncTransaction} variant is only used in v2 sync responses.
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record Transaction(
        @JsonProperty("id") String id,
        @JsonProperty("status") String status,
        @JsonProperty("data") TransactionData data,
        @JsonProperty("tesote_imported_at") String tesoteImportedAt,
        @JsonProperty("tesote_updated_at") String tesoteUpdatedAt,
        @JsonProperty("transaction_categories") List<Category> transactionCategories,
        @JsonProperty("counterparty") Counterparty counterparty
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record TransactionData(
            @JsonProperty("amount_cents") Long amountCents,
            @JsonProperty("currency") String currency,
            @JsonProperty("description") String description,
            @JsonProperty("transaction_date") String transactionDate,
            @JsonProperty("created_at") String createdAt,
            @JsonProperty("created_at_date") String createdAtDate,
            @JsonProperty("note") String note,
            @JsonProperty("external_service_id") String externalServiceId,
            @JsonProperty("running_balance_cents") Long runningBalanceCents
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Category(
            @JsonProperty("name") String name,
            @JsonProperty("external_category_code") String externalCategoryCode,
            @JsonProperty("created_at") String createdAt,
            @JsonProperty("updated_at") String updatedAt
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Counterparty(@JsonProperty("name") String name) {}
}
