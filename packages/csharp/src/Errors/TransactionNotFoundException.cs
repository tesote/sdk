using System;

namespace Tesote.Sdk.Errors;

/// <summary>404 TRANSACTION_NOT_FOUND.</summary>
public sealed class TransactionNotFoundException : NotFoundException
{
    /// <summary>Construct with the full required-field set.</summary>
    public TransactionNotFoundException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
