package com.tesote.sdk.v2;

import com.tesote.sdk.Transport;
import com.tesote.sdk.errors.InvalidOrderStateException;
import com.tesote.sdk.errors.TransactionOrderNotFoundException;
import com.tesote.sdk.errors.ValidationException;
import com.tesote.sdk.models.OffsetPage;
import com.tesote.sdk.models.TransactionOrder;
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

class V2TransactionOrdersClientTest {
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
    void listOffsetPagination() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"items\":[{\"id\":\"o1\",\"status\":\"draft\"}],"
                        + "\"limit\":50,\"offset\":0,\"has_more\":false}"));

        OffsetPage<TransactionOrder> page = client.transactionOrders().list("a1",
                new TransactionOrdersClient.ListParams().limit(50).offset(0).status("draft"));

        assertEquals(1, page.items().size());
        assertEquals("o1", page.items().get(0).id());

        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().contains("/v2/accounts/a1/transaction_orders"));
        assertTrue(rr.getPath().contains("status=draft"));
    }

    @Test
    void getReturnsOrder() {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"id\":\"o1\",\"status\":\"draft\",\"amount\":100,\"currency\":\"VES\"}"));
        TransactionOrder o = client.transactionOrders().get("a1", "o1");
        assertEquals("o1", o.id());
    }

    @Test
    void createSendsEnvelope() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(201)
                .setBody("{\"id\":\"o1\",\"status\":\"draft\",\"amount\":100,\"currency\":\"VES\"}"));

        TransactionOrder o = client.transactionOrders().create("a1",
                new TransactionOrdersClient.CreateRequest()
                        .amount("100")
                        .currency("VES")
                        .description("rent")
                        .idempotencyKey("k_1")
                        .beneficiary(new TransactionOrdersClient.Beneficiary()
                                .name("Alice").bankCode("0102")));
        assertEquals("o1", o.id());

        RecordedRequest rr = server.takeRequest();
        assertEquals("POST", rr.getMethod());
        String body = rr.getBody().readUtf8();
        assertTrue(body.contains("\"transaction_order\""));
        assertTrue(body.contains("\"amount\":\"100\""));
        assertTrue(body.contains("\"currency\":\"VES\""));
        assertTrue(body.contains("\"bank_code\":\"0102\""));
        assertEquals("k_1", rr.getHeader("Idempotency-Key"));
    }

    @Test
    void submitWithToken() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(202)
                .setBody("{\"id\":\"o1\",\"status\":\"processing\"}"));

        TransactionOrder o = client.transactionOrders().submit("a1", "o1", "123456");
        assertEquals("processing", o.status());

        RecordedRequest rr = server.takeRequest();
        assertEquals("/api/v2/accounts/a1/transaction_orders/o1/submit", rr.getPath());
        String body = rr.getBody().readUtf8();
        assertTrue(body.contains("\"token\":\"123456\""));
    }

    @Test
    void submitWithoutTokenSendsEmptyJsonBody() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(202)
                .setBody("{\"id\":\"o1\",\"status\":\"processing\"}"));

        client.transactionOrders().submit("a1", "o1");

        RecordedRequest rr = server.takeRequest();
        assertEquals("application/json", rr.getHeader("Content-Type"));
    }

    @Test
    void cancelSendsPostWithEmptyBody() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"id\":\"o1\",\"status\":\"cancelled\"}"));

        TransactionOrder o = client.transactionOrders().cancel("a1", "o1");
        assertEquals("cancelled", o.status());

        RecordedRequest rr = server.takeRequest();
        assertEquals("POST", rr.getMethod());
        assertEquals("/api/v2/accounts/a1/transaction_orders/o1/cancel", rr.getPath());
        // why: spec says POST/PUT/PATCH must always carry Content-Type.
        assertEquals("application/json", rr.getHeader("Content-Type"));
        assertNotNull(rr.getHeader("Idempotency-Key"));
    }

    @Test
    void invalidOrderStateOnSubmit() {
        server.enqueue(new MockResponse().setResponseCode(409)
                .setBody("{\"error\":\"bad state\",\"error_code\":\"INVALID_ORDER_STATE\"}"));

        assertThrows(InvalidOrderStateException.class,
                () -> client.transactionOrders().submit("a1", "o1"));
    }

    @Test
    void notFoundOnGet() {
        server.enqueue(new MockResponse().setResponseCode(404)
                .setBody("{\"error\":\"x\",\"error_code\":\"TRANSACTION_ORDER_NOT_FOUND\"}"));

        assertThrows(TransactionOrderNotFoundException.class,
                () -> client.transactionOrders().get("a1", "missing"));
    }

    @Test
    void validationOnCreate() {
        server.enqueue(new MockResponse().setResponseCode(400)
                .setBody("{\"error\":\"bad amount\",\"error_code\":\"VALIDATION_ERROR\"}"));

        assertThrows(ValidationException.class,
                () -> client.transactionOrders().create("a1",
                        new TransactionOrdersClient.CreateRequest().amount("-1")));
    }
}
