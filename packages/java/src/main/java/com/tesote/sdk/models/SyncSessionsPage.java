package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public record SyncSessionsPage(
        @JsonProperty("sync_sessions") List<SyncSession> syncSessions,
        @JsonProperty("limit") Integer limit,
        @JsonProperty("offset") Integer offset,
        @JsonProperty("has_more") Boolean hasMore
) {}
