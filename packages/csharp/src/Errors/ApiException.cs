using System;

namespace Tesote.Sdk.Errors;

/// <summary>
/// Server returned a usable HTTP response with an error envelope.
/// Concrete subclasses correspond 1:1 to API <c>error_code</c> values.
/// </summary>
public class ApiException : TesoteException
{
    /// <summary>Construct with the full required-field set.</summary>
    public ApiException(
        string? message,
        string? errorCode,
        int httpStatus,
        string? requestId,
        string? errorId,
        int? retryAfter,
        string? responseBody,
        RequestSummary? requestSummary,
        int attempts,
        Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
