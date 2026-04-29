package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * Response body for {@code POST /v2/accounts/{id}/sync}.
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record AccountSyncResponse(
        @JsonProperty("message") String message,
        @JsonProperty("sync_session_id") String syncSessionId,
        @JsonProperty("status") String status,
        @JsonProperty("started_at") String startedAt
) {}
