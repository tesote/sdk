package v3

import (
	"context"
	"errors"
)

// Per resources.md, v3 ships every v2 resource plus the v3-specific ones below.
// Each service mirrors the canonical method list. Stubs return ErrNotImplemented
// until they are wired end-to-end.
//
// Removing or renaming any exported symbol here is a breaking change.

// ErrNotImplemented is returned by every stub method.
var ErrNotImplemented = errors.New("tesote/v3: not implemented")

// TransactionsService groups transaction endpoints.
type TransactionsService struct{}

// ListForAccount lists transactions for an account. Stub.
func (TransactionsService) ListForAccount(_ context.Context, _ string) error {
	return ErrNotImplemented
}

// Get fetches a transaction. Stub.
func (TransactionsService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// Export exports transactions. Stub.
func (TransactionsService) Export(_ context.Context) error { return ErrNotImplemented }

// Sync triggers a transactions sync. Stub.
func (TransactionsService) Sync(_ context.Context, _ string) error { return ErrNotImplemented }

// Bulk performs a bulk transaction operation. Stub.
func (TransactionsService) Bulk(_ context.Context) error { return ErrNotImplemented }

// Search runs a transactions search. Stub.
func (TransactionsService) Search(_ context.Context) error { return ErrNotImplemented }

// SyncSessionsService groups sync-session endpoints.
type SyncSessionsService struct{}

// List lists sync sessions for an account. Stub.
func (SyncSessionsService) List(_ context.Context, _ string) error { return ErrNotImplemented }

// Get fetches a sync session. Stub.
func (SyncSessionsService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// TransactionOrdersService groups transaction-order endpoints.
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

// BatchesService groups batch endpoints.
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

// PaymentMethodsService groups payment-method endpoints.
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

// CategoriesService groups category endpoints.
type CategoriesService struct{}

// List lists categories. Stub.
func (CategoriesService) List(_ context.Context) error { return ErrNotImplemented }

// Get fetches a category. Stub.
func (CategoriesService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// Create creates a category. Stub.
func (CategoriesService) Create(_ context.Context) error { return ErrNotImplemented }

// Update updates a category. Stub.
func (CategoriesService) Update(_ context.Context, _ string) error { return ErrNotImplemented }

// Delete deletes a category. Stub.
func (CategoriesService) Delete(_ context.Context, _ string) error { return ErrNotImplemented }

// CounterpartiesService groups counterparty endpoints.
type CounterpartiesService struct{}

// List lists counterparties. Stub.
func (CounterpartiesService) List(_ context.Context) error { return ErrNotImplemented }

// Get fetches a counterparty. Stub.
func (CounterpartiesService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// Create creates a counterparty. Stub.
func (CounterpartiesService) Create(_ context.Context) error { return ErrNotImplemented }

// Update updates a counterparty. Stub.
func (CounterpartiesService) Update(_ context.Context, _ string) error { return ErrNotImplemented }

// Delete deletes a counterparty. Stub.
func (CounterpartiesService) Delete(_ context.Context, _ string) error { return ErrNotImplemented }

// LegalEntitiesService groups legal-entity endpoints (read-only).
type LegalEntitiesService struct{}

// List lists legal entities. Stub.
func (LegalEntitiesService) List(_ context.Context) error { return ErrNotImplemented }

// Get fetches a legal entity. Stub.
func (LegalEntitiesService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// ConnectionsService groups bank-connection endpoints.
type ConnectionsService struct{}

// List lists connections. Stub.
func (ConnectionsService) List(_ context.Context) error { return ErrNotImplemented }

// Get fetches a connection. Stub.
func (ConnectionsService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// Status returns connection status. Stub.
func (ConnectionsService) Status(_ context.Context, _ string) error { return ErrNotImplemented }

// WebhooksService groups webhook endpoints. Signature verification helper lives
// alongside this service (see verify.go).
type WebhooksService struct{}

// List lists webhooks. Stub.
func (WebhooksService) List(_ context.Context) error { return ErrNotImplemented }

// Get fetches a webhook. Stub.
func (WebhooksService) Get(_ context.Context, _ string) error { return ErrNotImplemented }

// Create creates a webhook. Stub.
func (WebhooksService) Create(_ context.Context) error { return ErrNotImplemented }

// Update updates a webhook. Stub.
func (WebhooksService) Update(_ context.Context, _ string) error { return ErrNotImplemented }

// Delete deletes a webhook. Stub.
func (WebhooksService) Delete(_ context.Context, _ string) error { return ErrNotImplemented }

// ReportsService groups report endpoints.
type ReportsService struct{}

// CashFlow returns the cash-flow report. Stub.
func (ReportsService) CashFlow(_ context.Context) error { return ErrNotImplemented }

// BalanceHistoryService groups balance-history endpoints.
type BalanceHistoryService struct{}

// ListForAccount lists balance history for an account. Stub.
func (BalanceHistoryService) ListForAccount(_ context.Context, _ string) error {
	return ErrNotImplemented
}

// WorkspaceService groups workspace endpoints.
type WorkspaceService struct{}

// Get returns the workspace summary. Stub.
func (WorkspaceService) Get(_ context.Context) error { return ErrNotImplemented }

// MCPService groups the v3 MCP pass-through endpoint.
type MCPService struct{}

// Handle proxies a raw MCP call. Stub.
func (MCPService) Handle(_ context.Context) error { return ErrNotImplemented }

// StatusService groups status endpoints.
type StatusService struct{}

// Status returns the API status. Stub.
func (StatusService) Status(_ context.Context) error { return ErrNotImplemented }

// Whoami returns identity info for the current API key. Stub.
func (StatusService) Whoami(_ context.Context) error { return ErrNotImplemented }
