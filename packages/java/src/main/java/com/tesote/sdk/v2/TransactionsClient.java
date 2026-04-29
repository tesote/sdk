package com.tesote.sdk.v2;

import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.Transport;
import com.tesote.sdk.internal.Json;
import com.tesote.sdk.internal.QueryParams;
import com.tesote.sdk.models.BulkResponse;
import com.tesote.sdk.models.SyncTransactionsResponse;
import com.tesote.sdk.models.Transaction;
import com.tesote.sdk.models.TransactionsExport;
import com.tesote.sdk.models.TransactionsPage;
import com.tesote.sdk.models.TransactionsSearchResponse;

import java.time.Duration;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

/**
 * v2 transactions: nested list, sync (cursor + legacy), get, bulk read,
 * search, export.
 */
public final class TransactionsClient {
    private static final Duration GET_CACHE = Duration.ofMinutes(5);
    private static final Duration LIST_CACHE = Duration.ofMinutes(1);

    private final Transport transport;

    public TransactionsClient(Transport transport) { this.transport = transport; }

    public TransactionsPage listForAccount(String accountId) {
        return listForAccount(accountId, new ListParams());
    }

    public TransactionsPage listForAccount(String accountId, ListParams params) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(params, "params");
        Transport.Options opts = Transport.Options.get(
                "/v2/accounts/" + AccountsClient.encode(accountId) + "/transactions")
                .query(params.toQuery())
                .cacheTtl(LIST_CACHE);
        return Json.treeToValue(transport.request(opts), TransactionsPage.class);
    }

    public Transaction get(String id) {
        Objects.requireNonNull(id, "id");
        Transport.Options opts = Transport.Options.get(
                "/v2/transactions/" + AccountsClient.encode(id))
                .cacheTtl(GET_CACHE);
        return Json.treeToValue(transport.request(opts), Transaction.class);
    }

    /** Cursor-based incremental sync scoped to one account. */
    public SyncTransactionsResponse sync(String accountId, SyncRequest request) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(request, "request");
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v2/accounts/" + AccountsClient.encode(accountId) + "/transactions/sync";
        opts.jsonBody(request.toMap());
        return Json.treeToValue(transport.request(opts), SyncTransactionsResponse.class);
    }

    /** Legacy non-nested sync. Caller specifies account context inside the body. */
    public SyncTransactionsResponse syncLegacy(SyncRequest request) {
        Objects.requireNonNull(request, "request");
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v2/transactions/sync";
        opts.jsonBody(request.toMap());
        return Json.treeToValue(transport.request(opts), SyncTransactionsResponse.class);
    }

    /** Fetch transactions across multiple accounts in one round-trip. */
    public BulkResponse bulk(BulkRequest request) {
        Objects.requireNonNull(request, "request");
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v2/transactions/bulk";
        opts.jsonBody(request.toMap());
        return Json.treeToValue(transport.request(opts), BulkResponse.class);
    }

    public TransactionsSearchResponse search(SearchParams params) {
        Objects.requireNonNull(params, "params");
        Transport.Options opts = Transport.Options.get("/v2/transactions/search")
                .query(params.toQuery());
        return Json.treeToValue(transport.request(opts), TransactionsSearchResponse.class);
    }

    /**
     * Download a CSV or JSON export of transactions for an account. Returns
     * the raw bytes plus the server-provided filename when present.
     */
    public TransactionsExport export(String accountId, ExportParams params) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(params, "params");
        Transport.Options opts = new Transport.Options();
        opts.method = "GET";
        opts.path = "/v2/accounts/" + AccountsClient.encode(accountId) + "/transactions/export";
        opts.query = params.toQuery();
        Transport.RawResponse raw = transport.requestRaw(opts);
        return new TransactionsExport(
                raw.body(),
                parseFilename(raw.contentDisposition()),
                raw.contentType());
    }

    private static String parseFilename(String contentDisposition) {
        if (contentDisposition == null) return null;
        // why: "attachment; filename=foo.csv" or filename*= variants.
        for (String part : contentDisposition.split(";")) {
            String trim = part.trim();
            if (trim.startsWith("filename=")) {
                String v = trim.substring("filename=".length()).trim();
                if (v.startsWith("\"") && v.endsWith("\"") && v.length() >= 2) {
                    v = v.substring(1, v.length() - 1);
                }
                return v;
            }
        }
        return null;
    }

    /** Body for sync / syncLegacy. */
    public static final class SyncRequest {
        public Integer count;
        public String cursor;
        public Boolean includeRunningBalance;

        public SyncRequest count(int v) { this.count = v; return this; }
        public SyncRequest cursor(String v) { this.cursor = v; return this; }
        public SyncRequest includeRunningBalance(boolean v) { this.includeRunningBalance = v; return this; }

        Map<String, Object> toMap() {
            Map<String, Object> m = new LinkedHashMap<>();
            if (count != null) m.put("count", count);
            if (cursor != null) m.put("cursor", cursor);
            if (includeRunningBalance != null) {
                Map<String, Object> opts = new LinkedHashMap<>();
                opts.put("include_running_balance", includeRunningBalance);
                m.put("options", opts);
            }
            return m;
        }
    }

    /** Body for {@link #bulk(BulkRequest)}. */
    public static final class BulkRequest {
        public List<String> accountIds;
        public Integer page;
        public Integer perPage;
        public Integer limit;
        public Integer offset;

        public BulkRequest accountIds(List<String> v) { this.accountIds = v; return this; }
        public BulkRequest page(int v) { this.page = v; return this; }
        public BulkRequest perPage(int v) { this.perPage = v; return this; }
        public BulkRequest limit(int v) { this.limit = v; return this; }
        public BulkRequest offset(int v) { this.offset = v; return this; }

        Map<String, Object> toMap() {
            Map<String, Object> m = new LinkedHashMap<>();
            if (accountIds != null) m.put("account_ids", accountIds);
            if (page != null) m.put("page", page);
            if (perPage != null) m.put("per_page", perPage);
            if (limit != null) m.put("limit", limit);
            if (offset != null) m.put("offset", offset);
            return m;
        }
    }

    /** Search filters; reuses the listForAccount filter shape. */
    public static final class SearchParams {
        public String q;
        public String accountId;
        public Integer limit;
        public Integer offset;
        public ListParams filters = new ListParams();

        public SearchParams q(String v) { this.q = v; return this; }
        public SearchParams accountId(String v) { this.accountId = v; return this; }
        public SearchParams limit(int v) { this.limit = v; return this; }
        public SearchParams offset(int v) { this.offset = v; return this; }
        public SearchParams filters(ListParams v) { this.filters = v; return this; }

        Map<String, String> toQuery() {
            Map<String, String> base = filters == null ? Map.of() : filters.toQuery();
            Map<String, String> result = new LinkedHashMap<>(base);
            QueryParams qp = QueryParams.of()
                    .put("q", q)
                    .put("account_id", accountId)
                    .put("limit", limit)
                    .put("offset", offset);
            qp.build().forEach(result::put);
            return result;
        }
    }

    /** Export query params: list filters plus {@code format}. */
    public static final class ExportParams {
        public ListParams filters = new ListParams();
        public TransactionsExport.Format format;

        public ExportParams filters(ListParams v) { this.filters = v; return this; }
        public ExportParams format(TransactionsExport.Format f) { this.format = f; return this; }

        Map<String, String> toQuery() {
            Map<String, String> base = filters == null ? Map.of() : filters.toQuery();
            Map<String, String> result = new LinkedHashMap<>(base);
            if (format != null) result.put("format", format.wire());
            return result;
        }
    }

    /** Full v2 transaction list filter set (used by list, search, export). */
    public static final class ListParams {
        public String startDate;
        public String endDate;
        public String scope;
        public Integer page;
        public Integer perPage;
        public String transactionsAfterId;
        public String transactionsBeforeId;
        public String transactionDateAfter;
        public String transactionDateBefore;
        public String createdAfter;
        public String updatedAfter;
        public String amountMin;
        public String amountMax;
        public String amount;
        public String status;
        public String categoryId;
        public String counterpartyId;
        public String q;
        public String type;
        public String referenceCode;

        public ListParams startDate(String v) { this.startDate = v; return this; }
        public ListParams endDate(String v) { this.endDate = v; return this; }
        public ListParams scope(String v) { this.scope = v; return this; }
        public ListParams page(int v) { this.page = v; return this; }
        public ListParams perPage(int v) { this.perPage = v; return this; }
        public ListParams transactionsAfterId(String v) { this.transactionsAfterId = v; return this; }
        public ListParams transactionsBeforeId(String v) { this.transactionsBeforeId = v; return this; }
        public ListParams transactionDateAfter(String v) { this.transactionDateAfter = v; return this; }
        public ListParams transactionDateBefore(String v) { this.transactionDateBefore = v; return this; }
        public ListParams createdAfter(String v) { this.createdAfter = v; return this; }
        public ListParams updatedAfter(String v) { this.updatedAfter = v; return this; }
        public ListParams amountMin(String v) { this.amountMin = v; return this; }
        public ListParams amountMax(String v) { this.amountMax = v; return this; }
        public ListParams amount(String v) { this.amount = v; return this; }
        public ListParams status(String v) { this.status = v; return this; }
        public ListParams categoryId(String v) { this.categoryId = v; return this; }
        public ListParams counterpartyId(String v) { this.counterpartyId = v; return this; }
        public ListParams q(String v) { this.q = v; return this; }
        public ListParams type(String v) { this.type = v; return this; }
        public ListParams referenceCode(String v) { this.referenceCode = v; return this; }

        Map<String, String> toQuery() {
            return QueryParams.of()
                    .put("start_date", startDate)
                    .put("end_date", endDate)
                    .put("scope", scope)
                    .put("page", page)
                    .put("per_page", perPage)
                    .put("transactions_after_id", transactionsAfterId)
                    .put("transactions_before_id", transactionsBeforeId)
                    .put("transaction_date_after", transactionDateAfter)
                    .put("transaction_date_before", transactionDateBefore)
                    .put("created_after", createdAfter)
                    .put("updated_after", updatedAfter)
                    .put("amount_min", amountMin)
                    .put("amount_max", amountMax)
                    .put("amount", amount)
                    .put("status", status)
                    .put("category_id", categoryId)
                    .put("counterparty_id", counterpartyId)
                    .put("q", q)
                    .put("type", type)
                    .put("reference_code", referenceCode)
                    .build();
        }
    }
}
