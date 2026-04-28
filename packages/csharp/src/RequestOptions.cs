using System;
using System.Collections.Generic;

namespace Tesote.Sdk;

/// <summary>Per-request options. Resource clients construct one for each call.</summary>
public sealed class RequestOptions
{
    /// <summary>HTTP method (default <c>GET</c>).</summary>
    public string Method { get; set; } = "GET";

    /// <summary>Request path appended to the base URL, e.g. <c>/v2/accounts</c>.</summary>
    public string Path { get; set; } = "/";

    /// <summary>Optional query parameters; null means none.</summary>
    public IReadOnlyDictionary<string, string>? Query { get; set; }

    /// <summary>Optional UTF-8 request body. Triggers <c>Content-Type: application/json</c>.</summary>
    public byte[]? Body { get; set; }

    /// <summary>Optional shape descriptor for error summaries (no PII).</summary>
    public string? BodyShape { get; set; }

    /// <summary>Caller-supplied idempotency key; auto-generated for mutations when null.</summary>
    public string? IdempotencyKey { get; set; }

    /// <summary>Opt-in TTL for the response cache. Null disables caching.</summary>
    public TimeSpan? CacheTtl { get; set; }

    /// <summary>Optional extra headers merged into the request.</summary>
    public IReadOnlyDictionary<string, string>? ExtraHeaders { get; set; }

    /// <summary>Convenience factory for a GET request to the given path.</summary>
    public static RequestOptions Get(string path)
    {
        ArgumentNullException.ThrowIfNull(path);
        return new RequestOptions { Method = "GET", Path = path };
    }
}
