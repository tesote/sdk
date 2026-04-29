package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;
import java.util.Map;

@JsonIgnoreProperties(ignoreUnknown = true)
public record BatchSummary(
        @JsonProperty("batch_id") String batchId,
        @JsonProperty("total_orders") Integer totalOrders,
        @JsonProperty("total_amount_cents") Long totalAmountCents,
        @JsonProperty("amount_currency") String amountCurrency,
        @JsonProperty("statuses") Map<String, Integer> statuses,
        @JsonProperty("batch_status") String batchStatus,
        @JsonProperty("created_at") String createdAt,
        @JsonProperty("orders") List<TransactionOrder> orders
) {}
