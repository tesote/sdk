package v2

import (
	"context"
	"strconv"

	tesote "github.com/tesote/sdk/go"
)

// SyncSessionsService groups v2 sync-session endpoints.
type SyncSessionsService struct {
	client *tesote.Client
}

// SyncSessionsListOptions tunes GET /v2/accounts/{id}/sync_sessions.
type SyncSessionsListOptions struct {
	Limit  int
	Offset int
	Status string
}

func (o SyncSessionsListOptions) query() map[string]string {
	q := map[string]string{}
	if o.Limit > 0 {
		q["limit"] = strconv.Itoa(o.Limit)
	}
	if o.Offset > 0 {
		q["offset"] = strconv.Itoa(o.Offset)
	}
	if o.Status != "" {
		q["status"] = o.Status
	}
	return q
}

// List lists sync sessions for an account.
// GET /v2/accounts/{id}/sync_sessions.
func (s *SyncSessionsService) List(ctx context.Context, accountID string, opts SyncSessionsListOptions) (*tesote.SyncSessionListResponse, error) {
	out := &tesote.SyncSessionListResponse{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts/"+accountID+"/sync_sessions", tesote.RequestOptions{
		Query: opts.query(),
		Out:   out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single sync session.
// GET /v2/accounts/{id}/sync_sessions/{session_id}.
func (s *SyncSessionsService) Get(ctx context.Context, accountID, sessionID string) (*tesote.SyncSession, error) {
	out := &tesote.SyncSession{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/accounts/"+accountID+"/sync_sessions/"+sessionID, tesote.RequestOptions{
		Out: out,
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
