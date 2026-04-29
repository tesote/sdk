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

func TestSyncSessions_List_Pagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("limit") != "10" || q.Get("offset") != "20" || q.Get("status") != "completed" {
			t.Errorf("query = %q", r.URL.RawQuery)
		}
		_, _ = io.WriteString(w, `{"sync_sessions":[{"id":"s1","status":"completed","started_at":"x","completed_at":null,"transactions_synced":0,"accounts_count":1,"error":null,"performance":null}],"limit":10,"offset":20,"has_more":true}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	out, err := c.SyncSessions.List(context.Background(), "a-1", v2.SyncSessionsListOptions{Limit: 10, Offset: 20, Status: "completed"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if !out.HasMore || len(out.SyncSessions) != 1 {
		t.Errorf("response = %+v", out)
	}
}

func TestSyncSessions_Get_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error":"missing","error_code":"SYNC_SESSION_NOT_FOUND"}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.SyncSessions.Get(context.Background(), "a-1", "nope")
	if !errors.Is(err, tesote.ErrSyncSessionNotFound) {
		t.Errorf("err = %v", err)
	}
}
