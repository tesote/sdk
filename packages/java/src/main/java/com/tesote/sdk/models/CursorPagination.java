package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * Pagination metadata for cursor-based endpoints (transactions list).
 *
 * <p>{@code afterId} is the last item in the page; {@code beforeId} is the
 * first. Pass {@code afterId} as {@code transactions_after_id} on the next
 * call to fetch the following page.
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record CursorPagination(
        @JsonProperty("has_more") Boolean hasMore,
        @JsonProperty("per_page") Integer perPage,
        @JsonProperty("after_id") String afterId,
        @JsonProperty("before_id") String beforeId
) {}
