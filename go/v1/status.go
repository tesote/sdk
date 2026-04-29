package v1

import (
	"context"

	tesote "github.com/tesote/sdk/go"
)

// StatusService groups v1 status endpoints.
type StatusService struct {
	client *tesote.Client
}

// Status returns API status. GET /status (unauthenticated).
func (s *StatusService) Status(ctx context.Context) (*tesote.StatusResponse, error) {
	out := &tesote.StatusResponse{}
	_, err := s.client.Do(ctx, "GET", "/status", tesote.RequestOptions{Out: out})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Whoami returns the API key's identity. GET /whoami.
func (s *StatusService) Whoami(ctx context.Context) (*tesote.WhoamiResponse, error) {
	out := &tesote.WhoamiResponse{}
	_, err := s.client.Do(ctx, "GET", "/whoami", tesote.RequestOptions{Out: out})
	if err != nil {
		return nil, err
	}
	return out, nil
}
