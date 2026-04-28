package com.tesote.sdk;

import com.fasterxml.jackson.databind.JsonNode;
import com.tesote.sdk.errors.ApiException;
import com.tesote.sdk.errors.ConfigException;
import com.tesote.sdk.errors.RateLimitExceededException;
import com.tesote.sdk.errors.UnauthorizedException;
import com.tesote.sdk.internal.InMemoryCacheBackend;
import okhttp3.mockwebserver.MockResponse;
import okhttp3.mockwebserver.MockWebServer;
import okhttp3.mockwebserver.RecordedRequest;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.time.Duration;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class TransportTest {
    private MockWebServer server;

    @BeforeEach
    void setUp() throws IOException {
        server = new MockWebServer();
        server.start();
    }

    @AfterEach
    void tearDown() throws IOException {
        server.shutdown();
    }

    private Transport client(Transport.RetryPolicy policy) {
        Transport.Builder b = Transport.builder()
                .apiKey("sk_test_abcd1234")
                .baseUrl(server.url("/api").toString())
                .requestTimeout(Duration.ofSeconds(2));
        if (policy != null) b.retryPolicy(policy);
        return b.build();
    }

    @Test
    void missingApiKeyThrowsConfigException() {
        ConfigException ex = assertThrows(ConfigException.class,
                () -> Transport.builder().apiKey("").build());
        assertEquals("CONFIG", ex.errorCode());
    }

    @Test
    void successfulGetReturnsParsedJsonAndInjectsBearer() throws Exception {
        server.enqueue(new MockResponse()
                .setResponseCode(200)
                .setHeader("Content-Type", "application/json")
                .setHeader("X-Request-Id", "req_123")
                .setBody("{\"data\":[{\"id\":\"acct_1\"}]}"));

        Transport t = client(null);
        JsonNode response = t.request(Transport.Options.get("/v3/accounts"));

        assertEquals("acct_1", response.get("data").get(0).get("id").asText());

        RecordedRequest rr = server.takeRequest();
        assertEquals("Bearer sk_test_abcd1234", rr.getHeader("Authorization"));
        assertEquals("application/json", rr.getHeader("Accept"));
        assertNotNull(rr.getHeader("User-Agent"));
        assertTrue(rr.getHeader("User-Agent").startsWith("tesote-sdk-java/"));
    }

    @Test
    void rateLimitHeadersCapturedIntoSnapshot() throws Exception {
        server.enqueue(new MockResponse()
                .setResponseCode(200)
                .setHeader("X-RateLimit-Limit", "200")
                .setHeader("X-RateLimit-Remaining", "199")
                .setHeader("X-RateLimit-Reset", "1700000000")
                .setBody("{}"));

        Transport t = client(null);
        t.request(Transport.Options.get("/v3/accounts"));

        RateLimitSnapshot snap = t.lastRateLimit();
        assertEquals(200, snap.limit());
        assertEquals(199, snap.remaining());
        assertEquals(1700000000L, snap.resetAt().getEpochSecond());
    }

    @Test
    void retriesOn503WithBackoffThenSucceeds() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(503).setBody("{}"));
        server.enqueue(new MockResponse().setResponseCode(503).setBody("{}"));
        server.enqueue(new MockResponse().setResponseCode(200).setBody("{\"ok\":true}"));

        Transport t = client(new Transport.RetryPolicy(
                3, Duration.ofMillis(1), Duration.ofMillis(5), true));
        JsonNode response = t.request(Transport.Options.get("/v3/accounts"));

        assertTrue(response.get("ok").asBoolean());
        assertEquals(3, server.getRequestCount());
    }

    @Test
    void retryExhaustionThrowsRateLimitWithAttempts() {
        for (int i = 0; i < 3; i++) {
            server.enqueue(new MockResponse()
                    .setResponseCode(429)
                    .setHeader("Retry-After", "0")
                    .setBody("{\"error\":\"rate limited\",\"error_code\":\"RATE_LIMIT_EXCEEDED\"}"));
        }

        Transport t = client(new Transport.RetryPolicy(
                3, Duration.ofMillis(1), Duration.ofMillis(5), true));

        RateLimitExceededException ex = assertThrows(RateLimitExceededException.class,
                () -> t.request(Transport.Options.get("/v3/accounts")));
        assertEquals(3, ex.attempts());
        assertEquals("RATE_LIMIT_EXCEEDED", ex.errorCode());
    }

    @Test
    void doesNotRetryOn4xxOtherThan429() {
        server.enqueue(new MockResponse()
                .setResponseCode(401)
                .setBody("{\"error\":\"bad key\",\"error_code\":\"UNAUTHORIZED\"}"));

        Transport t = client(null);
        UnauthorizedException ex = assertThrows(UnauthorizedException.class,
                () -> t.request(Transport.Options.get("/v3/accounts")));
        assertEquals(1, server.getRequestCount());
        assertEquals(1, ex.attempts());
    }

    @Test
    void requestIdAttachedToThrownException() {
        server.enqueue(new MockResponse()
                .setResponseCode(401)
                .setHeader("X-Request-Id", "req_abc")
                .setBody("{\"error\":\"bad key\",\"error_code\":\"UNAUTHORIZED\",\"error_id\":\"err_1\"}"));

        Transport t = client(null);
        UnauthorizedException ex = assertThrows(UnauthorizedException.class,
                () -> t.request(Transport.Options.get("/v3/accounts")));
        assertEquals("req_abc", ex.requestId());
        assertEquals("err_1", ex.errorId());
    }

    @Test
    void idempotencyKeyAutoGeneratedOnPost() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200).setBody("{}"));

        Transport t = client(null);
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v3/accounts/acct_1/sync";
        opts.body = "{}".getBytes();
        opts.bodyShape = "0 bytes";
        t.request(opts);

        RecordedRequest rr = server.takeRequest();
        String key = rr.getHeader("Idempotency-Key");
        assertNotNull(key);
        // Throws if not a UUID:
        UUID.fromString(key);
    }

    @Test
    void idempotencyKeyHonoredWhenProvided() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200).setBody("{}"));

        Transport t = client(null);
        Transport.Options opts = new Transport.Options();
        opts.method = "POST";
        opts.path = "/v3/accounts/acct_1/sync";
        opts.body = "{}".getBytes();
        opts.idempotencyKey = "my-key";
        t.request(opts);

        assertEquals("my-key", server.takeRequest().getHeader("Idempotency-Key"));
    }

    @Test
    void getDoesNotCarryIdempotencyKey() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200).setBody("{}"));

        Transport t = client(null);
        t.request(Transport.Options.get("/v3/accounts"));

        assertNull(server.takeRequest().getHeader("Idempotency-Key"));
    }

    @Test
    void cacheHitAvoidsSecondRequest() {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setHeader("Content-Type", "application/json")
                .setBody("{\"hits\":1}"));

        Transport t = Transport.builder()
                .apiKey("sk_test_abcd1234")
                .baseUrl(server.url("/api").toString())
                .cacheBackend(new InMemoryCacheBackend())
                .build();

        Transport.Options opts = Transport.Options.get("/v3/accounts")
                .cacheTtl(Duration.ofSeconds(30));

        JsonNode first = t.request(opts);
        JsonNode second = t.request(opts);

        assertEquals(1, server.getRequestCount());
        assertEquals(first.toString(), second.toString());
    }

    @Test
    void requestSummaryRedactsBearerToken() {
        server.enqueue(new MockResponse().setResponseCode(401)
                .setBody("{\"error\":\"x\",\"error_code\":\"UNAUTHORIZED\"}"));

        Transport t = client(null);
        UnauthorizedException ex = assertThrows(UnauthorizedException.class,
                () -> t.request(Transport.Options.get("/v3/accounts")));
        String redacted = ex.requestSummary().redactedAuthorization();
        assertTrue(redacted.startsWith("Bearer ****"));
        assertNotEquals(-1, redacted.indexOf("1234"));
        // why: original key must never appear verbatim.
        assertEquals(-1, redacted.indexOf("abcd1234"));
    }

    @Test
    void apiExceptionMessageIsHumanGreppable() {
        server.enqueue(new MockResponse()
                .setResponseCode(429)
                .setHeader("Retry-After", "42")
                .setHeader("X-Request-Id", "req_xyz")
                .setBody("{\"error\":\"Rate limit exceeded\",\"error_code\":\"RATE_LIMIT_EXCEEDED\"}"));

        Transport t = client(new Transport.RetryPolicy(
                1, Duration.ofMillis(1), Duration.ofMillis(2), true));

        ApiException ex = assertThrows(ApiException.class,
                () -> t.request(Transport.Options.get("/v3/accounts")));
        String msg = ex.getMessage();
        assertTrue(msg.contains("RateLimitExceededException"), msg);
        assertTrue(msg.contains("429"), msg);
        assertTrue(msg.contains("RATE_LIMIT_EXCEEDED"), msg);
        assertTrue(msg.contains("req_xyz"), msg);
        assertTrue(msg.contains("retry_after: 42"), msg);
    }
}
