package com.tesote.sdk.v1;

import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.Transport;
import com.tesote.sdk.internal.Json;
import com.tesote.sdk.internal.QueryParams;
import com.tesote.sdk.models.Transaction;
import com.tesote.sdk.models.TransactionsPage;

import java.time.Duration;
import java.util.Objects;

/**
 * v1 transactions: list-for-account (cursor) and read-by-id.
 */
public final class TransactionsClient {
    private static final Duration GET_CACHE = Duration.ofMinutes(5);

    private final Transport transport;

    public TransactionsClient(Transport transport) { this.transport = transport; }

    public TransactionsPage listForAccount(String accountId) {
        return listForAccount(accountId, new ListParams());
    }

    public TransactionsPage listForAccount(String accountId, ListParams params) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(params, "params");
        Transport.Options opts = Transport.Options.get(
                "/v1/accounts/" + AccountsClient.encode(accountId) + "/transactions")
                .query(params.toQuery());
        JsonNode node = transport.request(opts);
        return Json.treeToValue(node, TransactionsPage.class);
    }

    public Transaction get(String id) {
        Objects.requireNonNull(id, "id");
        Transport.Options opts = Transport.Options.get(
                "/v1/transactions/" + AccountsClient.encode(id))
                .cacheTtl(GET_CACHE);
        return Json.treeToValue(transport.request(opts), Transaction.class);
    }

    /** Cursor + date filter parameters for {@link #listForAccount(String, ListParams)}. */
    public static final class ListParams {
        public String startDate;
        public String endDate;
        public String scope;
        public Integer page;
        public Integer perPage;
        public String transactionsAfterId;
        public String transactionsBeforeId;

        public ListParams startDate(String v) { this.startDate = v; return this; }
        public ListParams endDate(String v) { this.endDate = v; return this; }
        public ListParams scope(String v) { this.scope = v; return this; }
        public ListParams page(int v) { this.page = v; return this; }
        public ListParams perPage(int v) { this.perPage = v; return this; }
        public ListParams transactionsAfterId(String v) { this.transactionsAfterId = v; return this; }
        public ListParams transactionsBeforeId(String v) { this.transactionsBeforeId = v; return this; }

        java.util.Map<String, String> toQuery() {
            return QueryParams.of()
                    .put("start_date", startDate)
                    .put("end_date", endDate)
                    .put("scope", scope)
                    .put("page", page)
                    .put("per_page", perPage)
                    .put("transactions_after_id", transactionsAfterId)
                    .put("transactions_before_id", transactionsBeforeId)
                    .build();
        }
    }
}
