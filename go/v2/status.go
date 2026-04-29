package v2

import (
	"context"

	tesote "github.com/tesote/sdk/go"
)

// StatusService groups v2 status endpoints.
type StatusService struct {
	client *tesote.Client
}

// Status returns API status. GET /v2/status (unauthenticated).
func (s *StatusService) Status(ctx context.Context) (*tesote.StatusResponse, error) {
	out := &tesote.StatusResponse{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/status", tesote.RequestOptions{Out: out})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Whoami returns the API key's identity. GET /v2/whoami.
func (s *StatusService) Whoami(ctx context.Context) (*tesote.WhoamiResponse, error) {
	out := &tesote.WhoamiResponse{}
	_, err := s.client.Do(ctx, "GET", pathPrefix+"/whoami", tesote.RequestOptions{Out: out})
	if err != nil {
		return nil, err
	}
	return out, nil
}
