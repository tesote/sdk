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

func TestBatches_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Idempotency-Key") == "" {
			t.Error("missing Idempotency-Key on POST")
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{"batch_id":"b1","orders":[],"errors":[]}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	out, err := c.Batches.Create(context.Background(), "a1", v2.BatchCreateOptions{
		Orders: []v2.BatchOrderInput{{Amount: "10", Currency: "VES", Description: "x"}},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if out.BatchID != "b1" {
		t.Errorf("batch = %+v", out)
	}
}

func TestBatches_Show_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error":"x","error_code":"BATCH_NOT_FOUND"}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.Batches.Show(context.Background(), "a1", "missing")
	if !errors.Is(err, tesote.ErrBatchNotFound) {
		t.Errorf("err = %v", err)
	}
}

func TestBatches_Approve(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"approved":3,"failed":0}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	out, err := c.Batches.Approve(context.Background(), "a1", "b1", v2.BatchActionOptions{})
	if err != nil {
		t.Fatalf("Approve: %v", err)
	}
	if out.Approved != 3 {
		t.Errorf("approved = %d", out.Approved)
	}
}

func TestBatches_Submit_BatchValidationError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":"bad","error_code":"BATCH_VALIDATION_ERROR"}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.Batches.Submit(context.Background(), "a1", "b1", v2.BatchSubmitOptions{Token: "tok"})
	var typed *tesote.BatchValidationError
	if !errors.As(err, &typed) {
		t.Fatalf("err = %T %v", err, err)
	}
}

func TestBatches_Cancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"cancelled":2,"skipped":1,"errors":[]}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	out, err := c.Batches.Cancel(context.Background(), "a1", "b1", v2.BatchActionOptions{})
	if err != nil {
		t.Fatalf("Cancel: %v", err)
	}
	if out.Cancelled != 2 || out.Skipped != 1 {
		t.Errorf("response = %+v", out)
	}
}
