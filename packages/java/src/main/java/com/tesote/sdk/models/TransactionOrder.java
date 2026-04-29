package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.math.BigDecimal;
import java.util.Map;

/**
 * Transaction order — the unit of payment intent in v2 batches and
 * transaction-order endpoints.
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record TransactionOrder(
        @JsonProperty("id") String id,
        @JsonProperty("status") String status,
        @JsonProperty("amount") BigDecimal amount,
        @JsonProperty("currency") String currency,
        @JsonProperty("description") String description,
        @JsonProperty("reference") String reference,
        @JsonProperty("external_reference") String externalReference,
        @JsonProperty("idempotency_key") String idempotencyKey,
        @JsonProperty("batch_id") String batchId,
        @JsonProperty("scheduled_for") String scheduledFor,
        @JsonProperty("approved_at") String approvedAt,
        @JsonProperty("submitted_at") String submittedAt,
        @JsonProperty("completed_at") String completedAt,
        @JsonProperty("failed_at") String failedAt,
        @JsonProperty("cancelled_at") String cancelledAt,
        @JsonProperty("source_account") SourceAccount sourceAccount,
        @JsonProperty("destination") Destination destination,
        @JsonProperty("fee") Fee fee,
        @JsonProperty("execution_strategy") String executionStrategy,
        @JsonProperty("tesote_transaction") TesoteTransaction tesoteTransaction,
        @JsonProperty("latest_attempt") LatestAttempt latestAttempt,
        @JsonProperty("metadata") Map<String, Object> metadata,
        @JsonProperty("created_at") String createdAt,
        @JsonProperty("updated_at") String updatedAt
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record SourceAccount(
            @JsonProperty("id") String id,
            @JsonProperty("name") String name,
            @JsonProperty("payment_method_id") String paymentMethodId
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Destination(
            @JsonProperty("payment_method_id") String paymentMethodId,
            @JsonProperty("counterparty_id") String counterpartyId,
            @JsonProperty("counterparty_name") String counterpartyName
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Fee(
            @JsonProperty("amount") BigDecimal amount,
            @JsonProperty("currency") String currency
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record TesoteTransaction(
            @JsonProperty("id") String id,
            @JsonProperty("status") String status
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record LatestAttempt(
            @JsonProperty("id") String id,
            @JsonProperty("status") String status,
            @JsonProperty("attempt_number") Integer attemptNumber,
            @JsonProperty("external_reference") String externalReference,
            @JsonProperty("submitted_at") String submittedAt,
            @JsonProperty("completed_at") String completedAt,
            @JsonProperty("error_code") String errorCode,
            @JsonProperty("error_message") String errorMessage
    ) {}
}
