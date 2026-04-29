package com.tesote.sdk.v2;

import com.fasterxml.jackson.core.type.TypeReference;
import com.tesote.sdk.Transport;
import com.tesote.sdk.internal.Json;
import com.tesote.sdk.internal.QueryParams;
import com.tesote.sdk.models.OffsetPage;
import com.tesote.sdk.models.PaymentMethod;

import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Objects;

/**
 * Payment methods (workspace-scoped). CRUD + soft-delete.
 */
public final class PaymentMethodsClient {
    private final Transport transport;

    public PaymentMethodsClient(Transport transport) { this.transport = transport; }

    public OffsetPage<PaymentMethod> list() { return list(new ListParams()); }

    public OffsetPage<PaymentMethod> list(ListParams params) {
        Objects.requireNonNull(params, "params");
        Transport.Options opts = Transport.Options.get("/v2/payment_methods")
                .query(QueryParams.of()
                        .put("limit", params.limit)
                        .put("offset", params.offset)
                        .put("method_type", params.methodType)
                        .put("currency", params.currency)
                        .put("counterparty_id", params.counterpartyId)
                        .put("verified", params.verified)
                        .build());
        return Json.treeToValue(transport.request(opts),
                new TypeReference<OffsetPage<PaymentMethod>>() {});
    }

    public PaymentMethod get(String id) {
        Objects.requireNonNull(id, "id");
        Transport.Options opts = Transport.Options.get(
                "/v2/payment_methods/" + AccountsClient.encode(id));
        return Json.treeToValue(transport.request(opts), PaymentMethod.class);
    }

    public PaymentMethod create(CreateRequest request) {
        Objects.requireNonNull(request, "request");
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v2/payment_methods";
        Map<String, Object> envelope = new LinkedHashMap<>();
        envelope.put("payment_method", request.toMap());
        opts.jsonBody(envelope);
        return Json.treeToValue(transport.request(opts), PaymentMethod.class);
    }

    public PaymentMethod update(String id, UpdateRequest request) {
        Objects.requireNonNull(id, "id");
        Objects.requireNonNull(request, "request");
        Transport.Options opts = new Transport.Options();
        opts.method = "PATCH";
        opts.path = "/v2/payment_methods/" + AccountsClient.encode(id);
        Map<String, Object> envelope = new LinkedHashMap<>();
        envelope.put("payment_method", request.toMap());
        opts.jsonBody(envelope);
        return Json.treeToValue(transport.request(opts), PaymentMethod.class);
    }

    public void delete(String id) {
        Objects.requireNonNull(id, "id");
        Transport.Options opts = new Transport.Options();
        opts.method = "DELETE";
        opts.path = "/v2/payment_methods/" + AccountsClient.encode(id);
        transport.request(opts);
    }

    public static final class ListParams {
        public Integer limit;
        public Integer offset;
        public String methodType;
        public String currency;
        public String counterpartyId;
        public Boolean verified;

        public ListParams limit(int v) { this.limit = v; return this; }
        public ListParams offset(int v) { this.offset = v; return this; }
        public ListParams methodType(String v) { this.methodType = v; return this; }
        public ListParams currency(String v) { this.currency = v; return this; }
        public ListParams counterpartyId(String v) { this.counterpartyId = v; return this; }
        public ListParams verified(boolean v) { this.verified = v; return this; }
    }

    /** Body for {@link #create(CreateRequest)}. */
    public static final class CreateRequest {
        public String methodType;
        public String currency;
        public String label;
        public String counterpartyId;
        public Counterparty counterparty;
        public Map<String, Object> details;

        public CreateRequest methodType(String v) { this.methodType = v; return this; }
        public CreateRequest currency(String v) { this.currency = v; return this; }
        public CreateRequest label(String v) { this.label = v; return this; }
        public CreateRequest counterpartyId(String v) { this.counterpartyId = v; return this; }
        public CreateRequest counterparty(Counterparty v) { this.counterparty = v; return this; }
        public CreateRequest details(Map<String, Object> v) { this.details = v; return this; }

        Map<String, Object> toMap() {
            Map<String, Object> m = new LinkedHashMap<>();
            if (methodType != null) m.put("method_type", methodType);
            if (currency != null) m.put("currency", currency);
            if (label != null) m.put("label", label);
            if (counterpartyId != null) m.put("counterparty_id", counterpartyId);
            if (counterparty != null) m.put("counterparty", counterparty.toMap());
            if (details != null) m.put("details", details);
            return m;
        }
    }

    /** Body for {@link #update(String, UpdateRequest)}; same shape as create, all optional. */
    public static final class UpdateRequest {
        public String methodType;
        public String currency;
        public String label;
        public String counterpartyId;
        public Counterparty counterparty;
        public Map<String, Object> details;

        public UpdateRequest methodType(String v) { this.methodType = v; return this; }
        public UpdateRequest currency(String v) { this.currency = v; return this; }
        public UpdateRequest label(String v) { this.label = v; return this; }
        public UpdateRequest counterpartyId(String v) { this.counterpartyId = v; return this; }
        public UpdateRequest counterparty(Counterparty v) { this.counterparty = v; return this; }
        public UpdateRequest details(Map<String, Object> v) { this.details = v; return this; }

        Map<String, Object> toMap() {
            Map<String, Object> m = new LinkedHashMap<>();
            if (methodType != null) m.put("method_type", methodType);
            if (currency != null) m.put("currency", currency);
            if (label != null) m.put("label", label);
            if (counterpartyId != null) m.put("counterparty_id", counterpartyId);
            if (counterparty != null) m.put("counterparty", counterparty.toMap());
            if (details != null) m.put("details", details);
            return m;
        }
    }

    public static final class Counterparty {
        public String name;

        public Counterparty name(String v) { this.name = v; return this; }

        Map<String, Object> toMap() {
            Map<String, Object> m = new LinkedHashMap<>();
            if (name != null) m.put("name", name);
            return m;
        }
    }
}
