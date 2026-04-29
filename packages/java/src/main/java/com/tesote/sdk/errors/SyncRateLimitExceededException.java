package com.tesote.sdk.errors;

/**
 * 429 — per-bank-connection sync was triggered too soon after the previous
 * one. Distinct from the per-API-key {@link RateLimitExceededException}.
 */
public final class SyncRateLimitExceededException extends ApiException {
    public SyncRateLimitExceededException(String message, String errorCode, int httpStatus,
                                          String requestId, String errorId, Integer retryAfter,
                                          String responseBody, RequestSummary requestSummary,
                                          int attempts, Throwable cause) {
        super(message, errorCode, httpStatus, requestId, errorId,
                retryAfter, responseBody, requestSummary, attempts, cause);
    }
}
