// Package v2 implements the v2 API surface (paths under /api/v2).
package v2

import (
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
		Transactions:      &TransactionsService{client: t},
		SyncSessions:      &SyncSessionsService{client: t},
		TransactionOrders: &TransactionOrdersService{client: t},
		Batches:           &BatchesService{client: t},
		PaymentMethods:    &PaymentMethodsService{client: t},
		Status:            &StatusService{client: t},
	}
}

// Transport exposes the underlying *tesote.Client.
func (c *V2Client) Transport() *tesote.Client { return c.transport }
