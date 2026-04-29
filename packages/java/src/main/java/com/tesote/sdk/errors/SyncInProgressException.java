package com.tesote.sdk.errors;

/**
 * 409 — a sync session for the bank connection is already running.
 * The active session id is in {@link #responseBody()} as
 * {@code current_session_id}.
 */
public final class SyncInProgressException extends ApiException {
    public SyncInProgressException(String message, String errorCode, int httpStatus,
                                   String requestId, String errorId, Integer retryAfter,
                                   String responseBody, RequestSummary requestSummary,
                                   int attempts, Throwable cause) {
        super(message, errorCode, httpStatus, requestId, errorId,
                retryAfter, responseBody, requestSummary, attempts, cause);
    }
}
