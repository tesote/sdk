using System;

namespace Tesote.Sdk.Errors;

/// <summary>503 BANK_UNDER_MAINTENANCE — upstream bank is unavailable.</summary>
public sealed class BankUnderMaintenanceException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public BankUnderMaintenanceException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
