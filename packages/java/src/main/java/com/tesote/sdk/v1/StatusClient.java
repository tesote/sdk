package com.tesote.sdk.v1;

import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.Transport;
import com.tesote.sdk.internal.Json;
import com.tesote.sdk.models.Status;
import com.tesote.sdk.models.Whoami;

/**
 * v1 status + whoami. {@code status()} works without an API key.
 */
public final class StatusClient {
    private final Transport transport;

    public StatusClient(Transport transport) { this.transport = transport; }

    /** Public health check. Auth not required server-side. */
    public Status status() {
        JsonNode node = transport.request(Transport.Options.get("/status"));
        return Json.treeToValue(node, Status.class);
    }

    /** Identity of the API key owner. Requires auth. */
    public Whoami whoami() {
        JsonNode node = transport.request(Transport.Options.get("/whoami"));
        return Json.treeToValue(node, Whoami.class);
    }
}
