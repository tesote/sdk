package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.math.BigDecimal;
import java.util.List;

/**
 * Plaid-compatible flat transaction shape returned by sync endpoints.
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record SyncTransaction(
        @JsonProperty("transaction_id") String transactionId,
        @JsonProperty("account_id") String accountId,
        @JsonProperty("amount") BigDecimal amount,
        @JsonProperty("iso_currency_code") String isoCurrencyCode,
        @JsonProperty("unofficial_currency_code") String unofficialCurrencyCode,
        @JsonProperty("date") String date,
        @JsonProperty("datetime") String datetime,
        @JsonProperty("name") String name,
        @JsonProperty("merchant_name") String merchantName,
        @JsonProperty("pending") Boolean pending,
        @JsonProperty("category") List<String> category,
        @JsonProperty("running_balance_cents") Long runningBalanceCents
) {}
