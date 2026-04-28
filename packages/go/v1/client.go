// Package v1 implements the v1 API surface (paths under /api/v1, read-only).
package v1

import (
	"context"
	"errors"

	tesote "github.com/tesote/sdk/go"
)

// pathPrefix is the URL prefix every v1 request uses.
const pathPrefix = "/v1"

// V1Client groups v1 resource services.
type V1Client struct {
	transport *tesote.Client

	Accounts     *AccountsService
	Transactions *TransactionsService
	Status       *StatusService
}

// New builds a V1Client around an existing tesote.Client.
func New(t *tesote.Client) *V1Client {
	return &V1Client{
		transport:    t,
		Accounts:     &AccountsService{client: t},
		Transactions: &TransactionsService{},
		Status:       &StatusService{},
	}
}

// Transport exposes the underlying *tesote.Client.
func (c *V1Client) Transport() *tesote.Client { return c.transport }

// ErrNotImplemented is returned by every stub method.
var ErrNotImplemented = errors.New("tesote/v1: not implemented")

// AccountsService groups v1 account endpoints.
type AccountsService struct {
	client *tesote.Client
}

// List lists accounts. Stub.
func (s *AccountsService) List(_ context.Context) error { return ErrNotImplemented }

// Get fetches an account. Stub.
func (s *AccountsService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// TransactionsService groups v1 transaction endpoints.
type TransactionsService struct{}

// ListForAccount lists transactions for an account. Stub.
func (TransactionsService) ListForAccount(_ context.Context, _ string) error {
	return ErrNotImplemented
}

// Get fetches a transaction. Stub.
func (TransactionsService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// StatusService groups v1 status endpoints.
type StatusService struct{}

// Status returns API status. Stub.
func (StatusService) Status(_ context.Context) error { return ErrNotImplemented }

// Whoami returns the API key's identity. Stub.
func (StatusService) Whoami(_ context.Context) error { return ErrNotImplemented }

// Suppress unused imports when only stubs exist.
var _ = pathPrefix
