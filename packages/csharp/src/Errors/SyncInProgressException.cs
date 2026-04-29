using System;

namespace Tesote.Sdk.Errors;

/// <summary>409 SYNC_IN_PROGRESS — another sync session is currently active.</summary>
public sealed class SyncInProgressException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public SyncInProgressException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
