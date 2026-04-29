package com.tesote.sdk.v2;

import com.tesote.sdk.errors.SyncInProgressException;
import com.tesote.sdk.errors.SyncRateLimitExceededException;
import com.tesote.sdk.models.Account;
import com.tesote.sdk.models.AccountSyncResponse;
import com.tesote.sdk.models.AccountsPage;
import okhttp3.mockwebserver.MockResponse;
import okhttp3.mockwebserver.MockWebServer;
import okhttp3.mockwebserver.RecordedRequest;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.time.Duration;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class V2AccountsClientTest {
    private MockWebServer server;
    private V2Client client;

    @BeforeEach
    void setUp() throws IOException {
        server = new MockWebServer();
        server.start();
        client = V2Client.builder()
                .apiKey("sk_test_abcd1234")
                .baseUrl(server.url("/api").toString())
                .retryPolicy(new com.tesote.sdk.Transport.RetryPolicy(
                        1, Duration.ofMillis(1), Duration.ofMillis(2), false))
                .build();
    }

    @AfterEach
    void tearDown() throws IOException { server.shutdown(); }

    @Test
    void listSendsToV2Path() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"total\":0,\"accounts\":[],\"pagination\":{}}"));

        AccountsPage page = client.accounts().list();
        assertNotNull(page);
        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().startsWith("/api/v2/accounts"), rr.getPath());
    }

    @Test
    void getDeserializesV2Account() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"id\":\"a1\",\"name\":\"v2 acct\","
                        + "\"data\":{\"currency\":\"VES\",\"balance_cents\":\"1000\"}}"));

        Account a = client.accounts().get("a1");
        assertEquals("a1", a.id());
        assertEquals("1000", a.data().balanceCents());
    }

    @Test
    void syncReturns202Body() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(202)
                .setBody("{\"message\":\"Sync started\",\"sync_session_id\":\"ss_1\","
                        + "\"status\":\"pending\",\"started_at\":\"2026-04-28T19:21:00Z\"}"));

        AccountSyncResponse resp = client.accounts().sync("a1");
        assertEquals("ss_1", resp.syncSessionId());
        assertEquals("pending", resp.status());

        RecordedRequest rr = server.takeRequest();
        assertEquals("POST", rr.getMethod());
        assertEquals("/api/v2/accounts/a1/sync", rr.getPath());
        assertEquals("application/json", rr.getHeader("Content-Type"));
        assertNotNull(rr.getHeader("Idempotency-Key"));
    }

    @Test
    void syncIdempotencyKeyHonored() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(202)
                .setBody("{\"sync_session_id\":\"ss_1\",\"status\":\"pending\"}"));

        client.accounts().sync("a1", "my-key");
        RecordedRequest rr = server.takeRequest();
        assertEquals("my-key", rr.getHeader("Idempotency-Key"));
    }

    @Test
    void syncInProgressMaps() {
        server.enqueue(new MockResponse().setResponseCode(409)
                .setBody("{\"error\":\"sync running\",\"error_code\":\"SYNC_IN_PROGRESS\","
                        + "\"current_session_id\":\"ss_running\"}"));

        SyncInProgressException ex = assertThrows(SyncInProgressException.class,
                () -> client.accounts().sync("a1"));
        assertEquals(409, ex.httpStatus());
        assertTrue(ex.responseBody().contains("ss_running"));
    }

    @Test
    void syncRateLimitMapsToTypedException() {
        server.enqueue(new MockResponse().setResponseCode(429)
                .setHeader("Retry-After", "300")
                .setBody("{\"error\":\"too soon\",\"error_code\":\"SYNC_RATE_LIMIT_EXCEEDED\","
                        + "\"retry_after\":300}"));

        SyncRateLimitExceededException ex = assertThrows(SyncRateLimitExceededException.class,
                () -> client.accounts().sync("a1"));
        assertEquals(300, ex.retryAfter());
    }
}
