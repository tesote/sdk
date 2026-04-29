package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;
import java.util.Map;

@JsonIgnoreProperties(ignoreUnknown = true)
public record BatchCreateResponse(
        @JsonProperty("batch_id") String batchId,
        @JsonProperty("orders") List<TransactionOrder> orders,
        @JsonProperty("errors") List<Map<String, Object>> errors
) {}
