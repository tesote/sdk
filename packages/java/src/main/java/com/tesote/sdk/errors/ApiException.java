package com.tesote.sdk.errors;

/**
 * Server returned a usable HTTP response with an error envelope. Concrete
 * subclasses correspond 1:1 to API {@code error_code} values.
 */
public class ApiException extends TesoteException {
    public ApiException(
            String message,
            String errorCode,
            int httpStatus,
            String requestId,
            String errorId,
            Integer retryAfter,
            String responseBody,
            RequestSummary requestSummary,
            int attempts,
            Throwable cause
    ) {
        super(message, errorCode, httpStatus, requestId, errorId,
                retryAfter, responseBody, requestSummary, attempts, cause);
    }
}
