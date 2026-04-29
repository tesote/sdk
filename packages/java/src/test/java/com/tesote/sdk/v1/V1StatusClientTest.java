package com.tesote.sdk.v1;

import com.tesote.sdk.models.Status;
import com.tesote.sdk.models.Whoami;
import okhttp3.mockwebserver.MockResponse;
import okhttp3.mockwebserver.MockWebServer;
import okhttp3.mockwebserver.RecordedRequest;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

class V1StatusClientTest {
    private MockWebServer server;
    private V1Client client;

    @BeforeEach
    void setUp() throws IOException {
        server = new MockWebServer();
        server.start();
        client = V1Client.builder()
                .apiKey("sk_test_abcd1234")
                .baseUrl(server.url("/api").toString())
                .build();
    }

    @AfterEach
    void tearDown() throws IOException { server.shutdown(); }

    @Test
    void statusReturnsTyped() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"status\":\"ok\",\"authenticated\":false}"));

        Status status = client.status().status();
        assertEquals("ok", status.status());
        assertFalse(status.authenticated());

        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().endsWith("/status"));
    }

    @Test
    void whoamiReturnsClient() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"client\":{\"id\":\"cli_1\",\"name\":\"Acme\",\"type\":\"workspace\"}}"));

        Whoami w = client.status().whoami();
        assertEquals("cli_1", w.client().id());
        assertEquals("workspace", w.client().type());
    }
}
