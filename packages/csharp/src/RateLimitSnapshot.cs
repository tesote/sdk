using System;

namespace Tesote.Sdk;

/// <summary>
/// Snapshot of the most recent <c>X-RateLimit-*</c> headers seen by the transport.
/// </summary>
/// <param name="Limit">Value of <c>X-RateLimit-Limit</c>, or -1 if absent.</param>
/// <param name="Remaining">Value of <c>X-RateLimit-Remaining</c>, or -1 if absent.</param>
/// <param name="ResetAt">Parsed <c>X-RateLimit-Reset</c>; null when absent or unparseable.</param>
public sealed record RateLimitSnapshot(int Limit, int Remaining, DateTimeOffset? ResetAt)
{
    /// <summary>An empty snapshot used before the first request lands.</summary>
    public static RateLimitSnapshot Empty { get; } = new(-1, -1, null);
}
