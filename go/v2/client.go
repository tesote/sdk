// Package v2 implements the v2 API surface (paths under /api/v2).
package v2

import (
	"context"
	"errors"

	tesote "github.com/tesote/sdk/go"
)

// pathPrefix is the URL prefix every v2 request uses.
const pathPrefix = "/v2"

// V2Client groups v2 resource services.
type V2Client struct {
	transport *tesote.Client

	Accounts          *AccountsService
	Transactions      *TransactionsService
	SyncSessions      *SyncSessionsService
	TransactionOrders *TransactionOrdersService
	Batches           *BatchesService
	PaymentMethods    *PaymentMethodsService
	Status            *StatusService
}

// New builds a V2Client around an existing tesote.Client.
func New(t *tesote.Client) *V2Client {
	return &V2Client{
		transport:         t,
		Accounts:          &AccountsService{client: t},
		Transactions:      &TransactionsService{},
		SyncSessions:      &SyncSessionsService{},
		TransactionOrders: &TransactionOrdersService{},
		Batches:           &BatchesService{},
		PaymentMethods:    &PaymentMethodsService{},
		Status:            &StatusService{},
	}
}

// Transport exposes the underlying *tesote.Client.
func (c *V2Client) Transport() *tesote.Client { return c.transport }

// ErrNotImplemented is returned by every stub method.
var ErrNotImplemented = errors.New("tesote/v2: not implemented")

// AccountsService groups v2 account endpoints.
type AccountsService struct {
	client *tesote.Client
}

// List lists accounts. Stub.
func (s *AccountsService) List(_ context.Context) error { return ErrNotImplemented }

// Get fetches an account. Stub.
func (s *AccountsService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// Sync triggers an account sync. Stub.
func (s *AccountsService) Sync(_ context.Context, _ string) error { return ErrNotImplemented }

// TransactionsService groups v2 transaction endpoints.
type TransactionsService struct{}

// ListForAccount lists transactions for an account. Stub.
func (TransactionsService) ListForAccount(_ context.Context, _ string) error {
	return ErrNotImplemented
}

// Get fetches a transaction. Stub.
func (TransactionsService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// Export exports transactions. Stub.
func (TransactionsService) Export(_ context.Context) error { return ErrNotImplemented }

// Sync triggers a transaction sync. Stub.
func (TransactionsService) Sync(_ context.Context, _ string) error { return ErrNotImplemented }

// Bulk performs a bulk transaction operation. Stub.
func (TransactionsService) Bulk(_ context.Context) error { return ErrNotImplemented }

// Search searches transactions. Stub.
func (TransactionsService) Search(_ context.Context) error { return ErrNotImplemented }

// SyncSessionsService groups v2 sync-session endpoints.
type SyncSessionsService struct{}

// List lists sync sessions for an account. Stub.
func (SyncSessionsService) List(_ context.Context, _ string) error { return ErrNotImplemented }

// Get fetches a sync session. Stub.
func (SyncSessionsService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// TransactionOrdersService groups v2 transaction-order endpoints.
type TransactionOrdersService struct{}

// List lists transaction orders. Stub.
func (TransactionOrdersService) List(_ context.Context, _ string) error { return ErrNotImplemented }

// Get fetches a transaction order. Stub.
func (TransactionOrdersService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// Create creates a transaction order. Stub.
func (TransactionOrdersService) Create(_ context.Context, _ string) error {
	return ErrNotImplemented
}

// Submit submits a transaction order. Stub.
func (TransactionOrdersService) Submit(_ context.Context, _ string) error {
	return ErrNotImplemented
}

// Cancel cancels a transaction order. Stub.
func (TransactionOrdersService) Cancel(_ context.Context, _ string) error {
	return ErrNotImplemented
}

// BatchesService groups v2 batch endpoints.
type BatchesService struct{}

// Create creates a batch. Stub.
func (BatchesService) Create(_ context.Context) error { return ErrNotImplemented }

// Get fetches a batch. Stub.
func (BatchesService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// Approve approves a batch. Stub.
func (BatchesService) Approve(_ context.Context, _ string) error { return ErrNotImplemented }

// Submit submits a batch. Stub.
func (BatchesService) Submit(_ context.Context, _ string) error { return ErrNotImplemented }

// Cancel cancels a batch. Stub.
func (BatchesService) Cancel(_ context.Context, _ string) error { return ErrNotImplemented }

// PaymentMethodsService groups v2 payment-method endpoints.
type PaymentMethodsService struct{}

// List lists payment methods. Stub.
func (PaymentMethodsService) List(_ context.Context) error { return ErrNotImplemented }

// Get fetches a payment method. Stub.
func (PaymentMethodsService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// Create creates a payment method. Stub.
func (PaymentMethodsService) Create(_ context.Context) error { return ErrNotImplemented }

// Update updates a payment method. Stub.
func (PaymentMethodsService) Update(_ context.Context, _ string) error { return ErrNotImplemented }

// Delete deletes a payment method. Stub.
func (PaymentMethodsService) Delete(_ context.Context, _ string) error { return ErrNotImplemented }

// StatusService groups v2 status endpoints.
type StatusService struct{}

// Status returns API status. Stub.
func (StatusService) Status(_ context.Context) error { return ErrNotImplemented }

// Whoami returns the API key identity. Stub.
func (StatusService) Whoami(_ context.Context) error { return ErrNotImplemented }

// Suppress unused warning until v2 wires endpoints.
var _ = pathPrefix
