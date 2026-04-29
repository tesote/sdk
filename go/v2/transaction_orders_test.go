package v2_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	tesote "github.com/tesote/sdk/go"
	"github.com/tesote/sdk/go/v2"
)

const orderJSON = `{"id":"o1","status":"draft","amount":1000.00,"currency":"VES","description":"x","reference":null,"external_reference":null,"idempotency_key":null,"batch_id":null,"scheduled_for":null,"approved_at":null,"submitted_at":null,"completed_at":null,"failed_at":null,"cancelled_at":null,"source_account":{"id":"a1","name":"checking","payment_method_id":"p1"},"destination":{"payment_method_id":"p2","counterparty_id":"c1","counterparty_name":"acme"},"fee":null,"execution_strategy":null,"tesote_transaction":null,"latest_attempt":null,"created_at":"x","updated_at":"x"}`

func TestTransactionOrders_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") != "draft" {
			t.Errorf("query = %q", r.URL.RawQuery)
		}
		_, _ = io.WriteString(w, `{"items":[`+orderJSON+`],"has_more":false,"limit":50,"offset":0}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	out, err := c.TransactionOrders.List(context.Background(), "a1", v2.TransactionOrdersListOptions{Status: "draft"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(out.Items) != 1 || out.Items[0].ID != "o1" {
		t.Errorf("items = %+v", out)
	}
}

func TestTransactionOrders_Create_BodyShape(t *testing.T) {
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, orderJSON)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	idem := "k-1"
	_, err := c.TransactionOrders.Create(context.Background(), "a1", v2.TransactionOrderCreateOptions{
		Amount:         "100.00",
		Currency:       "VES",
		Description:    "rent",
		IdempotencyKey: &idem,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	tx, ok := body["transaction_order"].(map[string]any)
	if !ok || tx["amount"] != "100.00" || tx["currency"] != "VES" {
		t.Errorf("body = %+v", body)
	}
}

func TestTransactionOrders_Submit_InvalidOrderState(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = io.WriteString(w, `{"error":"bad state","error_code":"INVALID_ORDER_STATE"}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.TransactionOrders.Submit(context.Background(), "a1", "o1", v2.SubmitOrderOptions{})
	if !errors.Is(err, tesote.ErrInvalidOrderState) {
		t.Errorf("err = %v", err)
	}
}

func TestTransactionOrders_Cancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/accounts/a1/transaction_orders/o1/cancel" {
			t.Errorf("path = %q", r.URL.Path)
		}
		_, _ = io.WriteString(w, orderJSON)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.TransactionOrders.Cancel(context.Background(), "a1", "o1", v2.CancelOrderOptions{})
	if err != nil {
		t.Fatalf("Cancel: %v", err)
	}
}

func TestTransactionOrders_Get_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error":"x","error_code":"TRANSACTION_ORDER_NOT_FOUND"}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.TransactionOrders.Get(context.Background(), "a1", "missing")
	var typed *tesote.TransactionOrderNotFoundError
	if !errors.As(err, &typed) {
		t.Errorf("err = %T %v", err, err)
	}
}
