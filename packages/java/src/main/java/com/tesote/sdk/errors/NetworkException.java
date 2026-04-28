package com.tesote.sdk.errors;

public final class NetworkException extends TransportException {
    public NetworkException(String message, RequestSummary requestSummary,
                            int attempts, Throwable cause) {
        super(message, "NETWORK_ERROR", 0, null, null, null, null,
                requestSummary, attempts, cause);
    }
}
