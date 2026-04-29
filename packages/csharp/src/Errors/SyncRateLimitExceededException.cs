using System;

namespace Tesote.Sdk.Errors;

/// <summary>429 SYNC_RATE_LIMIT_EXCEEDED — bank-connection sync rate limit hit.</summary>
public sealed class SyncRateLimitExceededException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public SyncRateLimitExceededException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
