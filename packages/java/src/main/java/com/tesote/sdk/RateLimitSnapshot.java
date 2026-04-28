package com.tesote.sdk;

import java.time.Instant;

/**
 * Snapshot of the most recent {@code X-RateLimit-*} headers seen by the transport.
 *
 * @param limit     value of {@code X-RateLimit-Limit}, or -1 if absent
 * @param remaining value of {@code X-RateLimit-Remaining}, or -1 if absent
 * @param resetAt   parsed {@code X-RateLimit-Reset}; null when absent or unparseable
 */
public record RateLimitSnapshot(int limit, int remaining, Instant resetAt) {
    public static RateLimitSnapshot empty() {
        return new RateLimitSnapshot(-1, -1, null);
    }
}
