package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

/**
 * Response body for {@code POST /v2/accounts/{id}/transactions/sync} and the
 * legacy {@code POST /v2/transactions/sync}.
 *
 * <p>Persist {@link #nextCursor()} between calls to resume.
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record SyncTransactionsResponse(
        @JsonProperty("added") List<SyncTransaction> added,
        @JsonProperty("modified") List<SyncTransaction> modified,
        @JsonProperty("removed") List<RemovedTransaction> removed,
        @JsonProperty("next_cursor") String nextCursor,
        @JsonProperty("has_more") Boolean hasMore
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record RemovedTransaction(
            @JsonProperty("transaction_id") String transactionId,
            @JsonProperty("account_id") String accountId
    ) {}
}
