package com.tesote.sdk.v2;

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
import static org.junit.jupiter.api.Assertions.assertTrue;

class V2StatusClientTest {
    private MockWebServer server;
    private V2Client client;

    @BeforeEach
    void setUp() throws IOException {
        server = new MockWebServer();
        server.start();
        client = V2Client.builder()
                .apiKey("sk_test_abcd1234")
                .baseUrl(server.url("/api").toString())
                .build();
    }

    @AfterEach
    void tearDown() throws IOException { server.shutdown(); }

    @Test
    void v2StatusUsesV2Path() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"status\":\"ok\",\"authenticated\":false}"));

        Status s = client.status().status();
        assertEquals("ok", s.status());
        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().endsWith("/v2/status"), rr.getPath());
    }

    @Test
    void v2WhoamiUsesV2Path() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"client\":{\"id\":\"c1\",\"name\":\"x\",\"type\":\"workspace\"}}"));

        Whoami w = client.status().whoami();
        assertEquals("c1", w.client().id());
        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().endsWith("/v2/whoami"), rr.getPath());
    }
}
