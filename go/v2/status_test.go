package v2_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tesote/sdk/go/v2"
)

func TestV2Status_Status(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/status" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"status":"ok","authenticated":false}`)
	}))
	defer srv.Close()

	c := v2.New(newClient(t, srv))
	out, err := c.Status.Status(context.Background())
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if out.Status != "ok" {
		t.Errorf("status = %q", out.Status)
	}
}

func TestV2Status_Whoami(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/whoami" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = io.WriteString(w, `{"client":{"id":"u","name":"n","type":"workspace"}}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	out, err := c.Status.Whoami(context.Background())
	if err != nil {
		t.Fatalf("Whoami: %v", err)
	}
	if out.Client.Type != "workspace" {
		t.Errorf("client = %+v", out.Client)
	}
}
