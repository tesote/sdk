package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

@JsonIgnoreProperties(ignoreUnknown = true)
public record AccountsPage(
        @JsonProperty("total") Integer total,
        @JsonProperty("accounts") List<Account> accounts,
        @JsonProperty("pagination") PagePagination pagination
) {}
