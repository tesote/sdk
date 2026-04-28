package com.tesote.sdk.v3;

import com.fasterxml.jackson.databind.JsonNode;
import okhttp3.mockwebserver.MockResponse;
import okhttp3.mockwebserver.MockWebServer;
import okhttp3.mockwebserver.RecordedRequest;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class AccountsTest {
    private MockWebServer server;
    private V3Client client;

    @BeforeEach
    void setUp() throws IOException {
        server = new MockWebServer();
        server.start();
        client = V3Client.builder()
                .apiKey("sk_test_abcd1234")
                .baseUrl(server.url("/api").toString())
                .build();
    }

    @AfterEach
    void tearDown() throws IOException {
        server.shutdown();
    }

    @Test
    void listHitsExpectedPath() throws Exception {
        server.enqueue(new MockResponse()
                .setResponseCode(200)
                .setHeader("Content-Type", "application/json")
                .setBody("{\"data\":[{\"id\":\"acct_1\"},{\"id\":\"acct_2\"}]}"));

        JsonNode response = client.accounts().list();
        assertEquals(2, response.get("data").size());

        RecordedRequest rr = server.takeRequest();
        assertEquals("GET", rr.getMethod());
        assertTrue(rr.getPath().endsWith("/v3/accounts"), rr.getPath());
    }

    @Test
    void getHitsResourcePath() throws Exception {
        server.enqueue(new MockResponse()
                .setResponseCode(200)
                .setHeader("Content-Type", "application/json")
                .setBody("{\"id\":\"acct_42\"}"));

        JsonNode response = client.accounts().get("acct_42");
        assertEquals("acct_42", response.get("id").asText());

        RecordedRequest rr = server.takeRequest();
        assertEquals("GET", rr.getMethod());
        assertTrue(rr.getPath().endsWith("/v3/accounts/acct_42"), rr.getPath());
    }

    @Test
    void unwiredMethodsThrowUnsupported() {
        assertThrows(UnsupportedOperationException.class,
                () -> client.accounts().sync("acct_1"));
        assertThrows(UnsupportedOperationException.class, client::transactions);
    }
}
