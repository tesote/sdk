package com.tesote.sdk.errors;

public final class TimeoutException extends TransportException {
    public TimeoutException(String message, RequestSummary requestSummary,
                            int attempts, Throwable cause) {
        super(message, "TIMEOUT", 0, null, null, null, null,
                requestSummary, attempts, cause);
    }
}
