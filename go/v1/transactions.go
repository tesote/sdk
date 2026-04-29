package v1

import (
	"context"
	"strconv"
	"time"

	tesote "github.com/tesote/sdk/go"
)

// TransactionsService groups v1 transaction endpoints.
type TransactionsService struct {
	client *tesote.Client
}

// TransactionsListOptions tunes GET /v1/accounts/{id}/transactions.
type TransactionsListOptions struct {
	StartDate            string
	EndDate              string
	Scope                string
	Page                 int
	PerPage              int
	TransactionsAfterID  string
	TransactionsBeforeID string
}

func (o TransactionsListOptions) query() map[string]string {
	q := map[string]string{}
	if o.StartDate != "" {
		q["start_date"] = o.StartDate
	}
	if o.EndDate != "" {
		q["end_date"] = o.EndDate
	}
	if o.Scope != "" {
		q["scope"] = o.Scope
	}
	if o.Page > 0 {
		q["page"] = strconv.Itoa(o.Page)
	}
	if o.PerPage > 0 {
		q["per_page"] = strconv.Itoa(o.PerPage)
	}
	if o.TransactionsAfterID != "" {
		q["transactions_after_id"] = o.TransactionsAfterID
	}
	if o.TransactionsBeforeID != "" {
		q["transactions_before_id"] = o.TransactionsBeforeID
	}
	return q
}

// ListForAccount lists transactions for an account. GET /v1/accounts/{id}/transactions.
func (s *TransactionsService) ListForAccount(ctx context.Context, accountID string, opts TransactionsListOptions) (*tesote.TransactionListResponse, error) {
	out := &tesote.TransactionListResponse{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts/"+accountID+"/transactions", tesote.RequestOptions{
		Query: opts.query(),
		Out:   out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single transaction. GET /v1/transactions/{id}.
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
