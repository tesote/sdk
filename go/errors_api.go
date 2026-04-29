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

// AccountNotFoundError is raised on 404 ACCOUNT_NOT_FOUND.
type AccountNotFoundError struct{ *APIError }

// Is matches the ErrAccountNotFound sentinel.
func (e *AccountNotFoundError) Is(target error) bool { return target == ErrAccountNotFound }

// Unwrap exposes the embedded *APIError.
func (e *AccountNotFoundError) Unwrap() error { return e.APIError }

// TransactionNotFoundError is raised on 404 TRANSACTION_NOT_FOUND.
type TransactionNotFoundError struct{ *APIError }

// Is matches the ErrTransactionNotFound sentinel.
func (e *TransactionNotFoundError) Is(target error) bool { return target == ErrTransactionNotFound }

// Unwrap exposes the embedded *APIError.
func (e *TransactionNotFoundError) Unwrap() error { return e.APIError }

// SyncSessionNotFoundError is raised on 404 SYNC_SESSION_NOT_FOUND.
type SyncSessionNotFoundError struct{ *APIError }

// Is matches the ErrSyncSessionNotFound sentinel.
func (e *SyncSessionNotFoundError) Is(target error) bool { return target == ErrSyncSessionNotFound }

// Unwrap exposes the embedded *APIError.
func (e *SyncSessionNotFoundError) Unwrap() error { return e.APIError }

// PaymentMethodNotFoundError is raised on 404 PAYMENT_METHOD_NOT_FOUND.
type PaymentMethodNotFoundError struct{ *APIError }

// Is matches the ErrPaymentMethodNotFound sentinel.
func (e *PaymentMethodNotFoundError) Is(target error) bool {
	return target == ErrPaymentMethodNotFound
}

// Unwrap exposes the embedded *APIError.
func (e *PaymentMethodNotFoundError) Unwrap() error { return e.APIError }

// TransactionOrderNotFoundError is raised on 404 TRANSACTION_ORDER_NOT_FOUND.
type TransactionOrderNotFoundError struct{ *APIError }

// Is matches the ErrTransactionOrderNotFound sentinel.
func (e *TransactionOrderNotFoundError) Is(target error) bool {
	return target == ErrTransactionOrderNotFound
}

// Unwrap exposes the embedded *APIError.
func (e *TransactionOrderNotFoundError) Unwrap() error { return e.APIError }

// BatchNotFoundError is raised on 404 BATCH_NOT_FOUND.
type BatchNotFoundError struct{ *APIError }

// Is matches the ErrBatchNotFound sentinel.
func (e *BatchNotFoundError) Is(target error) bool { return target == ErrBatchNotFound }

// Unwrap exposes the embedded *APIError.
func (e *BatchNotFoundError) Unwrap() error { return e.APIError }

// BankConnectionNotFoundError is raised on 404 BANK_CONNECTION_NOT_FOUND.
type BankConnectionNotFoundError struct{ *APIError }

// Is matches the ErrBankConnectionNotFound sentinel.
func (e *BankConnectionNotFoundError) Is(target error) bool {
	return target == ErrBankConnectionNotFound
}

// Unwrap exposes the embedded *APIError.
func (e *BankConnectionNotFoundError) Unwrap() error { return e.APIError }

// InvalidCursorError is raised on 422 INVALID_CURSOR.
type InvalidCursorError struct{ *APIError }

// Is matches the ErrInvalidCursor sentinel.
func (e *InvalidCursorError) Is(target error) bool { return target == ErrInvalidCursor }

// Unwrap exposes the embedded *APIError.
func (e *InvalidCursorError) Unwrap() error { return e.APIError }

// InvalidCountError is raised on 422 INVALID_COUNT.
type InvalidCountError struct{ *APIError }

// Is matches the ErrInvalidCount sentinel.
func (e *InvalidCountError) Is(target error) bool { return target == ErrInvalidCount }

// Unwrap exposes the embedded *APIError.
func (e *InvalidCountError) Unwrap() error { return e.APIError }

// InvalidLimitError is raised on 422 INVALID_LIMIT.
type InvalidLimitError struct{ *APIError }

// Is matches the ErrInvalidLimit sentinel.
func (e *InvalidLimitError) Is(target error) bool { return target == ErrInvalidLimit }

// Unwrap exposes the embedded *APIError.
func (e *InvalidLimitError) Unwrap() error { return e.APIError }

// InvalidQueryError is raised on 422 INVALID_QUERY.
type InvalidQueryError struct{ *APIError }

// Is matches the ErrInvalidQuery sentinel.
func (e *InvalidQueryError) Is(target error) bool { return target == ErrInvalidQuery }

// Unwrap exposes the embedded *APIError.
func (e *InvalidQueryError) Unwrap() error { return e.APIError }

// MissingDateRangeError is raised on 422 MISSING_DATE_RANGE.
type MissingDateRangeError struct{ *APIError }

// Is matches the ErrMissingDateRange sentinel.
func (e *MissingDateRangeError) Is(target error) bool { return target == ErrMissingDateRange }

// Unwrap exposes the embedded *APIError.
func (e *MissingDateRangeError) Unwrap() error { return e.APIError }

// SyncInProgressError is raised on 409 SYNC_IN_PROGRESS.
type SyncInProgressError struct{ *APIError }

// Is matches the ErrSyncInProgress sentinel.
func (e *SyncInProgressError) Is(target error) bool { return target == ErrSyncInProgress }

// Unwrap exposes the embedded *APIError.
func (e *SyncInProgressError) Unwrap() error { return e.APIError }

// SyncRateLimitExceededError is raised on 429 SYNC_RATE_LIMIT_EXCEEDED.
type SyncRateLimitExceededError struct{ *APIError }

// Is matches the ErrSyncRateLimitExceeded sentinel.
func (e *SyncRateLimitExceededError) Is(target error) bool {
	return target == ErrSyncRateLimitExceeded
}

// Unwrap exposes the embedded *APIError.
func (e *SyncRateLimitExceededError) Unwrap() error { return e.APIError }

// BankUnderMaintenanceError is raised on 503 BANK_UNDER_MAINTENANCE.
type BankUnderMaintenanceError struct{ *APIError }

// Is matches the ErrBankUnderMaintenance sentinel.
func (e *BankUnderMaintenanceError) Is(target error) bool { return target == ErrBankUnderMaintenance }

// Unwrap exposes the embedded *APIError.
func (e *BankUnderMaintenanceError) Unwrap() error { return e.APIError }

// ValidationError is raised on 400 VALIDATION_ERROR.
type ValidationError struct{ *APIError }

// Is matches the ErrValidation sentinel.
func (e *ValidationError) Is(target error) bool { return target == ErrValidation }

// Unwrap exposes the embedded *APIError.
func (e *ValidationError) Unwrap() error { return e.APIError }

// InvalidOrderStateError is raised on 409 INVALID_ORDER_STATE.
type InvalidOrderStateError struct{ *APIError }

// Is matches the ErrInvalidOrderState sentinel.
func (e *InvalidOrderStateError) Is(target error) bool { return target == ErrInvalidOrderState }

// Unwrap exposes the embedded *APIError.
func (e *InvalidOrderStateError) Unwrap() error { return e.APIError }

// BankSubmissionError is raised on 422 BANK_SUBMISSION_ERROR.
type BankSubmissionError struct{ *APIError }

// Is matches the ErrBankSubmission sentinel.
func (e *BankSubmissionError) Is(target error) bool { return target == ErrBankSubmission }

// Unwrap exposes the embedded *APIError.
func (e *BankSubmissionError) Unwrap() error { return e.APIError }

// BatchValidationError is raised on 400 BATCH_VALIDATION_ERROR.
type BatchValidationError struct{ *APIError }

// Is matches the ErrBatchValidation sentinel.
func (e *BatchValidationError) Is(target error) bool { return target == ErrBatchValidation }

// Unwrap exposes the embedded *APIError.
func (e *BatchValidationError) Unwrap() error { return e.APIError }

// InternalError is raised on 500 INTERNAL_ERROR.
type InternalError struct{ *APIError }

// Is matches the ErrInternal sentinel.
func (e *InternalError) Is(target error) bool { return target == ErrInternal }

// Unwrap exposes the embedded *APIError.
func (e *InternalError) Unwrap() error { return e.APIError }
