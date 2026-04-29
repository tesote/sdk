package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Map;

@JsonIgnoreProperties(ignoreUnknown = true)
public record PaymentMethod(
        @JsonProperty("id") String id,
        @JsonProperty("method_type") String methodType,
        @JsonProperty("currency") String currency,
        @JsonProperty("label") String label,
        @JsonProperty("details") Map<String, Object> details,
        @JsonProperty("verified") Boolean verified,
        @JsonProperty("verified_at") String verifiedAt,
        @JsonProperty("last_used_at") String lastUsedAt,
        @JsonProperty("counterparty") CounterpartyRef counterparty,
        @JsonProperty("tesote_account") TesoteAccountRef tesoteAccount,
        @JsonProperty("created_at") String createdAt,
        @JsonProperty("updated_at") String updatedAt
) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record CounterpartyRef(
            @JsonProperty("id") String id,
            @JsonProperty("name") String name
    ) {}

    @JsonIgnoreProperties(ignoreUnknown = true)
    public record TesoteAccountRef(
            @JsonProperty("id") String id,
            @JsonProperty("name") String name
    ) {}
}
