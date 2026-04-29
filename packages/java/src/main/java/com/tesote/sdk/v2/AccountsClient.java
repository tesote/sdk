package com.tesote.sdk.v2;

import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.Transport;
import com.tesote.sdk.internal.Json;
import com.tesote.sdk.internal.QueryParams;
import com.tesote.sdk.models.Account;
import com.tesote.sdk.models.AccountSyncResponse;
import com.tesote.sdk.models.AccountsPage;

import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.Objects;

/**
 * v2 accounts: list + get (same shape as v1) plus a {@code sync} mutation.
 */
public final class AccountsClient {
    private static final Duration LIST_CACHE = Duration.ofMinutes(1);
    private static final Duration GET_CACHE = Duration.ofMinutes(5);

    private final Transport transport;

    public AccountsClient(Transport transport) { this.transport = transport; }

    public AccountsPage list() { return list(new ListParams()); }

    public AccountsPage list(ListParams params) {
        Objects.requireNonNull(params, "params");
        Transport.Options opts = Transport.Options.get("/v2/accounts")
                .query(QueryParams.of()
                        .put("page", params.page)
                        .put("per_page", params.perPage)
                        .put("include", params.include)
                        .put("sort", params.sort)
                        .build())
                .cacheTtl(LIST_CACHE);
        JsonNode node = transport.request(opts);
        return Json.treeToValue(node, AccountsPage.class);
    }

    public Account get(String id) {
        Objects.requireNonNull(id, "id");
        Transport.Options opts = Transport.Options.get("/v2/accounts/" + encode(id))
                .cacheTtl(GET_CACHE);
        return Json.treeToValue(transport.request(opts), Account.class);
    }

    /**
     * Trigger an async bank sync for the account. The server returns 202 with
     * a {@code sync_session_id} that can be polled via the SyncSessionsClient.
     *
     * @param idempotencyKey optional caller-supplied key for safe retries.
     */
    public AccountSyncResponse sync(String id, String idempotencyKey) {
        Objects.requireNonNull(id, "id");
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v2/accounts/" + encode(id) + "/sync";
        opts.body = "{}".getBytes(StandardCharsets.UTF_8);
        opts.bodyShape = "0 fields";
        if (idempotencyKey != null) opts.idempotencyKey(idempotencyKey);
        return Json.treeToValue(transport.request(opts), AccountSyncResponse.class);
    }

    public AccountSyncResponse sync(String id) { return sync(id, null); }

    static String encode(String segment) {
        return URLEncoder.encode(segment, StandardCharsets.UTF_8);
    }

    public static final class ListParams {
        public Integer page;
        public Integer perPage;
        public String include;
        public String sort;

        public ListParams page(int v) { this.page = v; return this; }
        public ListParams perPage(int v) { this.perPage = v; return this; }
        public ListParams include(String v) { this.include = v; return this; }
        public ListParams sort(String v) { this.sort = v; return this; }
    }
}
