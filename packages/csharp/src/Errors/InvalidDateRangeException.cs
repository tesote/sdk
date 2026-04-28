using System;

namespace Tesote.Sdk.Errors;

/// <summary>422 INVALID_DATE_RANGE — subtype of <see cref="UnprocessableContentException"/>.</summary>
public sealed class InvalidDateRangeException : UnprocessableContentException
{
    /// <summary>Construct with the full required-field set.</summary>
    public InvalidDateRangeException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
