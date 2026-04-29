package com.tesote.sdk.errors;

/**
 * Single dispatcher: maps an API error envelope to its typed exception class.
 *
 * <p>If the {@code error_code} is unknown, falls back to {@link ApiException}
 * so callers still get a typed error and full context.
 */
public final class ErrorDispatcher {
    private ErrorDispatcher() {}

    public static ApiException dispatch(
            String message,
            String errorCode,
            int httpStatus,
            String requestId,
            String errorId,
            Integer retryAfter,
            String responseBody,
            RequestSummary requestSummary,
            int attempts,
            Throwable cause
    ) {
        String code = errorCode == null ? "" : errorCode;
        return switch (code) {
            case "UNAUTHORIZED" -> new UnauthorizedException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "API_KEY_REVOKED" -> new ApiKeyRevokedException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "WORKSPACE_SUSPENDED" -> new WorkspaceSuspendedException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "ACCOUNT_DISABLED" -> new AccountDisabledException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "HISTORY_SYNC_FORBIDDEN" -> new HistorySyncForbiddenException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "MUTATION_CONFLICT" -> new MutationDuringPaginationException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "INVALID_DATE_RANGE" -> new InvalidDateRangeException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "MISSING_DATE_RANGE" -> new MissingDateRangeException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "INVALID_CURSOR" -> new InvalidCursorException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "INVALID_COUNT" -> new InvalidCountException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "INVALID_LIMIT" -> new InvalidLimitException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "INVALID_QUERY" -> new InvalidQueryException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "UNPROCESSABLE_CONTENT" -> new UnprocessableContentException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "RATE_LIMIT_EXCEEDED" -> new RateLimitExceededException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "SYNC_RATE_LIMIT_EXCEEDED" -> new SyncRateLimitExceededException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "SYNC_IN_PROGRESS" -> new SyncInProgressException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "INVALID_ORDER_STATE" -> new InvalidOrderStateException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "VALIDATION_ERROR" -> new ValidationException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "BATCH_VALIDATION_ERROR" -> new BatchValidationException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "BANK_SUBMISSION_ERROR" -> new BankSubmissionException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "BANK_UNDER_MAINTENANCE" -> new BankUnderMaintenanceException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "INTERNAL_ERROR" -> new InternalErrorException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "ACCOUNT_NOT_FOUND" -> new AccountNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "TRANSACTION_NOT_FOUND" -> new TransactionNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "SYNC_SESSION_NOT_FOUND" -> new SyncSessionNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "PAYMENT_METHOD_NOT_FOUND" -> new PaymentMethodNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "TRANSACTION_ORDER_NOT_FOUND" -> new TransactionOrderNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "BATCH_NOT_FOUND" -> new BatchNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "BANK_CONNECTION_NOT_FOUND" -> new BankConnectionNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "CATEGORY_NOT_FOUND" -> new CategoryNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "COUNTERPARTY_NOT_FOUND" -> new CounterpartyNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "LEGAL_ENTITY_NOT_FOUND" -> new LegalEntityNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "WEBHOOK_NOT_FOUND" -> new WebhookNotFoundException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            default -> {
                // why: 503 with no envelope code is the documented "pause mode"
                // signal — surface it as ServiceUnavailableException so callers
                // can dispatch on type.
                if (httpStatus == 503) {
                    yield new ServiceUnavailableException(
                            message, code, httpStatus, requestId, errorId, retryAfter,
                            responseBody, requestSummary, attempts, cause);
                }
                yield new ApiException(
                        message, code, httpStatus, requestId, errorId, retryAfter,
                        responseBody, requestSummary, attempts, cause);
            }
        };
    }
}
