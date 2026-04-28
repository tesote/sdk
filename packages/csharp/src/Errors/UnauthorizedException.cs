using System;

namespace Tesote.Sdk.Errors;

/// <summary>401 UNAUTHORIZED — bearer token rejected.</summary>
public sealed class UnauthorizedException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public UnauthorizedException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
