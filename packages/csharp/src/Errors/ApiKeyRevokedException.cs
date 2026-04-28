using System;

namespace Tesote.Sdk.Errors;

/// <summary>401 API_KEY_REVOKED — key was explicitly revoked.</summary>
public sealed class ApiKeyRevokedException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public ApiKeyRevokedException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
