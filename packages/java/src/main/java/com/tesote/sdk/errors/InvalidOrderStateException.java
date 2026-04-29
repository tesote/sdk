package com.tesote.sdk.errors;

/**
 * 409 — order is not in a state that permits the attempted transition
 * (e.g., submitting a non-draft order, cancelling a completed order).
 */
public final class InvalidOrderStateException extends ApiException {
    public InvalidOrderStateException(String message, String errorCode, int httpStatus,
                                      String requestId, String errorId, Integer retryAfter,
                                      String responseBody, RequestSummary requestSummary,
                                      int attempts, Throwable cause) {
        super(message, errorCode, httpStatus, requestId, errorId,
                retryAfter, responseBody, requestSummary, attempts, cause);
    }
}
