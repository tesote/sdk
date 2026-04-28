using System;

namespace Tesote.Sdk.Errors;

/// <summary>409 MUTATION_CONFLICT — dataset changed mid-iteration.</summary>
public sealed class MutationDuringPaginationException : ApiException
{
    /// <summary>Construct with the full required-field set.</summary>
    public MutationDuringPaginationException(
        string? message, string? errorCode, int httpStatus,
        string? requestId, string? errorId, int? retryAfter,
        string? responseBody, RequestSummary? requestSummary,
        int attempts, Exception? cause)
        : base(message, errorCode, httpStatus, requestId, errorId,
               retryAfter, responseBody, requestSummary, attempts, cause)
    {
    }
}
