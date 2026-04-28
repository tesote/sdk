package com.tesote.sdk.errors;

import java.util.Map;

/**
 * Redacted snapshot of an outbound request, safe to include in error output.
 *
 * <p>Bearer tokens are redacted to {@code Bearer <last4>} before being captured here;
 * the raw token must never appear in a {@code RequestSummary}.
 */
public record RequestSummary(
        String method,
        String path,
        Map<String, String> query,
        String bodyShape,
        String redactedAuthorization
) {
    public RequestSummary {
        // why: defensive copy keeps callers from mutating after construction.
        query = query == null ? Map.of() : Map.copyOf(query);
    }
}
