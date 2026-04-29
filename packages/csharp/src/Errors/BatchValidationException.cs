using System;

namespace Tesote.Sdk.Errors;

/// <summary>400 BATCH_VALIDATION_ERROR — batch payload failed validation.</summary>
public sealed class BatchValidationException : ValidationException
{
    /// <summary>Construct with the full required-field set.</summary>
    public BatchValidationException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
