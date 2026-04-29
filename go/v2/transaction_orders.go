package v2

import (
	"context"
	"strconv"

	tesote "github.com/tesote/sdk/go"
)

// TransactionOrdersService groups v2 transaction-order endpoints.
type TransactionOrdersService struct {
	client *tesote.Client
}

// TransactionOrdersListOptions tunes GET /v2/accounts/{id}/transaction_orders.
type TransactionOrdersListOptions struct {
	Limit         int
	Offset        int
	Status        string
	CreatedAfter  string
	CreatedBefore string
	BatchID       string
}

func (o TransactionOrdersListOptions) query() map[string]string {
	q := map[string]string{}
	if o.Limit > 0 {
		q["limit"] = strconv.Itoa(o.Limit)
	}
	if o.Offset > 0 {
		q["offset"] = strconv.Itoa(o.Offset)
	}
	if o.Status != "" {
		q["status"] = o.Status
	}
	if o.CreatedAfter != "" {
		q["created_after"] = o.CreatedAfter
	}
	if o.CreatedBefore != "" {
		q["created_before"] = o.CreatedBefore
	}
	if o.BatchID != "" {
		q["batch_id"] = o.BatchID
	}
	return q
}

// List lists transaction orders for an account.
// GET /v2/accounts/{id}/transaction_orders.
func (s *TransactionOrdersService) List(ctx context.Context, accountID string, opts TransactionOrdersListOptions) (*tesote.TransactionOrderListResponse, error) {
	out := &tesote.TransactionOrderListResponse{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts/"+accountID+"/transaction_orders", tesote.RequestOptions{
		Query: opts.query(),
		Out:   out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single transaction order.
// GET /v2/accounts/{id}/transaction_orders/{order_id}.
func (s *TransactionOrdersService) Get(ctx context.Context, accountID, orderID string) (*tesote.TransactionOrder, error) {
	out := &tesote.TransactionOrder{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts/"+accountID+"/transaction_orders/"+orderID, tesote.RequestOptions{
		Out: out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Beneficiary is the on-the-fly beneficiary block accepted by Create.
type Beneficiary struct {
	Name                 string  `json:"name"`
	BankCode             *string `json:"bank_code,omitempty"`
	AccountNumber        *string `json:"account_number,omitempty"`
	IdentificationType   *string `json:"identification_type,omitempty"`
	IdentificationNumber *string `json:"identification_number,omitempty"`
}

// TransactionOrderCreateOptions is the body for POST .../transaction_orders.
type TransactionOrderCreateOptions struct {
	DestinationPaymentMethodID *string        `json:"destination_payment_method_id,omitempty"`
	Beneficiary                *Beneficiary   `json:"beneficiary,omitempty"`
	Amount                     string         `json:"amount"`
	Currency                   string         `json:"currency"`
	Description                string         `json:"description"`
	ScheduledFor               *string        `json:"scheduled_for,omitempty"`
	IdempotencyKey             *string        `json:"idempotency_key,omitempty"`
	Metadata                   map[string]any `json:"metadata,omitempty"`
}

// Create creates a draft transaction order.
// POST /v2/accounts/{id}/transaction_orders.
func (s *TransactionOrdersService) Create(ctx context.Context, accountID string, opts TransactionOrderCreateOptions) (*tesote.TransactionOrder, error) {
	body := map[string]any{"transaction_order": opts}
	idem := ""
	if opts.IdempotencyKey != nil {
		idem = *opts.IdempotencyKey
	}
	out := &tesote.TransactionOrder{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/accounts/"+accountID+"/transaction_orders", tesote.RequestOptions{
		Body:           body,
		IdempotencyKey: idem,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubmitOrderOptions is the body for POST .../transaction_orders/{id}/submit.
type SubmitOrderOptions struct {
	Token          string
	IdempotencyKey string
}

// Submit submits a transaction order for bank execution.
// POST /v2/accounts/{id}/transaction_orders/{order_id}/submit.
func (s *TransactionOrdersService) Submit(ctx context.Context, accountID, orderID string, opts SubmitOrderOptions) (*tesote.TransactionOrder, error) {
	body := map[string]any{}
	if opts.Token != "" {
		body["token"] = opts.Token
	}
	out := &tesote.TransactionOrder{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/accounts/"+accountID+"/transaction_orders/"+orderID+"/submit", tesote.RequestOptions{
		Body:           body,
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CancelOrderOptions wraps idempotency tuning for Cancel.
type CancelOrderOptions struct {
	IdempotencyKey string
}

// Cancel cancels a transaction order.
// POST /v2/accounts/{id}/transaction_orders/{order_id}/cancel.
func (s *TransactionOrdersService) Cancel(ctx context.Context, accountID, orderID string, opts CancelOrderOptions) (*tesote.TransactionOrder, error) {
	out := &tesote.TransactionOrder{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/accounts/"+accountID+"/transaction_orders/"+orderID+"/cancel", tesote.RequestOptions{
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
