package com.tesote.sdk.errors;

public class UnprocessableContentException extends ApiException {
    public UnprocessableContentException(String message, String errorCode, int httpStatus,
                                         String requestId, String errorId, Integer retryAfter,
                                         String responseBody, RequestSummary requestSummary,
                                         int attempts, Throwable cause) {
        super(message, errorCode, httpStatus, requestId, errorId,
                retryAfter, responseBody, requestSummary, attempts, cause);
    }
}
