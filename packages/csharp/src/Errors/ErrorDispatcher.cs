using System;

namespace Tesote.Sdk.Errors;

/// <summary>
/// Single dispatcher: maps an API error envelope to its typed exception class.
/// If the <c>error_code</c> is unknown, falls back to <see cref="ApiException"/>
/// so callers still get a typed error and full context.
/// </summary>
public static class ErrorDispatcher
{
    /// <summary>Build the matching typed exception for the supplied envelope fields.</summary>
    public static ApiException Dispatch(
        string? message,
        string? errorCode,
        int httpStatus,
        string? requestId,
        string? errorId,
        int? retryAfter,
        string? responseBody,
        RequestSummary? requestSummary,
        int attempts,
        Exception? cause)
    {
        var code = errorCode ?? string.Empty;
        return code switch
        {
            "UNAUTHORIZED" => new UnauthorizedException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
            "API_KEY_REVOKED" => new ApiKeyRevokedException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
            "WORKSPACE_SUSPENDED" => new WorkspaceSuspendedException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
            "ACCOUNT_DISABLED" => new AccountDisabledException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
            "HISTORY_SYNC_FORBIDDEN" => new HistorySyncForbiddenException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
            "MUTATION_CONFLICT" => new MutationDuringPaginationException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
            "INVALID_DATE_RANGE" => new InvalidDateRangeException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
            "UNPROCESSABLE_CONTENT" => new UnprocessableContentException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
            "RATE_LIMIT_EXCEEDED" => new RateLimitExceededException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
            // why: 503 with no envelope code is the documented "pause mode" signal —
            // surface it as ServiceUnavailableException so callers can dispatch on type.
            _ when httpStatus == 503 => new ServiceUnavailableException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
            _ => new ApiException(
                message, code, httpStatus, requestId, errorId, retryAfter,
                responseBody, requestSummary, attempts, cause),
        };
    }
}
