using System;

namespace Tesote.Sdk;

/// <summary>Configurable retry parameters for the transport layer.</summary>
/// <param name="MaxAttempts">Total attempts including the first try. Default 3.</param>
/// <param name="BaseDelay">Base delay for exponential backoff. Default 250ms.</param>
/// <param name="MaxDelay">Cap on backoff delay. Default 8s.</param>
/// <param name="RetryOnNetwork">Whether to retry transport-level network failures.</param>
public sealed record RetryPolicy(int MaxAttempts, TimeSpan BaseDelay, TimeSpan MaxDelay, bool RetryOnNetwork)
{
    /// <summary>Documented defaults: 3 attempts, 250ms base, 8s cap, retry on network.</summary>
    public static RetryPolicy Defaults { get; } =
        new(3, TimeSpan.FromMilliseconds(250), TimeSpan.FromSeconds(8), true);
}
