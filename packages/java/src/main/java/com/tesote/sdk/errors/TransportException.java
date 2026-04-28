package com.tesote.sdk.errors;

/**
 * Lower-level failure with no usable HTTP response (DNS, TLS, timeout, reset).
 */
public class TransportException extends TesoteException {
    public TransportException(String message, String errorCode, int httpStatus,
                              String requestId, String errorId, Integer retryAfter,
                              String responseBody, RequestSummary requestSummary,
                              int attempts, Throwable cause) {
        super(message, errorCode, httpStatus, requestId, errorId,
                retryAfter, responseBody, requestSummary, attempts, cause);
    }
}
