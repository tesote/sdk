package com.tesote.sdk.errors;

public final class TlsException extends TransportException {
    public TlsException(String message, RequestSummary requestSummary,
                        int attempts, Throwable cause) {
        super(message, "TLS_ERROR", 0, null, null, null, null,
                requestSummary, attempts, cause);
    }
}
