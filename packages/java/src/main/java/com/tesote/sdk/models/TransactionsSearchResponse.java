package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public record TransactionsSearchResponse(
        @JsonProperty("transactions") List<Transaction> transactions,
        @JsonProperty("total") Integer total
) {}
