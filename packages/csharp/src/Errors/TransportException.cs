using System;

namespace Tesote.Sdk.Errors;

/// <summary>
/// Lower-level failure with no usable HTTP response (DNS, TLS, timeout, reset).
/// </summary>
public class TransportException : TesoteException
{
    /// <summary>Construct with the full required-field set.</summary>
    public TransportException(
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
