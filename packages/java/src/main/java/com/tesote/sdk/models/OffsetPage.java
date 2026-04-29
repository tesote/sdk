package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

/**
 * Generic offset-paginated wrapper used by transaction-orders and
 * payment-methods list endpoints. {@code items} is the page payload.
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record OffsetPage<T>(
        @JsonProperty("items") List<T> items,
        @JsonProperty("limit") Integer limit,
        @JsonProperty("offset") Integer offset,
        @JsonProperty("has_more") Boolean hasMore
) {}
