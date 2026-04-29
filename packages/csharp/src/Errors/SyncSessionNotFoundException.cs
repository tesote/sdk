using System;

namespace Tesote.Sdk.Errors;

/// <summary>404 SYNC_SESSION_NOT_FOUND.</summary>
public sealed class SyncSessionNotFoundException : NotFoundException
{
    /// <summary>Construct with the full required-field set.</summary>
    public SyncSessionNotFoundException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
