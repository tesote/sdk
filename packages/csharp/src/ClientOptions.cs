using System;
using System.Net.Http;

namespace Tesote.Sdk;

/// <summary>
/// Configuration for <see cref="V1.V1Client"/> and <see cref="V2.V2Client"/>.
///
/// Resolution order at construction: explicit property → environment variable → default.
/// </summary>
public sealed class ClientOptions
{
    /// <summary>Bearer API key. Falls back to <c>TESOTE_SDK_API_KEY</c> env var.</summary>
    public string? ApiKey { get; set; }

    /// <summary>Override base URL. Falls back to <c>TESOTE_SDK_API_URL</c>, then <see cref="Transport.DefaultBaseUrl"/>.</summary>
    public string? BaseUrl { get; set; }

    /// <summary>Override User-Agent string. Defaults to <c>tesote-sdk-csharp/&lt;version&gt; (dotnet/&lt;runtime&gt;)</c>.</summary>
    public string? UserAgent { get; set; }

    /// <summary>Per-request timeout. Default 30s.</summary>
    public TimeSpan? RequestTimeout { get; set; }

    /// <summary>Retry policy. Default <see cref="RetryPolicy.Defaults"/>.</summary>
    public RetryPolicy? RetryPolicy { get; set; }

    /// <summary>Optional cache backend; required for the opt-in TTL response cache to function.</summary>
    public ICacheBackend? CacheBackend { get; set; }

    /// <summary>Optional logger callback invoked once per request attempt.</summary>
    public Action<LogEvent>? Logger { get; set; }

    /// <summary>Optional <see cref="HttpMessageHandler"/> override (e.g. for testing).</summary>
    public HttpMessageHandler? HttpHandler { get; set; }
}
