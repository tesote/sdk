package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public record TransactionsPage(
        @JsonProperty("total") Integer total,
        @JsonProperty("transactions") List<Transaction> transactions,
        @JsonProperty("pagination") CursorPagination pagination
) {}
