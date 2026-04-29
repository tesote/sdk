package com.tesote.sdk.errors;

/**
 * Resource-not-found family. Concrete subclasses correspond to specific
 * {@code *_NOT_FOUND} error codes; callers may catch this base when the
 * specific resource doesn't matter.
 */
public class NotFoundException extends ApiException {
    public NotFoundException(String message, String errorCode, int httpStatus,
                             String requestId, String errorId, Integer retryAfter,
                             String responseBody, RequestSummary requestSummary,
                             int attempts, Throwable cause) {
        super(message, errorCode, httpStatus, requestId, errorId,
                retryAfter, responseBody, requestSummary, attempts, cause);
    }
}
