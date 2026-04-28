using System;

namespace Tesote.Sdk.Errors;

/// <summary>503 — documented "pause mode"; raised when no envelope error_code is supplied.</summary>
public sealed class ServiceUnavailableException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public ServiceUnavailableException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
