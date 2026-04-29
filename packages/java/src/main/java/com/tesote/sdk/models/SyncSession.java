package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record SyncSession(
        @JsonProperty("id") String id,
        @JsonProperty("status") String status,
        @JsonProperty("started_at") String startedAt,
        @JsonProperty("completed_at") String completedAt,
        @JsonProperty("transactions_synced") Integer transactionsSynced,
        @JsonProperty("accounts_count") Integer accountsCount,
        @JsonProperty("error") SyncError error,
        @JsonProperty("performance") Performance performance
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record SyncError(
            @JsonProperty("type") String type,
            @JsonProperty("message") String message
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Performance(
            @JsonProperty("total_duration") Double totalDuration,
            @JsonProperty("complexity_score") Double complexityScore,
            @JsonProperty("sync_speed_score") Double syncSpeedScore
    ) {}
}
