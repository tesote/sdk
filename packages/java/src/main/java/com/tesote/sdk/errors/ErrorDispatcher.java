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
            case "UNPROCESSABLE_CONTENT" -> new UnprocessableContentException(
                    message, code, httpStatus, requestId, errorId, retryAfter,
                    responseBody, requestSummary, attempts, cause);
            case "RATE_LIMIT_EXCEEDED" -> new RateLimitExceededException(
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
