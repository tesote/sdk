package com.tesote.sdk.v2;

import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.Transport;
import com.tesote.sdk.internal.Json;
import com.tesote.sdk.models.Status;
import com.tesote.sdk.models.Whoami;

/**
 * v2 status + whoami. Mirrors v1; separate path prefix.
 */
public final class StatusClient {
    private final Transport transport;

    public StatusClient(Transport transport) { this.transport = transport; }

    public Status status() {
        JsonNode node = transport.request(Transport.Options.get("/v2/status"));
        return Json.treeToValue(node, Status.class);
    }

    public Whoami whoami() {
        JsonNode node = transport.request(Transport.Options.get("/v2/whoami"));
        return Json.treeToValue(node, Whoami.class);
    }
}
