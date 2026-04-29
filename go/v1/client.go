// Package v1 implements the v1 API surface (paths under /api/v1, read-only).
package v1

import (
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
		Transactions: &TransactionsService{client: t},
		Status:       &StatusService{client: t},
	}
}

// Transport exposes the underlying *tesote.Client.
func (c *V1Client) Transport() *tesote.Client { return c.transport }
