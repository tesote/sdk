using System;

namespace Tesote.Sdk.Errors;

/// <summary>403 HISTORY_SYNC_FORBIDDEN — caller cannot trigger history sync.</summary>
public sealed class HistorySyncForbiddenException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public HistorySyncForbiddenException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
