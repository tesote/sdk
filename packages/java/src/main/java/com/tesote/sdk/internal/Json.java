package com.tesote.sdk.internal;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

/**
 * Thin wrapper around a single shared {@link ObjectMapper}. ObjectMapper is
 * thread-safe once configured; sharing avoids re-allocating per request.
 */
public final class Json {
    public static final ObjectMapper MAPPER = new ObjectMapper();

    private Json() {}

    public static JsonNode parse(byte[] bytes) {
        if (bytes == null || bytes.length == 0) return MAPPER.nullNode();
        try {
            return MAPPER.readTree(bytes);
        } catch (JsonProcessingException e) {
            return MAPPER.nullNode();
        } catch (java.io.IOException e) {
            return MAPPER.nullNode();
        }
    }

    public static String stringify(Object value) {
        try {
            return MAPPER.writeValueAsString(value);
        } catch (JsonProcessingException e) {
            throw new IllegalStateException("failed to serialize " + value.getClass(), e);
        }
    }
}
