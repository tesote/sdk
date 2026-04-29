package com.tesote.sdk.v2;

import com.tesote.sdk.Transport;
import com.tesote.sdk.errors.PaymentMethodNotFoundException;
import com.tesote.sdk.errors.ValidationException;
import com.tesote.sdk.models.OffsetPage;
import com.tesote.sdk.models.PaymentMethod;
import okhttp3.mockwebserver.MockResponse;
import okhttp3.mockwebserver.MockWebServer;
import okhttp3.mockwebserver.RecordedRequest;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.time.Duration;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class V2PaymentMethodsClientTest {
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
    void listSendsFilters() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"items\":[{\"id\":\"pm1\",\"method_type\":\"bank_account\","
                        + "\"currency\":\"VES\"}],\"limit\":50,\"offset\":0,\"has_more\":false}"));

        OffsetPage<PaymentMethod> page = client.paymentMethods().list(
                new PaymentMethodsClient.ListParams()
                        .limit(50).offset(0)
                        .methodType("bank_account").verified(true));

        assertEquals(1, page.items().size());
        assertEquals("pm1", page.items().get(0).id());

        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().contains("method_type=bank_account"));
        assertTrue(rr.getPath().contains("verified=true"));
    }

    @Test
    void getReturnsTyped() {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"id\":\"pm1\",\"method_type\":\"bank_account\",\"currency\":\"VES\","
                        + "\"verified\":true}"));
        PaymentMethod pm = client.paymentMethods().get("pm1");
        assertEquals("pm1", pm.id());
        assertTrue(pm.verified());
    }

    @Test
    void createSendsEnvelope() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(201)
                .setBody("{\"id\":\"pm1\",\"method_type\":\"bank_account\",\"currency\":\"VES\"}"));

        PaymentMethod pm = client.paymentMethods().create(new PaymentMethodsClient.CreateRequest()
                .methodType("bank_account")
                .currency("VES")
                .label("Primary")
                .counterparty(new PaymentMethodsClient.Counterparty().name("Alice"))
                .details(Map.of("bank_code", "0102", "account_number", "1234")));

        assertEquals("pm1", pm.id());

        RecordedRequest rr = server.takeRequest();
        assertEquals("POST", rr.getMethod());
        String body = rr.getBody().readUtf8();
        assertTrue(body.contains("\"payment_method\""));
        assertTrue(body.contains("\"method_type\":\"bank_account\""));
        assertTrue(body.contains("\"name\":\"Alice\""));
    }

    @Test
    void updateSendsPatch() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"id\":\"pm1\",\"method_type\":\"bank_account\",\"label\":\"Renamed\"}"));

        PaymentMethod pm = client.paymentMethods().update("pm1",
                new PaymentMethodsClient.UpdateRequest().label("Renamed"));
        assertEquals("Renamed", pm.label());

        RecordedRequest rr = server.takeRequest();
        assertEquals("PATCH", rr.getMethod());
        assertEquals("/api/v2/payment_methods/pm1", rr.getPath());
        assertEquals("application/json", rr.getHeader("Content-Type"));
        assertNotNull(rr.getHeader("Idempotency-Key"));
    }

    @Test
    void deleteSendsDelete() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(204));

        client.paymentMethods().delete("pm1");

        RecordedRequest rr = server.takeRequest();
        assertEquals("DELETE", rr.getMethod());
        assertEquals("/api/v2/payment_methods/pm1", rr.getPath());
    }

    @Test
    void notFoundOnGet() {
        server.enqueue(new MockResponse().setResponseCode(404)
                .setBody("{\"error\":\"missing\",\"error_code\":\"PAYMENT_METHOD_NOT_FOUND\"}"));

        assertThrows(PaymentMethodNotFoundException.class,
                () -> client.paymentMethods().get("missing"));
    }

    @Test
    void validationOnCreate() {
        server.enqueue(new MockResponse().setResponseCode(400)
                .setBody("{\"error\":\"bad\",\"error_code\":\"VALIDATION_ERROR\"}"));

        assertThrows(ValidationException.class,
                () -> client.paymentMethods().create(new PaymentMethodsClient.CreateRequest()));
    }
}
