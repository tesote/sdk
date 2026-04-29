package com.tesote.sdk.v1;

import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.Transport;
import com.tesote.sdk.internal.Json;
import com.tesote.sdk.internal.QueryParams;
import com.tesote.sdk.models.Account;
import com.tesote.sdk.models.AccountsPage;

import java.time.Duration;
import java.util.Objects;

/**
 * Read-only accounts on {@code /v1/accounts}.
 */
public final class AccountsClient {
    private static final Duration LIST_CACHE = Duration.ofMinutes(1);
    private static final Duration GET_CACHE = Duration.ofMinutes(5);

    private final Transport transport;

    public AccountsClient(Transport transport) { this.transport = transport; }

    public AccountsPage list() { return list(new ListParams()); }

    public AccountsPage list(ListParams params) {
        Objects.requireNonNull(params, "params");
        Transport.Options opts = Transport.Options.get("/v1/accounts")
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
        Transport.Options opts = Transport.Options.get("/v1/accounts/" + encode(id))
                .cacheTtl(GET_CACHE);
        return Json.treeToValue(transport.request(opts), Account.class);
    }

    static String encode(String segment) {
        return java.net.URLEncoder.encode(segment, java.nio.charset.StandardCharsets.UTF_8);
    }

    /** Optional filter / sort / pagination params for {@link #list(ListParams)}. */
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
