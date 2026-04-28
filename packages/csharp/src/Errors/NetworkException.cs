using System;

namespace Tesote.Sdk.Errors;

/// <summary>DNS, connection refused, reset — no HTTP response received.</summary>
public sealed class NetworkException : TransportException
{
    /// <summary>Construct from a <see cref="RequestSummary"/> and underlying cause.</summary>
    public NetworkException(string? message, RequestSummary? requestSummary, int attempts, Exception? cause)
        : base(message, "NETWORK_ERROR", 0, null, null, null, null, requestSummary, attempts, cause)
    {
    }
}
