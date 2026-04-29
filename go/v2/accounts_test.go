package v2_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	tesote "github.com/tesote/sdk/go"
	"github.com/tesote/sdk/go/v2"
)

func TestAccounts_List_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/accounts" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"total":0,"accounts":[],"pagination":{"current_page":1,"per_page":50,"total_pages":0,"total_count":0}}`)
	}))
	defer srv.Close()

	c := v2.New(newClient(t, srv))
	_, err := c.Accounts.List(context.Background(), v2.AccountsListOptions{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
}

func TestAccounts_Get_404Typed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error":"x","error_code":"ACCOUNT_NOT_FOUND"}`)
	}))
	defer srv.Close()

	c := v2.New(newClient(t, srv))
	_, err := c.Accounts.Get(context.Background(), "nope")
	if !errors.Is(err, tesote.ErrAccountNotFound) {
		t.Errorf("err = %v", err)
	}
}

func TestAccounts_Sync_Success(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %s", r.Method)
			}
			if r.URL.Path != "/v2/accounts/abc/sync" {
				t.Errorf("path = %q", r.URL.Path)
			}
			if r.Header.Get("Idempotency-Key") == "" {
				t.Error("missing Idempotency-Key")
			}
			w.WriteHeader(http.StatusAccepted)
			_, _ = io.WriteString(w, `{"message":"Sync started","sync_session_id":"sid","status":"pending","started_at":"2026-04-28T19:21:00Z"}`)
		}))
		defer srv.Close()

		c := v2.New(newClient(t, srv))
		out, err := c.Accounts.Sync(context.Background(), "abc", v2.AccountSyncOptions{})
		if err != nil {
			t.Fatalf("Sync: %v", err)
		}
		if out.SyncSessionID != "sid" || out.Status != "pending" {
			t.Errorf("response = %+v", out)
		}
	})

	t.Run("sync in progress -> typed", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
			_, _ = io.WriteString(w, `{"error":"busy","error_code":"SYNC_IN_PROGRESS"}`)
		}))
		defer srv.Close()
		c := v2.New(newClient(t, srv))
		_, err := c.Accounts.Sync(context.Background(), "abc", v2.AccountSyncOptions{})
		if !errors.Is(err, tesote.ErrSyncInProgress) {
			t.Errorf("err = %v", err)
		}
		var typed *tesote.SyncInProgressError
		if !errors.As(err, &typed) {
			t.Errorf("not *SyncInProgressError: %T", err)
		}
	})

	t.Run("sync rate limit -> typed", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Retry-After", "300")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = io.WriteString(w, `{"error":"slow","error_code":"SYNC_RATE_LIMIT_EXCEEDED","retry_after":300}`)
		}))
		defer srv.Close()
		c := v2.New(newClient(t, srv))
		_, err := c.Accounts.Sync(context.Background(), "abc", v2.AccountSyncOptions{})
		var typed *tesote.SyncRateLimitExceededError
		if !errors.As(err, &typed) {
			t.Fatalf("not *SyncRateLimitExceededError: %T %v", err, err)
		}
		if typed.RetryAfter != 300 {
			t.Errorf("retry_after = %d", typed.RetryAfter)
		}
	})

	t.Run("bank under maintenance -> typed", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = io.WriteString(w, `{"error":"down","error_code":"BANK_UNDER_MAINTENANCE"}`)
		}))
		defer srv.Close()
		c := v2.New(newClient(t, srv))
		_, err := c.Accounts.Sync(context.Background(), "abc", v2.AccountSyncOptions{})
		var typed *tesote.BankUnderMaintenanceError
		if !errors.As(err, &typed) {
			t.Errorf("err = %T %v", err, err)
		}
	})
}

func TestAccounts_Sync_PreservesUserIdempotencyKey(t *testing.T) {
	var seen string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.Header.Get("Idempotency-Key")
		w.WriteHeader(http.StatusAccepted)
		_, _ = io.WriteString(w, `{"message":"x","sync_session_id":"y","status":"pending","started_at":"z"}`)
	}))
	defer srv.Close()

	c := v2.New(newClient(t, srv))
	_, err := c.Accounts.Sync(context.Background(), "abc", v2.AccountSyncOptions{IdempotencyKey: "user-key"})
	if err != nil {
		t.Fatalf("Sync: %v", err)
	}
	if seen != "user-key" {
		t.Errorf("idempotency-key = %q", seen)
	}
}
