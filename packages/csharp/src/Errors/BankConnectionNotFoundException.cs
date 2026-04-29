using System;

namespace Tesote.Sdk.Errors;

/// <summary>404 BANK_CONNECTION_NOT_FOUND.</summary>
public sealed class BankConnectionNotFoundException : NotFoundException
{
    /// <summary>Construct with the full required-field set.</summary>
    public BankConnectionNotFoundException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
