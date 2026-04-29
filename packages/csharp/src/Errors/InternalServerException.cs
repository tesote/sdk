using System;

namespace Tesote.Sdk.Errors;

/// <summary>500 INTERNAL_ERROR — server-side failure with an error_id for support.</summary>
public sealed class InternalServerException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public InternalServerException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
