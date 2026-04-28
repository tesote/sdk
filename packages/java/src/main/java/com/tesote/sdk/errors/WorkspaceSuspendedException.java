package com.tesote.sdk.errors;

public final class WorkspaceSuspendedException extends ApiException {
    public WorkspaceSuspendedException(String message, String errorCode, int httpStatus,
                                       String requestId, String errorId, Integer retryAfter,
                                       String responseBody, RequestSummary requestSummary,
                                       int attempts, Throwable cause) {
        super(message, errorCode, httpStatus, requestId, errorId,
                retryAfter, responseBody, requestSummary, attempts, cause);
    }
}
