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
	ErrAccountNotFound          = errors.New("tesote: account not found")
	ErrTransactionNotFound      = errors.New("tesote: transaction not found")
	ErrSyncSessionNotFound      = errors.New("tesote: sync session not found")
	ErrPaymentMethodNotFound    = errors.New("tesote: payment method not found")
	ErrTransactionOrderNotFound = errors.New("tesote: transaction order not found")
	ErrBatchNotFound            = errors.New("tesote: batch not found")
	ErrBankConnectionNotFound   = errors.New("tesote: bank connection not found")
	ErrInvalidCursor            = errors.New("tesote: invalid cursor")
	ErrInvalidCount             = errors.New("tesote: invalid count")
	ErrInvalidLimit             = errors.New("tesote: invalid limit")
	ErrInvalidQuery             = errors.New("tesote: invalid query")
	ErrMissingDateRange         = errors.New("tesote: missing date range")
	ErrSyncInProgress           = errors.New("tesote: sync in progress")
	ErrSyncRateLimitExceeded    = errors.New("tesote: sync rate limit exceeded")
	ErrBankUnderMaintenance     = errors.New("tesote: bank under maintenance")
	ErrValidation               = errors.New("tesote: validation error")
	ErrInvalidOrderState        = errors.New("tesote: invalid order state")
	ErrBankSubmission           = errors.New("tesote: bank submission error")
	ErrBatchValidation          = errors.New("tesote: batch validation error")
	ErrInternal                 = errors.New("tesote: internal error")
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
