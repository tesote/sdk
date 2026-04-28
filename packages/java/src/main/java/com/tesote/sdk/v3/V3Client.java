package com.tesote.sdk.v3;

import com.tesote.sdk.CacheBackend;
import com.tesote.sdk.Transport;

import java.net.http.HttpClient;
import java.time.Duration;
import java.util.function.Consumer;

/**
 * v3 client. Wires {@code accounts().list()} and {@code accounts().get(id)}
 * end-to-end; other resources stub with {@link UnsupportedOperationException}
 * until they're wired in subsequent commits.
 */
public final class V3Client {
    static final String VERSION_PATH = "/v3";

    private final Transport transport;
    private final AccountsClient accounts;

    private V3Client(Builder b) {
        this.transport = b.transportBuilder.build();
        this.accounts = new AccountsClient(transport);
    }

    public static Builder builder() {
        return new Builder();
    }

    public AccountsClient accounts() { return accounts; }

    public Transport transport() { return transport; }

    // why: stub sub-clients live as inner interfaces returning UnsupportedOperationException
    // so the public surface is visible without scaffolding empty files.
    public Object transactions() { throw new UnsupportedOperationException("not implemented"); }
    public Object syncSessions() { throw new UnsupportedOperationException("not implemented"); }
    public Object transactionOrders() { throw new UnsupportedOperationException("not implemented"); }
    public Object batches() { throw new UnsupportedOperationException("not implemented"); }
    public Object paymentMethods() { throw new UnsupportedOperationException("not implemented"); }
    public Object categories() { throw new UnsupportedOperationException("not implemented"); }
    public Object counterparties() { throw new UnsupportedOperationException("not implemented"); }
    public Object legalEntities() { throw new UnsupportedOperationException("not implemented"); }
    public Object connections() { throw new UnsupportedOperationException("not implemented"); }
    public Object webhooks() { throw new UnsupportedOperationException("not implemented"); }
    public Object reports() { throw new UnsupportedOperationException("not implemented"); }
    public Object balanceHistory() { throw new UnsupportedOperationException("not implemented"); }
    public Object workspace() { throw new UnsupportedOperationException("not implemented"); }
    public Object mcp() { throw new UnsupportedOperationException("not implemented"); }
    public Object status() { throw new UnsupportedOperationException("not implemented"); }

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

        public V3Client build() { return new V3Client(this); }
    }
}
