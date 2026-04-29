package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * Pagination metadata for page-based endpoints (accounts list).
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public record PagePagination(
        @JsonProperty("current_page") Integer currentPage,
        @JsonProperty("per_page") Integer perPage,
        @JsonProperty("total_pages") Integer totalPages,
        @JsonProperty("total_count") Integer totalCount
) {}
