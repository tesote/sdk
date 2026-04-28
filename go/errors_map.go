package tesote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

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
