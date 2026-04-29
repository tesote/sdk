package v1_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	tesote "github.com/tesote/sdk/go"
	"github.com/tesote/sdk/go/v1"
)

// newClient builds a test transport pinned to the given httptest.Server.
func newClient(t *testing.T, srv *httptest.Server) *tesote.Client {
	t.Helper()
	c, err := tesote.NewClient(tesote.Options{
		APIKey:     "secret-key",
		BaseURL:    srv.URL,
		Sleep:      func(_ context.Context, _ time.Duration) error { return nil },
		RandUint63: func() uint64 { return 0 },
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestStatus_Status(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/status" {
			t.Errorf("path = %q, want /status", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"status":"ok","authenticated":false}`)
	}))
	defer srv.Close()

	c := v1.New(newClient(t, srv))
	out, err := c.Status.Status(context.Background())
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if out.Status != "ok" {
		t.Errorf("status = %q", out.Status)
	}
}

func TestStatus_Whoami(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/whoami" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"client":{"id":"abc","name":"acme","type":"workspace"}}`)
	}))
	defer srv.Close()

	c := v1.New(newClient(t, srv))
	out, err := c.Status.Whoami(context.Background())
	if err != nil {
		t.Fatalf("Whoami: %v", err)
	}
	if out.Client.ID != "abc" || out.Client.Type != "workspace" {
		t.Errorf("client = %+v", out.Client)
	}
}

func TestAccounts_List_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/accounts" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("per_page") != "10" {
			t.Errorf("query = %q", r.URL.RawQuery)
		}
		_, _ = io.WriteString(w, `{"total":1,"accounts":[{"id":"a-1","name":"checking","data":{"masked_account_number":"1234","currency":"VES"},"bank":{"name":"BoA"},"legal_entity":{"id":null,"legal_name":null},"tesote_created_at":"2026-04-01T00:00:00Z","tesote_updated_at":"2026-04-01T00:00:00Z"}],"pagination":{"current_page":2,"per_page":10,"total_pages":3,"total_count":21}}`)
	}))
	defer srv.Close()

	c := v1.New(newClient(t, srv))
	out, err := c.Accounts.List(context.Background(), v1.AccountsListOptions{Page: 2, PerPage: 10})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if out.Total != 1 || len(out.Accounts) != 1 {
		t.Errorf("response = %+v", out)
	}
	if out.Accounts[0].ID != "a-1" || out.Accounts[0].Bank.Name != "BoA" {
		t.Errorf("account = %+v", out.Accounts[0])
	}
}

func TestAccounts_Get_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error":"missing","error_code":"ACCOUNT_NOT_FOUND"}`)
	}))
	defer srv.Close()

	c := v1.New(newClient(t, srv))
	_, err := c.Accounts.Get(context.Background(), "nope")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, tesote.ErrAccountNotFound) {
		t.Errorf("not ErrAccountNotFound: %v", err)
	}
	var typed *tesote.AccountNotFoundError
	if !errors.As(err, &typed) {
		t.Errorf("not *AccountNotFoundError: %T", err)
	}
}

func TestTransactions_ListForAccount_CursorPagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/accounts/acct-1/transactions" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("transactions_after_id") != "tx-9" {
			t.Errorf("missing cursor")
		}
		_, _ = io.WriteString(w, `{"total":2,"transactions":[],"pagination":{"has_more":true,"per_page":50,"after_id":"tx-99","before_id":"tx-50"}}`)
	}))
	defer srv.Close()

	c := v1.New(newClient(t, srv))
	out, err := c.Transactions.ListForAccount(context.Background(), "acct-1", v1.TransactionsListOptions{TransactionsAfterID: "tx-9"})
	if err != nil {
		t.Fatalf("ListForAccount: %v", err)
	}
	if !out.Pagination.HasMore || out.Pagination.AfterID != "tx-99" {
		t.Errorf("pagination = %+v", out.Pagination)
	}
}

func TestTransactions_Get_DateRangeError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = io.WriteString(w, `{"error":"bad","error_code":"INVALID_DATE_RANGE"}`)
	}))
	defer srv.Close()

	c := v1.New(newClient(t, srv))
	_, err := c.Transactions.Get(context.Background(), "tx-1")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, tesote.ErrInvalidDateRange) {
		t.Errorf("not ErrInvalidDateRange: %v", err)
	}
}

func TestAccounts_AuthRedaction(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error":"nope","error_code":"UNAUTHORIZED"}`)
	}))
	defer srv.Close()

	c := v1.New(newClient(t, srv))
	_, err := c.Accounts.List(context.Background(), v1.AccountsListOptions{})
	var apiErr *tesote.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("not *APIError: %T", err)
	}
	if !strings.HasPrefix(apiErr.RequestSummary.Authorization, "Bearer ") {
		t.Errorf("redacted auth = %q", apiErr.RequestSummary.Authorization)
	}
	if strings.Contains(apiErr.RequestSummary.Authorization, "secret-") {
		t.Errorf("auth leaked: %q", apiErr.RequestSummary.Authorization)
	}
}
