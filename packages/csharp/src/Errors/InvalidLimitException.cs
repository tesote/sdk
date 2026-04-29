using System;

namespace Tesote.Sdk.Errors;

/// <summary>422 INVALID_LIMIT.</summary>
public sealed class InvalidLimitException : UnprocessableContentException
{
    /// <summary>Construct with the full required-field set.</summary>
    public InvalidLimitException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
