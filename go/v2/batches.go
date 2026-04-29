package v2

import (
	"context"

	tesote "github.com/tesote/sdk/go"
)

// BatchesService groups v2 batch endpoints.
type BatchesService struct {
	client *tesote.Client
}

// BatchOrderInput is one entry in the orders array of a batch create.
type BatchOrderInput struct {
	DestinationPaymentMethodID *string        `json:"destination_payment_method_id,omitempty"`
	Beneficiary                *Beneficiary   `json:"beneficiary,omitempty"`
	Amount                     string         `json:"amount"`
	Currency                   string         `json:"currency"`
	Description                string         `json:"description"`
	ScheduledFor               *string        `json:"scheduled_for,omitempty"`
	Metadata                   map[string]any `json:"metadata,omitempty"`
}

// BatchCreateOptions is the body for POST /v2/accounts/{id}/batches.
type BatchCreateOptions struct {
	Orders         []BatchOrderInput
	IdempotencyKey string
}

// Create creates a batch of transaction orders.
// POST /v2/accounts/{id}/batches.
func (s *BatchesService) Create(ctx context.Context, accountID string, opts BatchCreateOptions) (*tesote.BatchCreateResponse, error) {
	body := map[string]any{"orders": opts.Orders}
	out := &tesote.BatchCreateResponse{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/accounts/"+accountID+"/batches", tesote.RequestOptions{
		Body:           body,
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Show fetches a batch summary.
// GET /v2/accounts/{id}/batches/{batch_id}.
func (s *BatchesService) Show(ctx context.Context, accountID, batchID string) (*tesote.BatchSummary, error) {
	out := &tesote.BatchSummary{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts/"+accountID+"/batches/"+batchID, tesote.RequestOptions{
		Out: out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BatchActionOptions tunes idempotency for batch state transitions.
type BatchActionOptions struct {
	IdempotencyKey string
}

// Approve transitions all draft orders in the batch to pending_approval.
// POST /v2/accounts/{id}/batches/{batch_id}/approve.
func (s *BatchesService) Approve(ctx context.Context, accountID, batchID string, opts BatchActionOptions) (*tesote.BatchApproveResponse, error) {
	out := &tesote.BatchApproveResponse{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/accounts/"+accountID+"/batches/"+batchID+"/approve", tesote.RequestOptions{
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BatchSubmitOptions is the body for POST .../batches/{id}/submit.
type BatchSubmitOptions struct {
	Token          string
	IdempotencyKey string
}

// Submit enqueues all orders in the batch for bank submission.
// POST /v2/accounts/{id}/batches/{batch_id}/submit.
func (s *BatchesService) Submit(ctx context.Context, accountID, batchID string, opts BatchSubmitOptions) (*tesote.BatchSubmitResponse, error) {
	body := map[string]any{}
	if opts.Token != "" {
		body["token"] = opts.Token
	}
	out := &tesote.BatchSubmitResponse{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/accounts/"+accountID+"/batches/"+batchID+"/submit", tesote.RequestOptions{
		Body:           body,
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Cancel cancels all eligible orders in the batch.
// POST /v2/accounts/{id}/batches/{batch_id}/cancel.
func (s *BatchesService) Cancel(ctx context.Context, accountID, batchID string, opts BatchActionOptions) (*tesote.BatchCancelResponse, error) {
	out := &tesote.BatchCancelResponse{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/accounts/"+accountID+"/batches/"+batchID+"/cancel", tesote.RequestOptions{
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
