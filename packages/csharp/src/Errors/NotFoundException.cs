using System;

namespace Tesote.Sdk.Errors;

/// <summary>404 family — resource not found. Concrete subclasses identify which resource.</summary>
public class NotFoundException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public NotFoundException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
