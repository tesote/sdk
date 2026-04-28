using System;

namespace Tesote.Sdk.Errors;

/// <summary>429 RATE_LIMIT_EXCEEDED — raised after retries are exhausted.</summary>
public sealed class RateLimitExceededException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public RateLimitExceededException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
