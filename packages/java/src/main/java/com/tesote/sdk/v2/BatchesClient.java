package com.tesote.sdk.v2;

import com.tesote.sdk.Transport;
import com.tesote.sdk.internal.Json;
import com.tesote.sdk.models.BatchActionResponse;
import com.tesote.sdk.models.BatchCreateResponse;
import com.tesote.sdk.models.BatchSummary;

import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;

/**
 * Multi-order batches scoped to a source account.
 */
public final class BatchesClient {
    private final Transport transport;

    public BatchesClient(Transport transport) { this.transport = transport; }

    public BatchCreateResponse create(String accountId, CreateRequest request) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(request, "request");
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v2/accounts/" + AccountsClient.encode(accountId) + "/batches";
        opts.jsonBody(request.toMap());
        return Json.treeToValue(transport.request(opts), BatchCreateResponse.class);
    }

    public BatchSummary show(String accountId, String batchId) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(batchId, "batchId");
        Transport.Options opts = Transport.Options.get(
                "/v2/accounts/" + AccountsClient.encode(accountId)
                        + "/batches/" + AccountsClient.encode(batchId));
        return Json.treeToValue(transport.request(opts), BatchSummary.class);
    }

    public BatchActionResponse approve(String accountId, String batchId) {
        return mutate(accountId, batchId, "approve", null);
    }

    public BatchActionResponse submit(String accountId, String batchId) {
        return submit(accountId, batchId, null);
    }

    public BatchActionResponse submit(String accountId, String batchId, String token) {
        Map<String, Object> body = new LinkedHashMap<>();
        if (token != null) body.put("token", token);
        return mutate(accountId, batchId, "submit", body);
    }

    public BatchActionResponse cancel(String accountId, String batchId) {
        return mutate(accountId, batchId, "cancel", null);
    }

    private BatchActionResponse mutate(String accountId, String batchId, String action,
                                       Map<String, Object> body) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(batchId, "batchId");
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v2/accounts/" + AccountsClient.encode(accountId)
                + "/batches/" + AccountsClient.encode(batchId) + "/" + action;
        if (body == null || body.isEmpty()) {
            opts.body = "{}".getBytes(StandardCharsets.UTF_8);
            opts.bodyShape = "0 fields";
        } else {
            opts.jsonBody(body);
        }
        return Json.treeToValue(transport.request(opts), BatchActionResponse.class);
    }

    public static final class CreateRequest {
        public List<TransactionOrdersClient.CreateRequest> orders = new ArrayList<>();

        public CreateRequest orders(List<TransactionOrdersClient.CreateRequest> v) {
            this.orders = v; return this;
        }

        public CreateRequest add(TransactionOrdersClient.CreateRequest order) {
            this.orders.add(order); return this;
        }

        Map<String, Object> toMap() {
            Map<String, Object> m = new LinkedHashMap<>();
            List<Map<String, Object>> wire = new ArrayList<>(orders.size());
            for (TransactionOrdersClient.CreateRequest o : orders) wire.add(o.toMap());
            m.put("orders", wire);
            return m;
        }
    }
}
