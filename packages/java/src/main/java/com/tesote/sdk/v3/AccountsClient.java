package com.tesote.sdk.v3;

import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.Transport;

import java.time.Duration;
import java.util.Objects;

/**
 * v3 accounts resource. Only {@code list} and {@code get} are wired for the
 * 0.1.0 bootstrap.
 */
public final class AccountsClient {
    private final Transport transport;

    AccountsClient(Transport transport) {
        this.transport = transport;
    }

    public JsonNode list() {
        return list(null);
    }

    /**
     * @param cacheTtl optional opt-in TTL for the response cache (null disables caching).
     */
    public JsonNode list(Duration cacheTtl) {
        Transport.Options o = Transport.Options.get(V3Client.VERSION_PATH + "/accounts");
        if (cacheTtl != null) o.cacheTtl(cacheTtl);
        return transport.request(o);
    }

    public JsonNode get(String id) {
        Objects.requireNonNull(id, "account id is required");
        Transport.Options o = Transport.Options.get(
                V3Client.VERSION_PATH + "/accounts/" + id);
        return transport.request(o);
    }

    public JsonNode sync(String id) {
        // why: stubbed; will return JsonNode once wired.
        throw new UnsupportedOperationException("not implemented");
    }
}
