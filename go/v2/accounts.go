package v2

import (
	"context"
	"strconv"
	"time"

	tesote "github.com/tesote/sdk/go"
)

// AccountsService groups v2 account endpoints.
type AccountsService struct {
	client *tesote.Client
}

// AccountsListOptions tunes GET /v2/accounts.
type AccountsListOptions struct {
	Page    int
	PerPage int
	Include string
	Sort    string
}

func (o AccountsListOptions) query() map[string]string {
	q := map[string]string{}
	if o.Page > 0 {
		q["page"] = strconv.Itoa(o.Page)
	}
	if o.PerPage > 0 {
		q["per_page"] = strconv.Itoa(o.PerPage)
	}
	if o.Include != "" {
		q["include"] = o.Include
	}
	if o.Sort != "" {
		q["sort"] = o.Sort
	}
	return q
}

// List lists accounts. GET /v2/accounts.
func (s *AccountsService) List(ctx context.Context, opts AccountsListOptions) (*tesote.AccountListResponse, error) {
	out := &tesote.AccountListResponse{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts", tesote.RequestOptions{
		Query:    opts.query(),
		Out:      out,
		CacheTTL: time.Minute,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single account. GET /v2/accounts/{id}.
func (s *AccountsService) Get(ctx context.Context, id string) (*tesote.Account, error) {
	out := &tesote.Account{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts/"+id, tesote.RequestOptions{
		Out:      out,
		CacheTTL: 5 * time.Minute,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountSyncOptions controls POST /v2/accounts/{id}/sync.
type AccountSyncOptions struct {
	// IdempotencyKey overrides the auto-generated key for the request.
	IdempotencyKey string
}

// Sync triggers an account sync. POST /v2/accounts/{id}/sync.
func (s *AccountsService) Sync(ctx context.Context, id string, opts AccountSyncOptions) (*tesote.AccountSyncResponse, error) {
	out := &tesote.AccountSyncResponse{}
	_, err := s.client.Do(ctx, "POST", pathPrefix+"/accounts/"+id+"/sync", tesote.RequestOptions{
		IdempotencyKey: opts.IdempotencyKey,
		Out:            out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
