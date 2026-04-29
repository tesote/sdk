package com.tesote.sdk.v2;

import com.tesote.sdk.Transport;
import com.tesote.sdk.errors.BatchNotFoundException;
import com.tesote.sdk.errors.BatchValidationException;
import com.tesote.sdk.errors.InvalidOrderStateException;
import com.tesote.sdk.models.BatchActionResponse;
import com.tesote.sdk.models.BatchCreateResponse;
import com.tesote.sdk.models.BatchSummary;
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

class V2BatchesClientTest {
    private MockWebServer server;
    private V2Client client;

    @BeforeEach
    void setUp() throws IOException {
        server = new MockWebServer();
        server.start();
        client = V2Client.builder()
                .apiKey("sk_test_abcd1234")
                .baseUrl(server.url("/api").toString())
                .retryPolicy(new Transport.RetryPolicy(
                        1, Duration.ofMillis(1), Duration.ofMillis(2), false))
                .build();
    }

    @AfterEach
    void tearDown() throws IOException { server.shutdown(); }

    @Test
    void createReturnsBatchIdAndOrders() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(201)
                .setBody("{\"batch_id\":\"b1\",\"orders\":[{\"id\":\"o1\",\"status\":\"draft\"}],\"errors\":[]}"));

        BatchCreateResponse resp = client.batches().create("a1",
                new BatchesClient.CreateRequest()
                        .add(new TransactionOrdersClient.CreateRequest()
                                .amount("10").currency("VES").description("x")));
        assertEquals("b1", resp.batchId());
        assertEquals(1, resp.orders().size());

        RecordedRequest rr = server.takeRequest();
        assertEquals("POST", rr.getMethod());
        assertEquals("/api/v2/accounts/a1/batches", rr.getPath());
        String body = rr.getBody().readUtf8();
        assertTrue(body.contains("\"orders\""));
    }

    @Test
    void showReturnsSummary() {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"batch_id\":\"b1\",\"total_orders\":2,\"total_amount_cents\":1000,"
                        + "\"amount_currency\":\"VES\",\"statuses\":{\"draft\":2},"
                        + "\"batch_status\":\"draft\",\"orders\":[]}"));
        BatchSummary s = client.batches().show("a1", "b1");
        assertEquals("b1", s.batchId());
        assertEquals(2, s.totalOrders());
        assertEquals(2, s.statuses().get("draft"));
    }

    @Test
    void approveSendsPost() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"approved\":3,\"failed\":0}"));
        BatchActionResponse r = client.batches().approve("a1", "b1");
        assertEquals(3, r.approved());

        RecordedRequest rr = server.takeRequest();
        assertEquals("POST", rr.getMethod());
        assertEquals("/api/v2/accounts/a1/batches/b1/approve", rr.getPath());
        assertEquals("application/json", rr.getHeader("Content-Type"));
        assertNotNull(rr.getHeader("Idempotency-Key"));
    }

    @Test
    void submitWithToken() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"enqueued\":3,\"failed\":0}"));

        BatchActionResponse r = client.batches().submit("a1", "b1", "tok_42");
        assertEquals(3, r.enqueued());

        RecordedRequest rr = server.takeRequest();
        String body = rr.getBody().readUtf8();
        assertTrue(body.contains("tok_42"));
    }

    @Test
    void cancelSendsPost() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"cancelled\":2,\"skipped\":1,\"errors\":[]}"));

        BatchActionResponse r = client.batches().cancel("a1", "b1");
        assertEquals(2, r.cancelled());
        assertEquals(1, r.skipped());
    }

    @Test
    void batchValidationMaps() {
        server.enqueue(new MockResponse().setResponseCode(400)
                .setBody("{\"error\":\"bad\",\"error_code\":\"BATCH_VALIDATION_ERROR\"}"));

        assertThrows(BatchValidationException.class,
                () -> client.batches().create("a1", new BatchesClient.CreateRequest()));
    }

    @Test
    void invalidOrderStateOnApprove() {
        server.enqueue(new MockResponse().setResponseCode(409)
                .setBody("{\"error\":\"x\",\"error_code\":\"INVALID_ORDER_STATE\"}"));

        assertThrows(InvalidOrderStateException.class,
                () -> client.batches().approve("a1", "b1"));
    }

    @Test
    void notFoundOnShow() {
        server.enqueue(new MockResponse().setResponseCode(404)
                .setBody("{\"error\":\"x\",\"error_code\":\"BATCH_NOT_FOUND\"}"));

        assertThrows(BatchNotFoundException.class,
                () -> client.batches().show("a1", "missing"));
    }
}
