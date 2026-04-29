package com.tesote.sdk.v2;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.Transport;
import com.tesote.sdk.internal.Json;
import com.tesote.sdk.internal.QueryParams;
import com.tesote.sdk.models.OffsetPage;
import com.tesote.sdk.models.TransactionOrder;

import java.nio.charset.StandardCharsets;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Objects;

/**
 * Manage transaction orders: read, create (draft), submit, cancel.
 */
public final class TransactionOrdersClient {
    private final Transport transport;

    public TransactionOrdersClient(Transport transport) { this.transport = transport; }

    public OffsetPage<TransactionOrder> list(String accountId) {
        return list(accountId, new ListParams());
    }

    public OffsetPage<TransactionOrder> list(String accountId, ListParams params) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(params, "params");
        Transport.Options opts = Transport.Options.get(
                "/v2/accounts/" + AccountsClient.encode(accountId) + "/transaction_orders")
                .query(QueryParams.of()
                        .put("limit", params.limit)
                        .put("offset", params.offset)
                        .put("status", params.status)
                        .put("created_after", params.createdAfter)
                        .put("created_before", params.createdBefore)
                        .put("batch_id", params.batchId)
                        .build());
        JsonNode node = transport.request(opts);
        return Json.treeToValue(node, new TypeReference<OffsetPage<TransactionOrder>>() {});
    }

    public TransactionOrder get(String accountId, String orderId) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(orderId, "orderId");
        Transport.Options opts = Transport.Options.get(
                "/v2/accounts/" + AccountsClient.encode(accountId)
                        + "/transaction_orders/" + AccountsClient.encode(orderId));
        return Json.treeToValue(transport.request(opts), TransactionOrder.class);
    }

    /**
     * Create a draft order. If {@code request.idempotencyKey} is set the
     * server returns the existing order on retry; we also forward it as the
     * {@code Idempotency-Key} HTTP header so transport-level retries are safe.
     */
    public TransactionOrder create(String accountId, CreateRequest request) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(request, "request");
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v2/accounts/" + AccountsClient.encode(accountId) + "/transaction_orders";
        Map<String, Object> envelope = new LinkedHashMap<>();
        envelope.put("transaction_order", request.toMap());
        opts.jsonBody(envelope);
        if (request.idempotencyKey != null) {
            opts.idempotencyKey(request.idempotencyKey);
        }
        return Json.treeToValue(transport.request(opts), TransactionOrder.class);
    }

    public TransactionOrder submit(String accountId, String orderId) {
        return submit(accountId, orderId, null);
    }

    public TransactionOrder submit(String accountId, String orderId, String token) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(orderId, "orderId");
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v2/accounts/" + AccountsClient.encode(accountId)
                + "/transaction_orders/" + AccountsClient.encode(orderId) + "/submit";
        Map<String, Object> body = new LinkedHashMap<>();
        if (token != null) body.put("token", token);
        opts.jsonBody(body);
        return Json.treeToValue(transport.request(opts), TransactionOrder.class);
    }

    public TransactionOrder cancel(String accountId, String orderId) {
        Objects.requireNonNull(accountId, "accountId");
        Objects.requireNonNull(orderId, "orderId");
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v2/accounts/" + AccountsClient.encode(accountId)
                + "/transaction_orders/" + AccountsClient.encode(orderId) + "/cancel";
        opts.body = "{}".getBytes(StandardCharsets.UTF_8);
        opts.bodyShape = "0 fields";
        return Json.treeToValue(transport.request(opts), TransactionOrder.class);
    }

    public static final class ListParams {
        public Integer limit;
        public Integer offset;
        public String status;
        public String createdAfter;
        public String createdBefore;
        public String batchId;

        public ListParams limit(int v) { this.limit = v; return this; }
        public ListParams offset(int v) { this.offset = v; return this; }
        public ListParams status(String v) { this.status = v; return this; }
        public ListParams createdAfter(String v) { this.createdAfter = v; return this; }
        public ListParams createdBefore(String v) { this.createdBefore = v; return this; }
        public ListParams batchId(String v) { this.batchId = v; return this; }
    }

    /** Body for {@link #create(String, CreateRequest)}. */
    public static final class CreateRequest {
        public String destinationPaymentMethodId;
        public Beneficiary beneficiary;
        public String amount;
        public String currency;
        public String description;
        public String scheduledFor;
        public String idempotencyKey;
        public Map<String, Object> metadata;

        public CreateRequest destinationPaymentMethodId(String v) {
            this.destinationPaymentMethodId = v; return this;
        }
        public CreateRequest beneficiary(Beneficiary v) { this.beneficiary = v; return this; }
        public CreateRequest amount(String v) { this.amount = v; return this; }
        public CreateRequest currency(String v) { this.currency = v; return this; }
        public CreateRequest description(String v) { this.description = v; return this; }
        public CreateRequest scheduledFor(String v) { this.scheduledFor = v; return this; }
        public CreateRequest idempotencyKey(String v) { this.idempotencyKey = v; return this; }
        public CreateRequest metadata(Map<String, Object> v) { this.metadata = v; return this; }

        Map<String, Object> toMap() {
            Map<String, Object> m = new LinkedHashMap<>();
            if (destinationPaymentMethodId != null) {
                m.put("destination_payment_method_id", destinationPaymentMethodId);
            }
            if (beneficiary != null) m.put("beneficiary", beneficiary.toMap());
            if (amount != null) m.put("amount", amount);
            if (currency != null) m.put("currency", currency);
            if (description != null) m.put("description", description);
            if (scheduledFor != null) m.put("scheduled_for", scheduledFor);
            if (idempotencyKey != null) m.put("idempotency_key", idempotencyKey);
            if (metadata != null) m.put("metadata", metadata);
            return m;
        }
    }

    /** Inline beneficiary (creates a payment method server-side on first use). */
    public static final class Beneficiary {
        public String name;
        public String bankCode;
        public String accountNumber;
        public String identificationType;
        public String identificationNumber;

        public Beneficiary name(String v) { this.name = v; return this; }
        public Beneficiary bankCode(String v) { this.bankCode = v; return this; }
        public Beneficiary accountNumber(String v) { this.accountNumber = v; return this; }
        public Beneficiary identificationType(String v) { this.identificationType = v; return this; }
        public Beneficiary identificationNumber(String v) { this.identificationNumber = v; return this; }

        Map<String, Object> toMap() {
            Map<String, Object> m = new LinkedHashMap<>();
            if (name != null) m.put("name", name);
            if (bankCode != null) m.put("bank_code", bankCode);
            if (accountNumber != null) m.put("account_number", accountNumber);
            if (identificationType != null) m.put("identification_type", identificationType);
            if (identificationNumber != null) m.put("identification_number", identificationNumber);
            return m;
        }
    }
}
