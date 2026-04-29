package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;
import java.util.Map;

/**
 * Response shape for batch approve / submit / cancel.
 * Different actions populate different counts; cancel additionally fills
 * {@code skipped} and may include per-order {@code errors}.
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record BatchActionResponse(
        @JsonProperty("approved") Integer approved,
        @JsonProperty("enqueued") Integer enqueued,
        @JsonProperty("cancelled") Integer cancelled,
        @JsonProperty("failed") Integer failed,
        @JsonProperty("skipped") Integer skipped,
        @JsonProperty("errors") List<Map<String, Object>> errors
) {}
