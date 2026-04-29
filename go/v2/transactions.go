package v2

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	tesote "github.com/tesote/sdk/go"
)

// TransactionsService groups v2 transaction endpoints.
type TransactionsService struct {
	client *tesote.Client
}

// TransactionsFilter holds the shared filter parameters for v2 transaction queries.
type TransactionsFilter struct {
	StartDate             string
	EndDate               string
	Scope                 string
	Page                  int
	PerPage               int
	TransactionsAfterID   string
	TransactionsBeforeID  string
	TransactionDateAfter  string
	TransactionDateBefore string
	CreatedAfter          string
	UpdatedAfter          string
	AmountMin             string
	AmountMax             string
	Amount                string
	Status                string
	CategoryID            string
	CounterpartyID        string
	Q                     string
	Type                  string
	ReferenceCode         string
}

func (f TransactionsFilter) query() map[string]string {
	q := map[string]string{}
	set := func(k, v string) {
		if v != "" {
			q[k] = v
		}
	}
	set("start_date", f.StartDate)
	set("end_date", f.EndDate)
	set("scope", f.Scope)
	if f.Page > 0 {
		q["page"] = strconv.Itoa(f.Page)
	}
	if f.PerPage > 0 {
		q["per_page"] = strconv.Itoa(f.PerPage)
	}
	set("transactions_after_id", f.TransactionsAfterID)
	set("transactions_before_id", f.TransactionsBeforeID)
	set("transaction_date_after", f.TransactionDateAfter)
	set("transaction_date_before", f.TransactionDateBefore)
	set("created_after", f.CreatedAfter)
	set("updated_after", f.UpdatedAfter)
	set("amount_min", f.AmountMin)
	set("amount_max", f.AmountMax)
	set("amount", f.Amount)
	set("status", f.Status)
	set("category_id", f.CategoryID)
	set("counterparty_id", f.CounterpartyID)
	set("q", f.Q)
	set("type", f.Type)
	set("reference_code", f.ReferenceCode)
	return q
}

// ListForAccount lists transactions for an account. GET /v2/accounts/{id}/transactions.
func (s *TransactionsService) ListForAccount(ctx context.Context, accountID string, filter TransactionsFilter) (*tesote.TransactionListResponse, error) {
	out := &tesote.TransactionListResponse{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts/"+accountID+"/transactions", tesote.RequestOptions{
		Query:    filter.query(),
		Out:      out,
		CacheTTL: time.Minute,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single transaction. GET /v2/transactions/{id}.
func (s *TransactionsService) Get(ctx context.Context, id string) (*tesote.Transaction, error) {
	out := &tesote.Transaction{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/transactions/"+id, tesote.RequestOptions{
		Out:      out,
		CacheTTL: 5 * time.Minute,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExportFormat is the wire enum for /v2/accounts/{id}/transactions/export `format`.
type ExportFormat string

// Supported export formats.
const (
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatJSON ExportFormat = "json"
)

// ExportOptions tunes GET /v2/accounts/{id}/transactions/export.
type ExportOptions struct {
	Filter TransactionsFilter
	Format ExportFormat
}

// ExportResult is the raw export payload (CSV or pretty JSON bytes).
type ExportResult struct {
	Body        []byte
	ContentType string
	Filename    string
}

// Export downloads transactions as CSV or JSON.
// GET /v2/accounts/{id}/transactions/export. Returns the raw payload bytes.
func (s *TransactionsService) Export(ctx context.Context, accountID string, opts ExportOptions) (*ExportResult, error) {
	q := opts.Filter.query()
	if opts.Format != "" {
		q["format"] = string(opts.Format)
	}
	resp, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts/"+accountID+"/transactions/export", tesote.RequestOptions{
		Query: q,
	})
	if err != nil {
		return nil, err
	}
	res := &ExportResult{
		Body:        append([]byte(nil), resp.Body...),
		ContentType: resp.Header.Get("Content-Type"),
	}
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		// why: extract filename from Content-Disposition without bringing in mime/multipart.
		res.Filename = parseFilename(cd)
	}
	return res, nil
}

// parseFilename pulls a filename from a Content-Disposition header (best-effort).
func parseFilename(cd string) string {
	const marker = "filename="
	idx := bytes.Index([]byte(cd), []byte(marker))
	if idx < 0 {
		return ""
	}
	rest := cd[idx+len(marker):]
	if len(rest) > 0 && rest[0] == '"' {
		end := bytes.IndexByte([]byte(rest[1:]), '"')
		if end < 0 {
			return rest[1:]
		}
		return rest[1 : 1+end]
	}
	end := bytes.IndexAny([]byte(rest), "; ")
	if end < 0 {
		return rest
	}
	return rest[:end]
}

// SyncOptions tunes the body of POST /v2/accounts/{id}/transactions/sync (and legacy).
type SyncOptions struct {
	Count          int
	Cursor         string
	Options        SyncSubOptions
	IdempotencyKey string
}

// SyncSubOptions is the nested `options` block on a sync request.
type SyncSubOptions struct {
	IncludeRunningBalance bool `json:"include_running_balance,omitempty"`
}

func (o SyncOptions) body() map[string]any {
	body := map[string]any{}
	if o.Count > 0 {
		body["count"] = o.Count
	}
	if o.Cursor != "" {
		body["cursor"] = o.Cursor
	}
	if o.Options.IncludeRunningBalance {
		body["options"] = o.Options
	}
	return body
}

// Sync runs a transaction sync against an account.
// POST /v2/accounts/{id}/transactions/sync.
func (s *TransactionsService) Sync(ctx context.Context, accountID string, opts SyncOptions) (*tesote.TransactionSyncResponse, error) {
	out := &tesote.TransactionSyncResponse{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/accounts/"+accountID+"/transactions/sync", tesote.RequestOptions{
		Body:           opts.body(),
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SyncLegacy hits the non-nested legacy sync path.
// POST /v2/transactions/sync.
func (s *TransactionsService) SyncLegacy(ctx context.Context, opts SyncOptions) (*tesote.TransactionSyncResponse, error) {
	out := &tesote.TransactionSyncResponse{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/transactions/sync", tesote.RequestOptions{
		Body:           opts.body(),
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BulkOptions tunes POST /v2/transactions/bulk.
type BulkOptions struct {
	AccountIDs     []string
	Page           int
	PerPage        int
	Limit          int
	Offset         int
	IdempotencyKey string
}

// Bulk fetches transactions for multiple accounts in one call.
// POST /v2/transactions/bulk.
func (s *TransactionsService) Bulk(ctx context.Context, opts BulkOptions) (*tesote.BulkTransactionsResponse, error) {
	if len(opts.AccountIDs) == 0 {
		return nil, fmt.Errorf("tesote/v2: BulkOptions.AccountIDs must not be empty")
	}
	body := map[string]any{
		"account_ids": opts.AccountIDs,
	}
	if opts.Page > 0 {
		body["page"] = opts.Page
	}
	if opts.PerPage > 0 {
		body["per_page"] = opts.PerPage
	}
	if opts.Limit > 0 {
		body["limit"] = opts.Limit
	}
	if opts.Offset > 0 {
		body["offset"] = opts.Offset
	}
	out := &tesote.BulkTransactionsResponse{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/transactions/bulk", tesote.RequestOptions{
		Body:           body,
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchOptions tunes GET /v2/transactions/search.
type SearchOptions struct {
	Q         string
	AccountID string
	Limit     int
	Offset    int
	Filter    TransactionsFilter
}

// Search runs a text search across transactions. GET /v2/transactions/search.
func (s *TransactionsService) Search(ctx context.Context, opts SearchOptions) (*tesote.SearchTransactionsResponse, error) {
	if opts.Q == "" {
		return nil, fmt.Errorf("tesote/v2: SearchOptions.Q is required")
	}
	q := opts.Filter.query()
	q["q"] = opts.Q
	if opts.AccountID != "" {
		q["account_id"] = opts.AccountID
	}
	if opts.Limit > 0 {
		q["limit"] = strconv.Itoa(opts.Limit)
	}
	if opts.Offset > 0 {
		q["offset"] = strconv.Itoa(opts.Offset)
	}
	out := &tesote.SearchTransactionsResponse{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/transactions/search", tesote.RequestOptions{
		Query: q,
		Out:   out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
