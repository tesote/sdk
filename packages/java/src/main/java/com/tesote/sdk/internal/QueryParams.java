package com.tesote.sdk.internal;

import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Tiny builder for query-string maps. Skips null values so callers don't
 * have to gate every {@code .put} on a not-null check.
 */
public final class QueryParams {
    private final Map<String, String> values = new LinkedHashMap<>();

    public static QueryParams of() { return new QueryParams(); }

    public QueryParams put(String key, String value) {
        if (value != null) values.put(key, value);
        return this;
    }

    public QueryParams put(String key, Number value) {
        if (value != null) values.put(key, value.toString());
        return this;
    }

    public QueryParams put(String key, Boolean value) {
        if (value != null) values.put(key, value.toString());
        return this;
    }

    public Map<String, String> build() { return values; }

    public boolean isEmpty() { return values.isEmpty(); }
}
