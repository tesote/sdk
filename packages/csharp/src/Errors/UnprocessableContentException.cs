using System;

namespace Tesote.Sdk.Errors;

/// <summary>422 UNPROCESSABLE_CONTENT — request body was syntactically valid but semantically rejected.</summary>
public class UnprocessableContentException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public UnprocessableContentException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
