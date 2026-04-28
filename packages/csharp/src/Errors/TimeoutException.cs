using System;

namespace Tesote.Sdk.Errors;

/// <summary>Connect or read timeout. Use <c>TesoteTimeoutException</c> in code to disambiguate from <see cref="System.TimeoutException"/>.</summary>
public sealed class TesoteTimeoutException : TransportException
{
    /// <summary>Construct from a <see cref="RequestSummary"/> and underlying cause.</summary>
    public TesoteTimeoutException(string? message, RequestSummary? requestSummary, int attempts, Exception? cause)
        : base(message, "TIMEOUT", 0, null, null, null, null, requestSummary, attempts, cause)
    {
    }
}
