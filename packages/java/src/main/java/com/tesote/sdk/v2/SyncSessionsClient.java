package com.tesote.sdk.v2;

import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.Transport;
import com.tesote.sdk.internal.Json;
import com.tesote.sdk.internal.QueryParams;
import com.tesote.sdk.models.SyncSession;
import com.tesote.sdk.models.SyncSessionsPage;

import java.util.Objects;

/**
 * Read-only access to per-account sync sessions.
 */
public final class SyncSessionsClient {
    private final Transport transport;

    public SyncSessionsClient(Transport transport) { this.transport = transport; }

    public SyncSessionsPage list(String accountId) {
        return list(accountId, new ListParams());
    }

    public SyncSessionsPage list(String accountId, ListParams params) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(params, "params");
        Transport.Options opts = Transport.Options.get(
                "/v2/accounts/" + AccountsClient.encode(accountId) + "/sync_sessions")
                .query(QueryParams.of()
                        .put("limit", params.limit)
                        .put("offset", params.offset)
                        .put("status", params.status)
                        .build());
        JsonNode node = transport.request(opts);
        return Json.treeToValue(node, SyncSessionsPage.class);
    }

    public SyncSession get(String accountId, String sessionId) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(sessionId, "sessionId");
        Transport.Options opts = Transport.Options.get(
                "/v2/accounts/" + AccountsClient.encode(accountId)
                        + "/sync_sessions/" + AccountsClient.encode(sessionId));
        return Json.treeToValue(transport.request(opts), SyncSession.class);
    }

    public static final class ListParams {
        public Integer limit;
        public Integer offset;
        public String status;

        public ListParams limit(int v) { this.limit = v; return this; }
        public ListParams offset(int v) { this.offset = v; return this; }
        public ListParams status(String v) { this.status = v; return this; }
    }
}
