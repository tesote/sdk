package com.tesote.sdk.v2;

import com.tesote.sdk.CacheBackend;
import com.tesote.sdk.Transport;

import java.net.http.HttpClient;
import java.time.Duration;
import java.util.function.Consumer;

/**
 * v2 client. Adds writes for payments + sync orchestration on top of v1.
 *
 * <p>0.1.0 ships the builder + transport plumbing; resource methods stub until
 * they're wired in subsequent commits.
 */
public final class V2Client {
    static final String VERSION_PATH = "/v2";

    private final Transport transport;

    private V2Client(Builder b) {
        this.transport = b.transportBuilder.build();
    }

    public static Builder builder() {
        return new Builder();
    }

    public Transport transport() { return transport; }

    public Object status() { throw new UnsupportedOperationException("not implemented"); }
    public Object accounts() { throw new UnsupportedOperationException("not implemented"); }
    public Object transactions() { throw new UnsupportedOperationException("not implemented"); }
    public Object syncSessions() { throw new UnsupportedOperationException("not implemented"); }
    public Object transactionOrders() { throw new UnsupportedOperationException("not implemented"); }
    public Object batches() { throw new UnsupportedOperationException("not implemented"); }
    public Object paymentMethods() { throw new UnsupportedOperationException("not implemented"); }

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
