package v2_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	tesote "github.com/tesote/sdk/go"
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
