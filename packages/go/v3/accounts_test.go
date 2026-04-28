package v3

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
)

func newClientFor(t *testing.T, h http.Handler) (*V3Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(h)
	tc, err := tesote.NewClient(tesote.Options{
		APIKey:     "secret-abcd",
		BaseURL:    srv.URL,
		Sleep:      func(_ context.Context, _ time.Duration) error { return nil },
		RandUint63: func() uint64 { return 0 },
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return New(tc), srv
}

func TestAccounts_List(t *testing.T) {
	cases := []struct {
		name      string
		opts      ListOptions
		wantQuery string
	}{
		{"no opts", ListOptions{}, ""},
		{"cursor", ListOptions{Cursor: "abc"}, "cursor=abc"},
		{"page size", ListOptions{PageSize: 50}, "page_size=50"},
		{"both", ListOptions{Cursor: "x", PageSize: 25}, "cursor=x&page_size=25"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var seenPath, seenQuery, seenAuth string
			c, srv := newClientFor(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				seenPath = r.URL.Path
				seenQuery = r.URL.RawQuery
				seenAuth = r.Header.Get("Authorization")
				_, _ = io.WriteString(w, `{"data":[{"id":"acct_1","name":"checking","currency":"USD"}],"next_cursor":"nx"}`)
			}))
			defer srv.Close()

			out, err := c.Accounts.List(context.Background(), tc.opts)
			if err != nil {
				t.Fatalf("List: %v", err)
			}
			if seenPath != "/v3/accounts" {
				t.Errorf("path = %q", seenPath)
			}
			if seenQuery != tc.wantQuery {
				t.Errorf("query = %q, want %q", seenQuery, tc.wantQuery)
			}
			if seenAuth != "Bearer secret-abcd" {
				t.Errorf("auth = %q", seenAuth)
			}
			if len(out.Data) != 1 || out.Data[0].ID != "acct_1" {
				t.Errorf("decode = %+v", out)
			}
			if out.NextCursor != "nx" {
				t.Errorf("cursor = %q", out.NextCursor)
			}
		})
	}
}

func TestAccounts_Get(t *testing.T) {
	c, srv := newClientFor(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/accounts/acct_42" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"id":"acct_42","name":"savings","currency":"EUR"}`)
	}))
	defer srv.Close()

	got, err := c.Accounts.Get(context.Background(), "acct_42")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != "acct_42" || got.Name != "savings" || got.Currency != "EUR" {
		t.Errorf("decode = %+v", got)
	}
}

func TestAccounts_Get_EmptyIDIsConfigError(t *testing.T) {
	c, srv := newClientFor(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected request to %s", r.URL.Path)
	}))
	defer srv.Close()

	_, err := c.Accounts.Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, tesote.ErrConfig) {
		t.Errorf("expected ErrConfig, got %v", err)
	}
}

func TestAccounts_List_PropagatesAPIError(t *testing.T) {
	c, srv := newClientFor(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", "req-acct-err")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error":"nope","error_code":"UNAUTHORIZED"}`)
	}))
	defer srv.Close()

	_, err := c.Accounts.List(context.Background(), ListOptions{})
	if err == nil {
		t.Fatal("expected error")
	}
	var u *tesote.UnauthorizedError
	if !errors.As(err, &u) {
		t.Fatalf("expected *UnauthorizedError, got %T", err)
	}
	if u.RequestID != "req-acct-err" {
		t.Errorf("request id = %q", u.RequestID)
	}
	if !strings.Contains(u.RequestSummary.Path, "/v3/accounts") {
		t.Errorf("summary path = %q", u.RequestSummary.Path)
	}
}

func TestAccounts_Sync_Stub(t *testing.T) {
	tc, _ := tesote.NewClient(tesote.Options{APIKey: "x"})
	c := New(tc)
	if err := c.Accounts.Sync(context.Background(), "acct_1"); err == nil {
		t.Error("expected stub error")
	}
}
