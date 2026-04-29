using System;

namespace Tesote.Sdk.Errors;

/// <summary>400 VALIDATION_ERROR — request body failed validation.</summary>
public class ValidationException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public ValidationException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
