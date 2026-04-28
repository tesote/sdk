package tesote

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

// RateLimitExceededError is raised when retries against 429 are exhausted.
type RateLimitExceededError struct{ *APIError }

// Is matches the ErrRateLimitExceeded sentinel for errors.Is.
func (e *RateLimitExceededError) Is(target error) bool { return target == ErrRateLimitExceeded }

// Unwrap exposes the embedded *APIError for errors.As.
func (e *RateLimitExceededError) Unwrap() error { return e.APIError }

// UnauthorizedError is raised on 401 UNAUTHORIZED.
type UnauthorizedError struct{ *APIError }

// Is matches the ErrUnauthorized sentinel.
func (e *UnauthorizedError) Is(target error) bool { return target == ErrUnauthorized }

// Unwrap exposes the embedded *APIError.
func (e *UnauthorizedError) Unwrap() error { return e.APIError }

// APIKeyRevokedError is raised on 401 API_KEY_REVOKED.
type APIKeyRevokedError struct{ *APIError }

// Is matches the ErrAPIKeyRevoked sentinel.
func (e *APIKeyRevokedError) Is(target error) bool { return target == ErrAPIKeyRevoked }

// Unwrap exposes the embedded *APIError.
func (e *APIKeyRevokedError) Unwrap() error { return e.APIError }

// WorkspaceSuspendedError is raised on 403 WORKSPACE_SUSPENDED.
type WorkspaceSuspendedError struct{ *APIError }

// Is matches the ErrWorkspaceSuspended sentinel.
func (e *WorkspaceSuspendedError) Is(target error) bool { return target == ErrWorkspaceSuspended }

// Unwrap exposes the embedded *APIError.
func (e *WorkspaceSuspendedError) Unwrap() error { return e.APIError }

// AccountDisabledError is raised on 403 ACCOUNT_DISABLED.
type AccountDisabledError struct{ *APIError }

// Is matches the ErrAccountDisabled sentinel.
func (e *AccountDisabledError) Is(target error) bool { return target == ErrAccountDisabled }

// Unwrap exposes the embedded *APIError.
func (e *AccountDisabledError) Unwrap() error { return e.APIError }

// HistorySyncForbiddenError is raised on 403 HISTORY_SYNC_FORBIDDEN.
type HistorySyncForbiddenError struct{ *APIError }

// Is matches the ErrHistorySyncForbidden sentinel.
func (e *HistorySyncForbiddenError) Is(target error) bool { return target == ErrHistorySyncForbidden }

// Unwrap exposes the embedded *APIError.
func (e *HistorySyncForbiddenError) Unwrap() error { return e.APIError }

// MutationDuringPaginationError is raised on 409 MUTATION_CONFLICT.
type MutationDuringPaginationError struct{ *APIError }

// Is matches the ErrMutationDuringPagination sentinel.
func (e *MutationDuringPaginationError) Is(target error) bool {
	return target == ErrMutationDuringPagination
}

// Unwrap exposes the embedded *APIError.
func (e *MutationDuringPaginationError) Unwrap() error { return e.APIError }

// UnprocessableContentError is raised on 422 UNPROCESSABLE_CONTENT.
type UnprocessableContentError struct{ *APIError }

// Is matches the ErrUnprocessableContent sentinel.
func (e *UnprocessableContentError) Is(target error) bool { return target == ErrUnprocessableContent }

// Unwrap exposes the embedded *APIError.
func (e *UnprocessableContentError) Unwrap() error { return e.APIError }

// InvalidDateRangeError is raised on 422 INVALID_DATE_RANGE.
type InvalidDateRangeError struct{ *APIError }

// Is matches the ErrInvalidDateRange sentinel.
func (e *InvalidDateRangeError) Is(target error) bool { return target == ErrInvalidDateRange }

// Unwrap exposes the embedded *APIError.
func (e *InvalidDateRangeError) Unwrap() error { return e.APIError }

// ServiceUnavailableError is raised on 503 (platform pause mode).
type ServiceUnavailableError struct{ *APIError }

// Is matches the ErrServiceUnavailable sentinel.
func (e *ServiceUnavailableError) Is(target error) bool { return target == ErrServiceUnavailable }

// Unwrap exposes the embedded *APIError.
func (e *ServiceUnavailableError) Unwrap() error { return e.APIError }

// TransportError is the base for network-level failures (no usable HTTP
// response). Subtypes embed it and override Unwrap.
type TransportError struct {
	Op             string
	Message        string
	RequestSummary RequestSummary
	Attempts       int
	Cause          error
}

// Error implements error.
func (e *TransportError) Error() string {
	return fmt.Sprintf("tesote: %s: %s (attempts=%d)", e.Op, e.Message, e.Attempts)
}

// NetworkError is raised on DNS, connection refused, reset.
type NetworkError struct{ *TransportError }

// Unwrap returns the wrapped cause and matches ErrNetwork via errors.Is.
func (e *NetworkError) Unwrap() error { return e.Cause }

// Is satisfies errors.Is for the ErrNetwork sentinel.
func (e *NetworkError) Is(target error) bool { return target == ErrNetwork }

// TimeoutError is raised on connect/read timeout.
type TimeoutError struct{ *TransportError }

// Unwrap returns the wrapped cause.
func (e *TimeoutError) Unwrap() error { return e.Cause }

// Is satisfies errors.Is for the ErrTimeout sentinel.
func (e *TimeoutError) Is(target error) bool { return target == ErrTimeout }

// TLSError is raised on certificate / handshake failures.
type TLSError struct{ *TransportError }

// Unwrap returns the wrapped cause.
func (e *TLSError) Unwrap() error { return e.Cause }

// Is satisfies errors.Is for the ErrTLS sentinel.
func (e *TLSError) Is(target error) bool { return target == ErrTLS }

// ConfigError is raised at client construction for bad SDK config.
type ConfigError struct {
	Field   string
	Message string
}

// Error implements error.
func (e *ConfigError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("tesote: config error: %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("tesote: config error: %s", e.Message)
}

// Is satisfies errors.Is for the ErrConfig sentinel.
func (e *ConfigError) Is(target error) bool { return target == ErrConfig }

// EndpointRemovedError is returned when the SDK exposes a method whose upstream
// endpoint has been removed from this API version.
type EndpointRemovedError struct {
	Method      string
	Replacement string
}

// Error implements error.
func (e *EndpointRemovedError) Error() string {
	if e.Replacement != "" {
		return fmt.Sprintf("tesote: %s removed; use %s", e.Method, e.Replacement)
	}
	return fmt.Sprintf("tesote: %s removed", e.Method)
}

// Is satisfies errors.Is for the ErrEndpointRemoved sentinel.
func (e *EndpointRemovedError) Is(target error) bool { return target == ErrEndpointRemoved }

// envelope mirrors the API error envelope shape.
type envelope struct {
	Error      string `json:"error"`
	ErrorCode  string `json:"error_code"`
	ErrorID    string `json:"error_id"`
	RetryAfter int    `json:"retry_after"`
}

// MapAPIError dispatches a non-2xx response into the narrowest typed error.
// It expects the response body to already be drained into body.
func MapAPIError(resp *http.Response, body []byte, summary RequestSummary) error {
	var env envelope
	// why: best-effort parse; non-JSON bodies fall through to status-based dispatch.
	if len(body) > 0 {
		_ = json.Unmarshal(body, &env)
	}

	requestID := resp.Header.Get("X-Request-Id")
	retryAfter := env.RetryAfter
	if hdr := resp.Header.Get("Retry-After"); hdr != "" {
		if n, err := strconv.Atoi(hdr); err == nil {
			retryAfter = n
		}
	}

	message := env.Error
	if message == "" {
		message = http.StatusText(resp.StatusCode)
		if message == "" {
			message = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
	}

	base := &APIError{
		ErrorCode:      env.ErrorCode,
		Message:        message,
		HTTPStatus:     resp.StatusCode,
		RequestID:      requestID,
		ErrorID:        env.ErrorID,
		RetryAfter:     retryAfter,
		ResponseBody:   string(body),
		RequestSummary: summary,
		Attempts:       1,
	}
	if base.ErrorCode == "" {
		base.ErrorCode = fallbackCode(resp.StatusCode)
	}

	return wrapTyped(base)
}

// wrapTyped picks the narrowest *APIError subtype based on error_code, with
// HTTP-status fallback when the server omits or sends an unknown code.
func wrapTyped(base *APIError) error {
	switch base.ErrorCode {
	case "UNAUTHORIZED":
		return &UnauthorizedError{APIError: base}
	case "API_KEY_REVOKED":
		return &APIKeyRevokedError{APIError: base}
	case "WORKSPACE_SUSPENDED":
		return &WorkspaceSuspendedError{APIError: base}
	case "ACCOUNT_DISABLED":
		return &AccountDisabledError{APIError: base}
	case "HISTORY_SYNC_FORBIDDEN":
		return &HistorySyncForbiddenError{APIError: base}
	case "MUTATION_CONFLICT":
		return &MutationDuringPaginationError{APIError: base}
	case "UNPROCESSABLE_CONTENT":
		return &UnprocessableContentError{APIError: base}
	case "INVALID_DATE_RANGE":
		return &InvalidDateRangeError{APIError: base}
	case "RATE_LIMIT_EXCEEDED":
		return &RateLimitExceededError{APIError: base}
	}
	switch base.HTTPStatus {
	case http.StatusUnauthorized:
		return &UnauthorizedError{APIError: base}
	case http.StatusConflict:
		return &MutationDuringPaginationError{APIError: base}
	case http.StatusUnprocessableEntity:
		return &UnprocessableContentError{APIError: base}
	case http.StatusTooManyRequests:
		return &RateLimitExceededError{APIError: base}
	case http.StatusServiceUnavailable:
		return &ServiceUnavailableError{APIError: base}
	}
	return base
}

func fallbackCode(status int) string {
	switch status {
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusConflict:
		return "MUTATION_CONFLICT"
	case http.StatusUnprocessableEntity:
		return "UNPROCESSABLE_CONTENT"
	case http.StatusTooManyRequests:
		return "RATE_LIMIT_EXCEEDED"
	case http.StatusServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	}
	return fmt.Sprintf("HTTP_%d", status)
}

// RedactBearer formats an API key as "Bearer <last4>" for safe logging.
// Short keys collapse to "Bearer <redacted>".
func RedactBearer(apiKey string) string {
	if len(apiKey) < 4 {
		return "Bearer <redacted>"
	}
	return "Bearer " + apiKey[len(apiKey)-4:]
}
