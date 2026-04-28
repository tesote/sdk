package com.tesote.sdk.errors;

/**
 * Raised when calling a method whose upstream endpoint is gone in this API
 * version. The SDK keeps the method per the back-compat policy but throws this
 * pointing at the replacement.
 */
public final class EndpointRemovedException extends TesoteException {
    public EndpointRemovedException(String message) {
        super(message, "ENDPOINT_REMOVED", 0, null, null, null, null, null, 0, null);
    }
}
