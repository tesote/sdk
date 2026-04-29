package com.tesote.sdk.v2;

import com.tesote.sdk.CacheBackend;
import com.tesote.sdk.Transport;

import java.net.http.HttpClient;
import java.time.Duration;
import java.util.function.Consumer;

/**
 * v2 client. Adds writes for payments + sync orchestration on top of v1.
 *
 * <p>Each accessor returns the same instance for the lifetime of the client.
 * The client is thread-safe; share it.
 */
public final class V2Client {
    static final String VERSION_PATH = "/v2";

    private final Transport transport;
    private final StatusClient status;
    private final AccountsClient accounts;
    private final TransactionsClient transactions;
    private final SyncSessionsClient syncSessions;
    private final TransactionOrdersClient transactionOrders;
    private final BatchesClient batches;
    private final PaymentMethodsClient paymentMethods;

    private V2Client(Builder b) {
        this.transport = b.transportBuilder.build();
        this.status = new StatusClient(this.transport);
        this.accounts = new AccountsClient(this.transport);
        this.transactions = new TransactionsClient(this.transport);
        this.syncSessions = new SyncSessionsClient(this.transport);
        this.transactionOrders = new TransactionOrdersClient(this.transport);
        this.batches = new BatchesClient(this.transport);
        this.paymentMethods = new PaymentMethodsClient(this.transport);
    }

    public static Builder builder() {
        return new Builder();
    }

    public Transport transport() { return transport; }

    public StatusClient status() { return status; }
    public AccountsClient accounts() { return accounts; }
    public TransactionsClient transactions() { return transactions; }
    public SyncSessionsClient syncSessions() { return syncSessions; }
    public TransactionOrdersClient transactionOrders() { return transactionOrders; }
    public BatchesClient batches() { return batches; }
    public PaymentMethodsClient paymentMethods() { return paymentMethods; }

    public static final class Builder {
        private final Transport.Builder transportBuilder = Transport.builder();

        public Builder apiKey(String v) { transportBuilder.apiKey(v); return this; }
        public Builder baseUrl(String v) { transportBuilder.baseUrl(v); return this; }
        public Builder userAgent(String v) { transportBuilder.userAgent(v); return this; }
        public Builder requestTimeout(Duration v) { transportBuilder.requestTimeout(v); return this; }
        public Builder retryPolicy(Transport.RetryPolicy v) { transportBuilder.retryPolicy(v); return this; }
        public Builder cacheBackend(CacheBackend v) { transportBuilder.cacheBackend(v); return this; }
        public Builder logger(Consumer<Transport.LogEvent> v) { transportBuilder.logger(v); return this; }
        public Builder httpClientBuilder(HttpClient.Builder v) {
            transportBuilder.httpClientBuilder(v); return this;
        }

        public V2Client build() { return new V2Client(this); }
    }
}
