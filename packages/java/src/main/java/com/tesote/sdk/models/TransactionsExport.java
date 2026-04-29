package com.tesote.sdk.models;

/**
 * File export from {@code GET /v2/accounts/{id}/transactions/export}.
 * Holds raw bytes plus the server-suggested filename and content type.
 */
public record TransactionsExport(
        byte[] body,
        String filename,
        String contentType
) {
    public enum Format {
        CSV("csv"),
        JSON("json");

        private final String wire;

        Format(String wire) { this.wire = wire; }

        public String wire() { return wire; }
    }
}
