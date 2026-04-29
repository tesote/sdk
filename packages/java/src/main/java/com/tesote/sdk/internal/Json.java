package com.tesote.sdk.internal;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.core.type.TypeReference;

/**
 * Thin wrapper around a single shared {@link ObjectMapper}. ObjectMapper is
 * thread-safe once configured; sharing avoids re-allocating per request.
 *
 * <p>No JSR-310 module: the SDK keeps jackson-databind as its only runtime
 * dep, so date/time API fields are exposed as {@code String} on the model
 * records and callers parse with {@link java.time.OffsetDateTime#parse}.
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

    public static <T> T treeToValue(JsonNode node, Class<T> type) {
        if (node == null || node.isMissingNode() || node.isNull()) return null;
        try {
            return MAPPER.treeToValue(node, type);
        } catch (JsonProcessingException e) {
            throw new IllegalStateException(
                    "failed to deserialize " + type.getSimpleName() + ": " + e.getMessage(), e);
        }
    }

    public static <T> T treeToValue(JsonNode node, TypeReference<T> type) {
        if (node == null || node.isMissingNode() || node.isNull()) return null;
        try {
            return MAPPER.readValue(MAPPER.treeAsTokens(node), type);
        } catch (java.io.IOException e) {
            throw new IllegalStateException("failed to deserialize: " + e.getMessage(), e);
        }
    }

    public static byte[] toBytes(Object value) {
        try {
            return MAPPER.writeValueAsBytes(value);
        } catch (JsonProcessingException e) {
            throw new IllegalStateException("failed to serialize " + value.getClass(), e);
        }
    }
}
