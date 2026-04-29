package com.tesote.sdk.v1;

import com.tesote.sdk.errors.AccountNotFoundException;
import com.tesote.sdk.errors.InvalidDateRangeException;
import com.tesote.sdk.errors.TransactionNotFoundException;
import com.tesote.sdk.models.Transaction;
import com.tesote.sdk.models.TransactionsPage;
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

class V1TransactionsClientTest {
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
    void listForAccountWithCursor() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"total\":1,\"transactions\":[{\"id\":\"t1\",\"status\":\"posted\","
                        + "\"data\":{\"amount_cents\":1500}}],"
                        + "\"pagination\":{\"has_more\":true,\"per_page\":50,\"after_id\":\"t1\",\"before_id\":\"t1\"}}"));

        TransactionsPage page = client.transactions().listForAccount("acct_1",
                new TransactionsClient.ListParams()
                        .startDate("2026-01-01")
                        .endDate("2026-01-31")
                        .perPage(50));

        assertEquals(1, page.transactions().size());
        assertEquals("t1", page.transactions().get(0).id());
        assertTrue(page.pagination().hasMore());
        assertEquals(1500L, page.transactions().get(0).data().amountCents());

        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().contains("/v1/accounts/acct_1/transactions"));
        assertTrue(rr.getPath().contains("start_date=2026-01-01"));
        assertTrue(rr.getPath().contains("per_page=50"));
    }

    @Test
    void getTransactionByIdMaps() {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"id\":\"t99\",\"status\":\"posted\","
                        + "\"data\":{\"amount_cents\":1,\"currency\":\"VES\",\"description\":\"x\"}}"));

        Transaction t = client.transactions().get("t99");
        assertEquals("t99", t.id());
        assertEquals("posted", t.status());
        assertEquals("VES", t.data().currency());
    }

    @Test
    void invalidDateRangeMapsTo422() {
        server.enqueue(new MockResponse().setResponseCode(422)
                .setBody("{\"error\":\"bad range\",\"error_code\":\"INVALID_DATE_RANGE\"}"));

        InvalidDateRangeException ex = assertThrows(InvalidDateRangeException.class,
                () -> client.transactions().listForAccount("acct_1",
                        new TransactionsClient.ListParams().startDate("2030-01-01").endDate("2020-01-01")));
        assertEquals("INVALID_DATE_RANGE", ex.errorCode());
    }

    @Test
    void accountNotFoundOnList() {
        server.enqueue(new MockResponse().setResponseCode(404)
                .setBody("{\"error\":\"x\",\"error_code\":\"ACCOUNT_NOT_FOUND\"}"));
        assertThrows(AccountNotFoundException.class,
                () -> client.transactions().listForAccount("missing"));
    }

    @Test
    void transactionNotFoundOnGet() {
        server.enqueue(new MockResponse().setResponseCode(404)
                .setBody("{\"error\":\"x\",\"error_code\":\"TRANSACTION_NOT_FOUND\"}"));
        assertThrows(TransactionNotFoundException.class,
                () -> client.transactions().get("missing"));
    }
}
