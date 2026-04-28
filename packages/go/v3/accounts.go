package v3

import (
	"context"
	"errors"
	"strconv"
	"time"

	tesote "github.com/tesote/sdk/go"
)

// Account is the v3 account resource.
type Account struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Currency  string `json:"currency"`
	IBAN      string `json:"iban,omitempty"`
	Balance   string `json:"balance,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ListOptions tunes a paged Accounts.List call.
type ListOptions struct {
	Cursor   string
	PageSize int
	// CacheTTL, if > 0, opts in to the transport's TTL cache for this call.
	CacheTTL time.Duration
}

// AccountsList is one page of the cursor-paginated accounts listing.
type AccountsList struct {
	Data       []Account `json:"data"`
	NextCursor string    `json:"next_cursor,omitempty"`
}

// AccountsService groups account endpoints.
type AccountsService struct {
	client *tesote.Client
}

// List returns one page of accounts.
func (s *AccountsService) List(ctx context.Context, opts ListOptions) (*AccountsList, error) {
	query := map[string]string{}
	if opts.Cursor != "" {
		query["cursor"] = opts.Cursor
	}
	if opts.PageSize > 0 {
		query["page_size"] = strconv.Itoa(opts.PageSize)
	}
	out := &AccountsList{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts", tesote.RequestOptions{
		Query:    query,
		CacheTTL: opts.CacheTTL,
		Out:      out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns a single account by ID.
func (s *AccountsService) Get(ctx context.Context, id string) (*Account, error) {
	if id == "" {
		return nil, &tesote.ConfigError{Field: "id", Message: "must not be empty"}
	}
	out := &Account{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts/"+id, tesote.RequestOptions{Out: out})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Sync triggers a sync of the given account. Stub.
func (s *AccountsService) Sync(_ context.Context, _ string) error {
	return errors.New("not implemented")
}
