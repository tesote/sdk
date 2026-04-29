package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record Whoami(@JsonProperty("client") Client client) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    public record Client(
            @JsonProperty("id") String id,
            @JsonProperty("name") String name,
            @JsonProperty("type") String type
    ) {}
}
