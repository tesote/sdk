package com.tesote.sdk.errors;

/**
 * 500 — server-side error. Includes {@link #errorId()} for support tickets.
 */
public final class InternalErrorException extends ApiException {
    public InternalErrorException(String message, String errorCode, int httpStatus,
                                  String requestId, String errorId, Integer retryAfter,
                                  String responseBody, RequestSummary requestSummary,
                                  int attempts, Throwable cause) {
        super(message, errorCode, httpStatus, requestId, errorId,
                retryAfter, responseBody, requestSummary, attempts, cause);
    }
}
