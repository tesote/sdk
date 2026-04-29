package com.tesote.sdk.errors;

/**
 * 503 — the upstream bank reported maintenance. Retry after the configured
 * window. Distinct from the broader {@link ServiceUnavailableException}
 * (which signals API-side pause mode).
 */
public final class BankUnderMaintenanceException extends ApiException {
    public BankUnderMaintenanceException(String message, String errorCode, int httpStatus,
                                         String requestId, String errorId, Integer retryAfter,
                                         String responseBody, RequestSummary requestSummary,
                                         int attempts, Throwable cause) {
        super(message, errorCode, httpStatus, requestId, errorId,
                retryAfter, responseBody, requestSummary, attempts, cause);
    }
}
