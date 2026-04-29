package v2

import (
	"context"
	"strconv"

	tesote "github.com/tesote/sdk/go"
)

// PaymentMethodsService groups v2 payment-method endpoints.
type PaymentMethodsService struct {
	client *tesote.Client
}

// PaymentMethodsListOptions tunes GET /v2/payment_methods.
type PaymentMethodsListOptions struct {
	Limit          int
	Offset         int
	MethodType     string
	Currency       string
	CounterpartyID string
	Verified       *bool
}

func (o PaymentMethodsListOptions) query() map[string]string {
	q := map[string]string{}
	if o.Limit > 0 {
		q["limit"] = strconv.Itoa(o.Limit)
	}
	if o.Offset > 0 {
		q["offset"] = strconv.Itoa(o.Offset)
	}
	if o.MethodType != "" {
		q["method_type"] = o.MethodType
	}
	if o.Currency != "" {
		q["currency"] = o.Currency
	}
	if o.CounterpartyID != "" {
		q["counterparty_id"] = o.CounterpartyID
	}
	if o.Verified != nil {
		if *o.Verified {
			q["verified"] = "true"
		} else {
			q["verified"] = "false"
		}
	}
	return q
}

// List lists payment methods. GET /v2/payment_methods.
func (s *PaymentMethodsService) List(ctx context.Context, opts PaymentMethodsListOptions) (*tesote.PaymentMethodListResponse, error) {
	out := &tesote.PaymentMethodListResponse{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/payment_methods", tesote.RequestOptions{
		Query: opts.query(),
		Out:   out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single payment method. GET /v2/payment_methods/{id}.
func (s *PaymentMethodsService) Get(ctx context.Context, id string) (*tesote.PaymentMethod, error) {
	out := &tesote.PaymentMethod{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/payment_methods/"+id, tesote.RequestOptions{
		Out: out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentMethodCounterpartyInput is the inline counterparty block on create.
type PaymentMethodCounterpartyInput struct {
	Name string `json:"name"`
}

// PaymentMethodInput is the body wrapped under "payment_method" for create/update.
type PaymentMethodInput struct {
	MethodType     string                          `json:"method_type,omitempty"`
	Currency       string                          `json:"currency,omitempty"`
	Label          *string                         `json:"label,omitempty"`
	CounterpartyID *string                         `json:"counterparty_id,omitempty"`
	Counterparty   *PaymentMethodCounterpartyInput `json:"counterparty,omitempty"`
	Details        map[string]any                  `json:"details,omitempty"`
}

// PaymentMethodMutateOptions wraps the body with idempotency control.
type PaymentMethodMutateOptions struct {
	Input          PaymentMethodInput
	IdempotencyKey string
}

// Create creates a payment method. POST /v2/payment_methods.
func (s *PaymentMethodsService) Create(ctx context.Context, opts PaymentMethodMutateOptions) (*tesote.PaymentMethod, error) {
	body := map[string]any{"payment_method": opts.Input}
	out := &tesote.PaymentMethod{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/payment_methods", tesote.RequestOptions{
		Body:           body,
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Update partially updates a payment method. PATCH /v2/payment_methods/{id}.
func (s *PaymentMethodsService) Update(ctx context.Context, id string, opts PaymentMethodMutateOptions) (*tesote.PaymentMethod, error) {
	body := map[string]any{"payment_method": opts.Input}
	out := &tesote.PaymentMethod{}
	_, err := s.client.Do(ctx, "PATCH", pathPrefix+"/payment_methods/"+id, tesote.RequestOptions{
		Body:           body,
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentMethodDeleteOptions tunes idempotency for Delete.
type PaymentMethodDeleteOptions struct {
	IdempotencyKey string
}

// Delete soft-deletes a payment method. DELETE /v2/payment_methods/{id}.
func (s *PaymentMethodsService) Delete(ctx context.Context, id string, opts PaymentMethodDeleteOptions) error {
	_, err := s.client.Do(ctx, "DELETE", pathPrefix+"/payment_methods/"+id, tesote.RequestOptions{
		IdempotencyKey: opts.IdempotencyKey,
	})
	return err
}
