package com.tesote.sdk.errors;

/**
 * Bad SDK configuration; raised at construction.
 */
public final class ConfigException extends TesoteException {
    public ConfigException(String message) {
        super(message, "CONFIG", 0, null, null, null, null, null, 0, null);
    }

    public ConfigException(String message, Throwable cause) {
        super(message, "CONFIG", 0, null, null, null, null, null, 0, cause);
    }
}
