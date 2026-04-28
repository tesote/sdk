package com.tesote.sdk;

import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.errors.ApiException;
import com.tesote.sdk.errors.ConfigException;
import com.tesote.sdk.errors.ErrorDispatcher;
import com.tesote.sdk.errors.NetworkException;
import com.tesote.sdk.errors.RequestSummary;
import com.tesote.sdk.errors.TimeoutException;
import com.tesote.sdk.errors.TlsException;
import com.tesote.sdk.errors.TransportException;
import com.tesote.sdk.internal.Json;

import javax.net.ssl.SSLException;
import java.io.IOException;
import java.net.ConnectException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.UnknownHostException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.Instant;
import java.time.format.DateTimeParseException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.ThreadLocalRandom;
import java.util.concurrent.atomic.AtomicReference;
import java.util.function.Consumer;

/**
 * Single HTTP entry point for every resource client.
 *
 * <p>Owns: bearer injection, retries with exponential backoff + jitter,
 * rate-limit header capture, idempotency-key generation, request-id
 * propagation into thrown exceptions, and the opt-in TTL response cache.
 *
 * <p>Resource clients call {@link #request(Options)} and receive a parsed
 * {@link JsonNode} or a typed {@link com.tesote.sdk.errors.TesoteException}.
 */
public final class Transport {
    public static final String DEFAULT_BASE_URL = "https://equipo.tesote.com/api";
    public static final String SDK_VERSION = "0.1.0";

    private static final List<String> MUTATING_METHODS =
            List.of("POST", "PUT", "PATCH", "DELETE");

    private final HttpClient httpClient;
    private final String baseUrl;
    private final String apiKey;
    private final String userAgent;
    private final RetryPolicy retryPolicy;
    private final Duration requestTimeout;
    private final CacheBackend cacheBackend;
    private final Consumer<LogEvent> logger;
    private final AtomicReference<RateLimitSnapshot> lastRateLimit =
            new AtomicReference<>(RateLimitSnapshot.empty());

    private Transport(Builder b) {
        if (b.apiKey == null || b.apiKey.isBlank()) {
            throw new ConfigException("apiKey is required");
        }
        this.apiKey = b.apiKey;
        this.baseUrl = trimTrailingSlash(b.baseUrl == null ? DEFAULT_BASE_URL : b.baseUrl);
        this.userAgent = b.userAgent != null ? b.userAgent : defaultUserAgent();
        this.retryPolicy = b.retryPolicy == null ? RetryPolicy.defaults() : b.retryPolicy;
        this.requestTimeout = b.requestTimeout == null ? Duration.ofSeconds(30) : b.requestTimeout;
        this.cacheBackend = b.cacheBackend;
        this.logger = b.logger;

        HttpClient.Builder hb = b.httpClientBuilder != null
                ? b.httpClientBuilder
                : HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(5));
        this.httpClient = hb.build();
    }

    public static Builder builder() {
        return new Builder();
    }

    /** Last captured rate-limit snapshot. Empty values until the first request lands. */
    public RateLimitSnapshot lastRateLimit() {
        return lastRateLimit.get();
    }

    /**
     * Send a request, applying retries / rate-limit awareness / caching.
     *
     * @return parsed JSON, or a missing {@link com.fasterxml.jackson.databind.JsonNode} on 204.
     */
    public JsonNode request(Options opts) {
        Map<String, String> queryView = opts.query == null ? Map.of() : Map.copyOf(opts.query);
        RequestSummary summary = new RequestSummary(
                opts.method, opts.path, queryView,
                opts.bodyShape, redactBearer(apiKey));

        boolean cacheable = "GET".equals(opts.method)
                && cacheBackend != null
                && opts.cacheTtl != null
                && !opts.cacheTtl.isZero()
                && !opts.cacheTtl.isNegative();

        String cacheKey = cacheable ? cacheKey(opts) : null;
        if (cacheable) {
            Optional<byte[]> hit = cacheBackend.get(cacheKey);
            if (hit.isPresent()) return Json.parse(hit.get());
        }

        String idempotencyKey = opts.idempotencyKey;
        if (idempotencyKey == null && MUTATING_METHODS.contains(opts.method)) {
            idempotencyKey = UUID.randomUUID().toString();
        }

        IOException lastIo = null;
        ApiException lastApi = null;
        for (int attempt = 1; attempt <= retryPolicy.maxAttempts(); attempt++) {
            HttpRequest httpRequest = build(opts, idempotencyKey);
            HttpResponse<byte[]> response;
            try {
                response = httpClient.send(httpRequest, HttpResponse.BodyHandlers.ofByteArray());
            } catch (java.net.http.HttpTimeoutException e) {
                lastIo = e;
                if (logger != null) logger.accept(new LogEvent(summary, attempt, -1, e));
                if (!retryPolicy.retryOnNetwork() || attempt == retryPolicy.maxAttempts()) {
                    throw new TimeoutException("request timed out", summary, attempt, e);
                }
                sleep(backoff(attempt, null));
                continue;
            } catch (SSLException e) {
                throw new TlsException("TLS error: " + e.getMessage(), summary, attempt, e);
            } catch (ConnectException | UnknownHostException e) {
                lastIo = e;
                if (logger != null) logger.accept(new LogEvent(summary, attempt, -1, e));
                if (!retryPolicy.retryOnNetwork() || attempt == retryPolicy.maxAttempts()) {
                    throw new NetworkException(e.getMessage(), summary, attempt, e);
                }
                sleep(backoff(attempt, null));
                continue;
            } catch (IOException e) {
                lastIo = e;
                if (logger != null) logger.accept(new LogEvent(summary, attempt, -1, e));
                if (!retryPolicy.retryOnNetwork() || attempt == retryPolicy.maxAttempts()) {
                    throw new NetworkException(e.getMessage(), summary, attempt, e);
                }
                sleep(backoff(attempt, null));
                continue;
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                throw new TransportException("interrupted", "INTERRUPTED",
                        0, null, null, null, null, summary, attempt, e) {};
            }

            captureRateLimit(response);
            int status = response.statusCode();
            String requestId = header(response, "X-Request-Id");

            if (logger != null) logger.accept(new LogEvent(summary, attempt, status, null));

            if (status >= 200 && status < 300) {
                JsonNode parsed = status == 204
                        ? Json.MAPPER.missingNode()
                        : Json.parse(response.body());
                if (cacheable) {
                    cacheBackend.put(cacheKey, response.body() == null ? new byte[0] : response.body(),
                            opts.cacheTtl);
                }
                if (MUTATING_METHODS.contains(opts.method) && cacheBackend != null) {
                    // why: any mutation invalidates GET caches under the same path prefix.
                    cacheBackend.invalidatePrefix("GET " + opts.path);
                }
                return parsed;
            }

            ApiException api = buildApiException(opts, summary, response, requestId, attempt);
            lastApi = api;

            if (shouldRetry(status, attempt)) {
                Duration sleepFor = backoff(attempt, retryAfterSeconds(response));
                sleep(sleepFor);
                continue;
            }
            throw api;
        }

        // why: loop only exits via throw on terminal failure or return on success.
        // Reaching here means retries exhausted on a transient response/network.
        if (lastApi != null) throw lastApi;
        if (lastIo != null) throw new NetworkException("retries exhausted", summary,
                retryPolicy.maxAttempts(), lastIo);
        throw new NetworkException("unexpected transport state", summary,
                retryPolicy.maxAttempts(), null);
    }

    private HttpRequest build(Options opts, String idempotencyKey) {
        URI uri = buildUri(opts);
        HttpRequest.Builder rb = HttpRequest.newBuilder()
                .uri(uri)
                .timeout(requestTimeout)
                .header("Authorization", "Bearer " + apiKey)
                .header("Accept", "application/json")
                .header("User-Agent", userAgent);

        if (idempotencyKey != null) {
            rb.header("Idempotency-Key", idempotencyKey);
        }
        if (opts.extraHeaders != null) {
            opts.extraHeaders.forEach(rb::header);
        }

        HttpRequest.BodyPublisher body = HttpRequest.BodyPublishers.noBody();
        if (opts.body != null) {
            body = HttpRequest.BodyPublishers.ofByteArray(opts.body);
            rb.header("Content-Type", "application/json");
        }

        return rb.method(opts.method, body).build();
    }

    private URI buildUri(Options opts) {
        StringBuilder sb = new StringBuilder(baseUrl);
        if (!opts.path.startsWith("/")) sb.append('/');
        sb.append(opts.path);
        if (opts.query != null && !opts.query.isEmpty()) {
            sb.append('?');
            boolean first = true;
            for (Map.Entry<String, String> e : opts.query.entrySet()) {
                if (!first) sb.append('&');
                first = false;
                sb.append(URLEncoder.encode(e.getKey(), StandardCharsets.UTF_8));
                sb.append('=');
                sb.append(URLEncoder.encode(e.getValue(), StandardCharsets.UTF_8));
            }
        }
        return URI.create(sb.toString());
    }

    private boolean shouldRetry(int status, int attempt) {
        if (attempt >= retryPolicy.maxAttempts()) return false;
        return status == 429 || status == 502 || status == 503 || status == 504;
    }

    private Duration backoff(int attempt, Integer retryAfterSeconds) {
        if (retryAfterSeconds != null && retryAfterSeconds > 0) {
            return Duration.ofSeconds(retryAfterSeconds);
        }
        long base = retryPolicy.baseDelay().toMillis();
        long cap = retryPolicy.maxDelay().toMillis();
        long exp = (long) (base * Math.pow(2, attempt - 1));
        long capped = Math.min(cap, exp);
        long jitter = ThreadLocalRandom.current().nextLong(capped / 4 + 1);
        return Duration.ofMillis(capped + jitter);
    }

    private void sleep(Duration d) {
        try {
            Thread.sleep(d.toMillis());
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    private void captureRateLimit(HttpResponse<?> r) {
        int limit = parseInt(header(r, "X-RateLimit-Limit"), -1);
        int remaining = parseInt(header(r, "X-RateLimit-Remaining"), -1);
        Instant resetAt = parseResetAt(header(r, "X-RateLimit-Reset"));
        lastRateLimit.set(new RateLimitSnapshot(limit, remaining, resetAt));
    }

    private static Instant parseResetAt(String raw) {
        if (raw == null || raw.isBlank()) return null;
        try {
            long maybeEpoch = Long.parseLong(raw.trim());
            return Instant.ofEpochSecond(maybeEpoch);
        } catch (NumberFormatException ignored) {
            // why: server sometimes sends ISO-8601; fall through.
        }
        try {
            return Instant.parse(raw.trim());
        } catch (DateTimeParseException ignored) {
            return null;
        }
    }

    private Integer retryAfterSeconds(HttpResponse<?> r) {
        String h = header(r, "Retry-After");
        if (h == null || h.isBlank()) return null;
        try {
            return Integer.parseInt(h.trim());
        } catch (NumberFormatException e) {
            return null;
        }
    }

    private static int parseInt(String s, int dflt) {
        if (s == null) return dflt;
        try { return Integer.parseInt(s.trim()); }
        catch (NumberFormatException e) { return dflt; }
    }

    private static String header(HttpResponse<?> r, String name) {
        return r.headers().firstValue(name).orElse(null);
    }

    private ApiException buildApiException(Options opts, RequestSummary summary,
                                            HttpResponse<byte[]> response,
                                            String requestId, int attempt) {
        byte[] bytes = response.body() == null ? new byte[0] : response.body();
        String bodyStr = new String(bytes, StandardCharsets.UTF_8);
        JsonNode envelope = Json.parse(bytes);

        String message = textField(envelope, "error", "HTTP " + response.statusCode());
        String errorCode = textField(envelope, "error_code", null);
        String errorId = textField(envelope, "error_id", null);
        Integer retryAfter = retryAfterSeconds(response);
        if (retryAfter == null && envelope.hasNonNull("retry_after")) {
            retryAfter = envelope.get("retry_after").asInt();
        }

        return ErrorDispatcher.dispatch(
                message, errorCode, response.statusCode(), requestId, errorId,
                retryAfter, bodyStr, summary, attempt, null);
    }

    private static String textField(JsonNode node, String field, String dflt) {
        if (node == null || !node.hasNonNull(field)) return dflt;
        return node.get(field).asText(dflt);
    }

    private String cacheKey(Options opts) {
        StringBuilder sb = new StringBuilder("GET ").append(opts.path);
        if (opts.query != null && !opts.query.isEmpty()) {
            sb.append('?');
            opts.query.entrySet().stream()
                    .sorted(Map.Entry.comparingByKey())
                    .forEach(e -> sb.append(e.getKey()).append('=').append(e.getValue()).append('&'));
        }
        sb.append("|key=").append(apiKeyHash());
        return sb.toString();
    }

    private String apiKeyHash() {
        // why: scope cache by API-key hash to avoid cross-tenant bleed without
        // logging the bearer token itself.
        return Integer.toHexString(apiKey.hashCode());
    }

    static String redactBearer(String key) {
        if (key == null || key.length() < 4) return "Bearer ****";
        return "Bearer ****" + key.substring(key.length() - 4);
    }

    private static String trimTrailingSlash(String s) {
        return s.endsWith("/") ? s.substring(0, s.length() - 1) : s;
    }

    private static String defaultUserAgent() {
        String javaVersion = System.getProperty("java.version", "unknown");
        return "tesote-sdk-java/" + SDK_VERSION + " (java/" + javaVersion + ")";
    }

    /** Per-request options. Resource clients construct this for each call. */
    public static final class Options {
        public String method = "GET";
        public String path = "/";
        public Map<String, String> query;
        public byte[] body;
        public String bodyShape;
        public String idempotencyKey;
        public Duration cacheTtl;
        public Map<String, String> extraHeaders;

        public static Options get(String path) {
            Options o = new Options();
            o.method = "GET";
            o.path = path;
            return o;
        }

        public Options query(String k, String v) {
            if (query == null) query = new HashMap<>();
            query.put(k, v);
            return this;
        }

        public Options cacheTtl(Duration d) { this.cacheTtl = d; return this; }
        public Options idempotencyKey(String k) { this.idempotencyKey = k; return this; }
    }

    /** Configurable retry parameters. */
    public record RetryPolicy(int maxAttempts, Duration baseDelay,
                              Duration maxDelay, boolean retryOnNetwork) {
        public static RetryPolicy defaults() {
            return new RetryPolicy(3, Duration.ofMillis(250), Duration.ofSeconds(8), true);
        }
    }

    /** Single-event log payload. {@code error} is null on success. */
    public record LogEvent(RequestSummary request, int attempt, int status, Throwable error) {}

    public static final class Builder {
        private String apiKey;
        private String baseUrl;
        private String userAgent;
        private RetryPolicy retryPolicy;
        private Duration requestTimeout;
        private CacheBackend cacheBackend;
        private Consumer<LogEvent> logger;
        private HttpClient.Builder httpClientBuilder;

        public Builder apiKey(String v) { this.apiKey = v; return this; }
        public Builder baseUrl(String v) { this.baseUrl = v; return this; }
        public Builder userAgent(String v) { this.userAgent = v; return this; }
        public Builder retryPolicy(RetryPolicy v) { this.retryPolicy = v; return this; }
        public Builder requestTimeout(Duration v) { this.requestTimeout = v; return this; }
        public Builder cacheBackend(CacheBackend v) { this.cacheBackend = v; return this; }
        public Builder logger(Consumer<LogEvent> v) { this.logger = v; return this; }
        public Builder httpClientBuilder(HttpClient.Builder v) { this.httpClientBuilder = v; return this; }

        public Transport build() { return new Transport(this); }
    }

    // why: package-private accessor for test convenience without exposing fields.
    List<String> mutatingMethodsView() { return new ArrayList<>(MUTATING_METHODS); }
}
