using System;

namespace Tesote.Sdk.Errors;

/// <summary>422 INVALID_QUERY.</summary>
public sealed class InvalidQueryException : UnprocessableContentException
{
    /// <summary>Construct with the full required-field set.</summary>
    public InvalidQueryException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
