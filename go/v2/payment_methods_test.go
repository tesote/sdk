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

const pmJSON = `{"id":"p1","method_type":"bank_account","currency":"VES","label":null,"details":{},"verified":false,"verified_at":null,"last_used_at":null,"counterparty":null,"tesote_account":null,"created_at":"x","updated_at":"x"}`

func TestPaymentMethods_List_VerifiedFilter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("verified") != "true" {
			t.Errorf("verified = %q", r.URL.Query().Get("verified"))
		}
		_, _ = io.WriteString(w, `{"items":[`+pmJSON+`],"has_more":false,"limit":50,"offset":0}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	verified := true
	out, err := c.PaymentMethods.List(context.Background(), v2.PaymentMethodsListOptions{Verified: &verified})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(out.Items) != 1 {
		t.Errorf("items = %+v", out)
	}
}

func TestPaymentMethods_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, pmJSON)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	out, err := c.PaymentMethods.Get(context.Background(), "p1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if out.ID != "p1" || out.MethodType != "bank_account" {
		t.Errorf("pm = %+v", out)
	}
}

func TestPaymentMethods_Get_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error":"x","error_code":"PAYMENT_METHOD_NOT_FOUND"}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.PaymentMethods.Get(context.Background(), "missing")
	if !errors.Is(err, tesote.ErrPaymentMethodNotFound) {
		t.Errorf("err = %v", err)
	}
}

func TestPaymentMethods_Create_ValidationError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":"bad","error_code":"VALIDATION_ERROR"}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	_, err := c.PaymentMethods.Create(context.Background(), v2.PaymentMethodMutateOptions{
		Input: v2.PaymentMethodInput{MethodType: "bank_account"},
	})
	var typed *tesote.ValidationError
	if !errors.As(err, &typed) {
		t.Errorf("err = %T %v", err, err)
	}
}

func TestPaymentMethods_Update_PATCH(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %s", r.Method)
		}
		_, _ = io.WriteString(w, pmJSON)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	label := "primary"
	_, err := c.PaymentMethods.Update(context.Background(), "p1", v2.PaymentMethodMutateOptions{
		Input: v2.PaymentMethodInput{Label: &label},
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
}

func TestPaymentMethods_Delete_204(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	if err := c.PaymentMethods.Delete(context.Background(), "p1", v2.PaymentMethodDeleteOptions{}); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestPaymentMethods_Delete_InUseConflict(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = io.WriteString(w, `{"error":"in use","error_code":"VALIDATION_ERROR"}`)
	}))
	defer srv.Close()
	c := v2.New(newClient(t, srv))
	err := c.PaymentMethods.Delete(context.Background(), "p1", v2.PaymentMethodDeleteOptions{})
	var typed *tesote.ValidationError
	if !errors.As(err, &typed) {
		t.Errorf("err = %T %v", err, err)
	}
}
