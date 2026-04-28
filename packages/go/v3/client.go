// Package v3 implements the v3 API surface (paths under /api/v3).
package v3

import (
	tesote "github.com/tesote/sdk/go"
)

// pathPrefix is the URL prefix every v3 request uses, appended after the base URL.
const pathPrefix = "/v3"

// V3Client groups the v3 resource services. All services share the same
// underlying transport.
type V3Client struct {
	transport *tesote.Client

	Accounts *AccountsService
}

// New builds a V3Client around an existing tesote.Client.
func New(t *tesote.Client) *V3Client {
	c := &V3Client{transport: t}
	c.Accounts = &AccountsService{client: t}
	return c
}

// Transport exposes the underlying *tesote.Client for advanced callers (e.g.
// reading LastRateLimit).
func (c *V3Client) Transport() *tesote.Client { return c.transport }
