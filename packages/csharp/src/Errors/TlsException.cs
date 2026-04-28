using System;

namespace Tesote.Sdk.Errors;

/// <summary>Certificate or TLS handshake failure.</summary>
public sealed class TlsException : TransportException
{
    /// <summary>Construct from a <see cref="RequestSummary"/> and underlying cause.</summary>
    public TlsException(string? message, RequestSummary? requestSummary, int attempts, Exception? cause)
        : base(message, "TLS_ERROR", 0, null, null, null, null, requestSummary, attempts, cause)
    {
    }
}
