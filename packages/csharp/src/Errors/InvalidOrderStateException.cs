using System;

namespace Tesote.Sdk.Errors;

/// <summary>409 INVALID_ORDER_STATE — order cannot transition from current state.</summary>
public sealed class InvalidOrderStateException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public InvalidOrderStateException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
