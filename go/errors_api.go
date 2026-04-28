package tesote

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
