package com.tesote.sdk.models;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record Status(
        @JsonProperty("status") String status,
        @JsonProperty("authenticated") Boolean authenticated
) {}
