using System;

namespace Tesote.Sdk.Errors;

/// <summary>404 PAYMENT_METHOD_NOT_FOUND.</summary>
public sealed class PaymentMethodNotFoundException : NotFoundException
{
    /// <summary>Construct with the full required-field set.</summary>
    public PaymentMethodNotFoundException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
