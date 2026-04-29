package v2_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	tesote "github.com/tesote/sdk/go"
	"github.com/tesote/sdk/go/v2"
)

func TestTransactions_ListForAccount_FilterQuery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("amount_min") != "10.00" || q.Get("category_id") != "cat-1" {
			t.Errorf("query = %q", r.URL.RawQuery)
		}
		_, _ = io.WriteString(w, `{"total":0,"transactions":[],"pagination":{"has_more":false,"per_page":50,"after_id":"","before_id":""}}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.Transactions.ListForAccount(context.Background(), "a-1", v2.TransactionsFilter{
		AmountMin:  "10.00",
		CategoryID: "cat-1",
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestTransactions_Get_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"id":"tx-1","status":"posted","data":{"amount_cents":1000,"currency":"VES","description":"x","transaction_date":"2026-04-01"},"tesote_imported_at":"2026-04-01","tesote_updated_at":"2026-04-01","transaction_categories":[],"counterparty":null}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	out, err := c.Transactions.Get(context.Background(), "tx-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if out.ID != "tx-1" || out.Data.AmountCents != 1000 {
		t.Errorf("tx = %+v", out)
	}
}

func TestTransactions_Sync_Body(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("content-type = %q", ct)
		}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		_, _ = io.WriteString(w, `{"added":[],"modified":[],"removed":[],"next_cursor":null,"has_more":false}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.Transactions.Sync(context.Background(), "a-1", v2.SyncOptions{Count: 100, Cursor: "now"})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if got["count"].(float64) != 100 || got["cursor"] != "now" {
		t.Errorf("body = %+v", got)
	}
}

func TestTransactions_SyncLegacy_Path(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/transactions/sync" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"added":[],"modified":[],"removed":[],"next_cursor":null,"has_more":false}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.Transactions.SyncLegacy(context.Background(), v2.SyncOptions{})
	if err != nil {
		t.Fatalf("SyncLegacy: %v", err)
	}
}

func TestTransactions_Sync_HistoryForbidden(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, `{"error":"too far back","error_code":"HISTORY_SYNC_FORBIDDEN"}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.Transactions.Sync(context.Background(), "a-1", v2.SyncOptions{})
	if !errors.Is(err, tesote.ErrHistorySyncForbidden) {
		t.Errorf("err = %v", err)
	}
}

func TestTransactions_Sync_InvalidCount(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = io.WriteString(w, `{"error":"x","error_code":"INVALID_COUNT"}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.Transactions.Sync(context.Background(), "a-1", v2.SyncOptions{Count: 5000})
	var typed *tesote.InvalidCountError
	if !errors.As(err, &typed) {
		t.Fatalf("err = %T %v", err, err)
	}
}

func TestTransactions_Bulk_RequiresAccountIDs(t *testing.T) {
	c := v2.New(newClient(t, httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))))
	_, err := c.Transactions.Bulk(context.Background(), v2.BulkOptions{})
	if err == nil || !strings.Contains(err.Error(), "AccountIDs") {
		t.Errorf("err = %v", err)
	}
}

func TestTransactions_Bulk_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/transactions/bulk" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"bulk_results":[{"account_id":"a-1","transactions":[],"pagination":{"has_more":false,"per_page":50,"after_id":"","before_id":""}}]}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	out, err := c.Transactions.Bulk(context.Background(), v2.BulkOptions{AccountIDs: []string{"a-1"}})
	if err != nil {
		t.Fatalf("Bulk: %v", err)
	}
	if len(out.BulkResults) != 1 || out.BulkResults[0].AccountID != "a-1" {
		t.Errorf("bulk = %+v", out)
	}
}

func TestTransactions_Search_RequiresQ(t *testing.T) {
	c := v2.New(newClient(t, httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))))
	_, err := c.Transactions.Search(context.Background(), v2.SearchOptions{})
	if err == nil || !strings.Contains(err.Error(), "Q is required") {
		t.Errorf("err = %v", err)
	}
}

func TestTransactions_Search_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "starbucks" {
			t.Errorf("q = %q", r.URL.Query().Get("q"))
		}
		_, _ = io.WriteString(w, `{"transactions":[],"total":0}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.Transactions.Search(context.Background(), v2.SearchOptions{Q: "starbucks", Limit: 10})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
}

func TestTransactions_Export_CSV(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("format") != "csv" {
			t.Errorf("format = %q", r.URL.Query().Get("format"))
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", `attachment; filename="transactions_a-1_2026-04-28.csv"`)
		_, _ = io.WriteString(w, "Transaction ID,Date\nx,y\n")
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	res, err := c.Transactions.Export(context.Background(), "a-1", v2.ExportOptions{Format: v2.ExportFormatCSV})
	if err != nil {
		t.Fatalf("Export: %v", err)
	}
	if res.ContentType != "text/csv" {
		t.Errorf("content-type = %q", res.ContentType)
	}
	if res.Filename != "transactions_a-1_2026-04-28.csv" {
		t.Errorf("filename = %q", res.Filename)
	}
	if !strings.Contains(string(res.Body), "Transaction ID") {
		t.Errorf("body = %q", string(res.Body))
	}
}

func TestTransactions_Sync_415RequiresContentType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			_, _ = io.WriteString(w, `{"error":"need json","error_code":"UNPROCESSABLE_CONTENT"}`)
			return
		}
		_, _ = io.WriteString(w, `{"added":[],"modified":[],"removed":[],"next_cursor":null,"has_more":false}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	// Sync always sends a body, so Content-Type should be set automatically.
	_, err := c.Transactions.Sync(context.Background(), "a-1", v2.SyncOptions{Count: 10})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
}
