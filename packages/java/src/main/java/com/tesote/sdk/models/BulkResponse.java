package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public record BulkResponse(@JsonProperty("bulk_results") List<BulkResult> bulkResults) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record BulkResult(
            @JsonProperty("account_id") String accountId,
            @JsonProperty("transactions") List<Transaction> transactions,
            @JsonProperty("pagination") CursorPagination pagination
    ) {}
}
