package tesote

import (
	"errors"
	"fmt"
)

// Sentinel errors. Use with errors.Is for programmatic dispatch:
//
//	if errors.Is(err, tesote.ErrRateLimitExceeded) { ... }
var (
	ErrUnauthorized             = errors.New("tesote: unauthorized")
	ErrAPIKeyRevoked            = errors.New("tesote: api key revoked")
	ErrWorkspaceSuspended       = errors.New("tesote: workspace suspended")
	ErrAccountDisabled          = errors.New("tesote: account disabled")
	ErrHistorySyncForbidden     = errors.New("tesote: history sync forbidden")
	ErrMutationDuringPagination = errors.New("tesote: mutation during pagination")
	ErrUnprocessableContent     = errors.New("tesote: unprocessable content")
	ErrInvalidDateRange         = errors.New("tesote: invalid date range")
	ErrRateLimitExceeded        = errors.New("tesote: rate limit exceeded")
	ErrServiceUnavailable       = errors.New("tesote: service unavailable")
	ErrNetwork                  = errors.New("tesote: network error")
	ErrTimeout                  = errors.New("tesote: timeout")
	ErrTLS                      = errors.New("tesote: tls error")
	ErrConfig                   = errors.New("tesote: config error")
	ErrEndpointRemoved          = errors.New("tesote: endpoint removed")
)

// RequestSummary is the redacted request snapshot attached to every error.
// Bearer tokens MUST be redacted to "Bearer <last4>" before construction.
type RequestSummary struct {
	Method        string            `json:"method"`
	Path          string            `json:"path"`
	Query         map[string]string `json:"query,omitempty"`
	BodyShape     string            `json:"body_shape,omitempty"`
	Authorization string            `json:"authorization,omitempty"`
}

// APIError is the base type for every server-returned typed error. Subtypes
// (RateLimitExceededError, etc.) embed *APIError and override Unwrap to expose
// the matching sentinel for errors.Is.
type APIError struct {
	ErrorCode      string
	Message        string
	HTTPStatus     int
	RequestID      string
	ErrorID        string
	RetryAfter     int
	ResponseBody   string
	RequestSummary RequestSummary
	Attempts       int
}

// Error implements the error interface with a greppable first line.
func (e *APIError) Error() string {
	return fmt.Sprintf(
		"tesote: %s (error_code=%s http_status=%d request_id=%s attempts=%d)",
		e.Message, e.ErrorCode, e.HTTPStatus, e.RequestID, e.Attempts,
	)
}

// RedactBearer formats an API key as "Bearer <last4>" for safe logging.
// Short keys collapse to "Bearer <redacted>".
func RedactBearer(apiKey string) string {
	if len(apiKey) < 4 {
		return "Bearer <redacted>"
	}
	return "Bearer " + apiKey[len(apiKey)-4:]
}
