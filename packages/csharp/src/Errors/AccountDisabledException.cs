using System;

namespace Tesote.Sdk.Errors;

/// <summary>403 ACCOUNT_DISABLED — caller's account is disabled.</summary>
public sealed class AccountDisabledException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public AccountDisabledException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
