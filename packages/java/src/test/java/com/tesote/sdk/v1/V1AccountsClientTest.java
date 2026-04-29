package com.tesote.sdk.v1;

import com.tesote.sdk.errors.AccountNotFoundException;
import com.tesote.sdk.errors.UnauthorizedException;
import com.tesote.sdk.models.Account;
import com.tesote.sdk.models.AccountsPage;
import okhttp3.mockwebserver.MockResponse;
import okhttp3.mockwebserver.MockWebServer;
import okhttp3.mockwebserver.RecordedRequest;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class V1AccountsClientTest {
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
    void listSendsPageQueryAndDeserializes() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"total\":1,\"accounts\":[{\"id\":\"a1\",\"name\":\"Bank A\","
                        + "\"data\":{\"currency\":\"VES\"}}],"
                        + "\"pagination\":{\"current_page\":1,\"per_page\":50,\"total_pages\":1,\"total_count\":1}}"));

        AccountsPage page = client.accounts().list(new AccountsClient.ListParams().page(1).perPage(50));

        assertEquals(1, page.total());
        assertEquals("a1", page.accounts().get(0).id());
        assertEquals("Bank A", page.accounts().get(0).name());
        assertEquals("VES", page.accounts().get(0).data().currency());

        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().contains("/v1/accounts"), rr.getPath());
        assertTrue(rr.getPath().contains("page=1"), rr.getPath());
        assertTrue(rr.getPath().contains("per_page=50"), rr.getPath());
    }

    @Test
    void getReturnsTypedAccount() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"id\":\"a1\",\"name\":\"Bank A\",\"data\":{\"currency\":\"USD\"},"
                        + "\"bank\":{\"name\":\"Banco Test\"}}"));

        Account acct = client.accounts().get("a1");
        assertEquals("a1", acct.id());
        assertEquals("Banco Test", acct.bank().name());

        RecordedRequest rr = server.takeRequest();
        assertEquals("/api/v1/accounts/a1", rr.getPath());
    }

    @Test
    void notFoundMapsToTypedException() {
        server.enqueue(new MockResponse().setResponseCode(404)
                .setBody("{\"error\":\"missing\",\"error_code\":\"ACCOUNT_NOT_FOUND\"}"));

        AccountNotFoundException ex = assertThrows(AccountNotFoundException.class,
                () -> client.accounts().get("a-missing"));
        assertEquals("ACCOUNT_NOT_FOUND", ex.errorCode());
        assertEquals(404, ex.httpStatus());
    }

    @Test
    void unauthorizedOnList() {
        server.enqueue(new MockResponse().setResponseCode(401)
                .setBody("{\"error\":\"bad key\",\"error_code\":\"UNAUTHORIZED\"}"));

        UnauthorizedException ex = assertThrows(UnauthorizedException.class,
                () -> client.accounts().list());
        assertEquals(401, ex.httpStatus());
    }

    @Test
    void listIdLookupSendsBearer() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"total\":0,\"accounts\":[],\"pagination\":{}}"));
        client.accounts().list();
        RecordedRequest rr = server.takeRequest();
        assertNotNull(rr.getHeader("Authorization"));
        assertTrue(rr.getHeader("Authorization").startsWith("Bearer "));
    }
}
