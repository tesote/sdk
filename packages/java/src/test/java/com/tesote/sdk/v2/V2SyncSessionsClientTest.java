package com.tesote.sdk.v2;

import com.tesote.sdk.errors.SyncSessionNotFoundException;
import com.tesote.sdk.models.SyncSession;
import com.tesote.sdk.models.SyncSessionsPage;
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

class V2SyncSessionsClientTest {
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
    void listAndPagination() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"sync_sessions\":[{\"id\":\"ss1\",\"status\":\"completed\"}],"
                        + "\"limit\":50,\"offset\":0,\"has_more\":false}"));

        SyncSessionsPage page = client.syncSessions().list("a1",
                new SyncSessionsClient.ListParams().limit(50).offset(0).status("completed"));

        assertEquals(1, page.syncSessions().size());
        assertEquals("ss1", page.syncSessions().get(0).id());

        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().contains("/v2/accounts/a1/sync_sessions"));
        assertTrue(rr.getPath().contains("status=completed"));
    }

    @Test
    void getReturnsSession() {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"id\":\"ss1\",\"status\":\"failed\","
                        + "\"error\":{\"type\":\"BankError\",\"message\":\"down\"}}"));

        SyncSession s = client.syncSessions().get("a1", "ss1");
        assertEquals("ss1", s.id());
        assertEquals("BankError", s.error().type());
    }

    @Test
    void notFoundMaps() {
        server.enqueue(new MockResponse().setResponseCode(404)
                .setBody("{\"error\":\"missing\",\"error_code\":\"SYNC_SESSION_NOT_FOUND\"}"));

        assertThrows(SyncSessionNotFoundException.class,
                () -> client.syncSessions().get("a1", "missing"));
    }
}
