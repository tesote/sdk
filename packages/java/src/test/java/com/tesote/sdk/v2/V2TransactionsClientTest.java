package com.tesote.sdk.v2;

import com.tesote.sdk.Transport;
import com.tesote.sdk.errors.InvalidCountException;
import com.tesote.sdk.errors.UnprocessableContentException;
import com.tesote.sdk.errors.HistorySyncForbiddenException;
import com.tesote.sdk.models.BulkResponse;
import com.tesote.sdk.models.SyncTransactionsResponse;
import com.tesote.sdk.models.Transaction;
import com.tesote.sdk.models.TransactionsExport;
import com.tesote.sdk.models.TransactionsPage;
import com.tesote.sdk.models.TransactionsSearchResponse;
import okhttp3.mockwebserver.MockResponse;
import okhttp3.mockwebserver.MockWebServer;
import okhttp3.mockwebserver.RecordedRequest;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.time.Duration;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class V2TransactionsClientTest {
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
    void listForAccountSendsAllFilters() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"total\":0,\"transactions\":[],\"pagination\":{\"has_more\":false}}"));

        TransactionsClient.ListParams params = new TransactionsClient.ListParams()
                .startDate("2026-01-01").endDate("2026-01-31")
                .amountMin("10").amountMax("100")
                .status("posted")
                .categoryId("c1");

        TransactionsPage page = client.transactions().listForAccount("acct_1", params);
        assertNotNull(page);

        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().contains("/v2/accounts/acct_1/transactions"));
        assertTrue(rr.getPath().contains("amount_min=10"));
        assertTrue(rr.getPath().contains("amount_max=100"));
        assertTrue(rr.getPath().contains("status=posted"));
        assertTrue(rr.getPath().contains("category_id=c1"));
    }

    @Test
    void getReturnsTypedTransaction() {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"id\":\"t1\",\"status\":\"posted\","
                        + "\"data\":{\"amount_cents\":2500,\"currency\":\"VES\"}}"));
        Transaction t = client.transactions().get("t1");
        assertEquals("t1", t.id());
        assertEquals(2500L, t.data().amountCents());
    }

    @Test
    void syncSendsNestedOptionsAndDeserializes() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"added\":[{\"transaction_id\":\"t1\",\"account_id\":\"a1\","
                        + "\"amount\":1000.0,\"name\":\"x\"}],"
                        + "\"modified\":[],\"removed\":[],\"next_cursor\":\"c2\",\"has_more\":false}"));

        SyncTransactionsResponse resp = client.transactions().sync("a1",
                new TransactionsClient.SyncRequest()
                        .count(100).cursor("now").includeRunningBalance(true));

        assertEquals(1, resp.added().size());
        assertEquals("t1", resp.added().get(0).transactionId());
        assertEquals("c2", resp.nextCursor());

        RecordedRequest rr = server.takeRequest();
        assertEquals("POST", rr.getMethod());
        assertEquals("application/json", rr.getHeader("Content-Type"));
        String body = rr.getBody().readUtf8();
        assertTrue(body.contains("\"count\":100"), body);
        assertTrue(body.contains("\"cursor\":\"now\""), body);
        assertTrue(body.contains("\"include_running_balance\":true"), body);
    }

    @Test
    void syncLegacyHitsRootPath() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"added\":[],\"modified\":[],\"removed\":[],\"next_cursor\":null,\"has_more\":false}"));

        client.transactions().syncLegacy(new TransactionsClient.SyncRequest().count(50));

        RecordedRequest rr = server.takeRequest();
        assertEquals("/api/v2/transactions/sync", rr.getPath());
    }

    @Test
    void syncCountValidationMaps() {
        server.enqueue(new MockResponse().setResponseCode(422)
                .setBody("{\"error\":\"bad count\",\"error_code\":\"INVALID_COUNT\"}"));

        InvalidCountException ex = assertThrows(InvalidCountException.class,
                () -> client.transactions().sync("a1",
                        new TransactionsClient.SyncRequest().count(0)));
        assertEquals(422, ex.httpStatus());
    }

    @Test
    void historicalSyncForbidden() {
        server.enqueue(new MockResponse().setResponseCode(403)
                .setBody("{\"error\":\"too far back\",\"error_code\":\"HISTORY_SYNC_FORBIDDEN\"}"));
        assertThrows(HistorySyncForbiddenException.class,
                () -> client.transactions().sync("a1",
                        new TransactionsClient.SyncRequest().cursor("ancient")));
    }

    @Test
    void bulkSendsAccountIds() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"bulk_results\":[{\"account_id\":\"a1\","
                        + "\"transactions\":[],\"pagination\":{\"has_more\":false}}]}"));

        BulkResponse resp = client.transactions().bulk(
                new TransactionsClient.BulkRequest().accountIds(List.of("a1", "a2")).limit(50));
        assertEquals(1, resp.bulkResults().size());

        RecordedRequest rr = server.takeRequest();
        String body = rr.getBody().readUtf8();
        assertTrue(body.contains("\"a1\""));
        assertTrue(body.contains("\"a2\""));
        assertTrue(body.contains("\"limit\":50"));
    }

    @Test
    void bulkUnprocessableMapsTypedException() {
        server.enqueue(new MockResponse().setResponseCode(422)
                .setBody("{\"error\":\"too many\",\"error_code\":\"UNPROCESSABLE_CONTENT\"}"));

        assertThrows(UnprocessableContentException.class,
                () -> client.transactions().bulk(new TransactionsClient.BulkRequest()));
    }

    @Test
    void searchSendsQueryAndDeserializes() throws Exception {
        server.enqueue(new MockResponse().setResponseCode(200)
                .setBody("{\"transactions\":[{\"id\":\"t1\",\"status\":\"posted\","
                        + "\"data\":{\"amount_cents\":100}}],\"total\":1}"));

        TransactionsSearchResponse resp = client.transactions().search(
                new TransactionsClient.SearchParams().q("groceries").limit(25));

        assertEquals(1, resp.total());

        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().contains("q=groceries"));
        assertTrue(rr.getPath().contains("limit=25"));
    }

    @Test
    void exportReturnsRawBytesAndFilename() throws Exception {
        byte[] csv = "id,amount\nt1,100\n".getBytes();
        server.enqueue(new MockResponse().setResponseCode(200)
                .setHeader("Content-Type", "text/csv")
                .setHeader("Content-Disposition", "attachment; filename=\"transactions_a1_2026-04-28.csv\"")
                .setBody(new okio.Buffer().write(csv)));

        TransactionsExport export = client.transactions().export("a1",
                new TransactionsClient.ExportParams()
                        .format(TransactionsExport.Format.CSV));

        assertArrayEquals(csv, export.body());
        assertEquals("text/csv", export.contentType());
        assertEquals("transactions_a1_2026-04-28.csv", export.filename());

        RecordedRequest rr = server.takeRequest();
        assertTrue(rr.getPath().contains("/v2/accounts/a1/transactions/export"));
        assertTrue(rr.getPath().contains("format=csv"));
    }
}
