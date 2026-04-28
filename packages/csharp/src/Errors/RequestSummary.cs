using System.Collections.Generic;
using System.Collections.ObjectModel;

namespace Tesote.Sdk.Errors;

/// <summary>
/// Redacted snapshot of an outbound request, safe to include in error output.
/// Bearer tokens are redacted to <c>Bearer ****&lt;last4&gt;</c> before being captured here;
/// the raw token must never appear in a <see cref="RequestSummary"/>.
/// </summary>
/// <param name="Method">HTTP method (GET, POST, ...).</param>
/// <param name="Path">Request path, e.g. <c>/v2/accounts</c>.</param>
/// <param name="Query">Query parameters, defensively copied.</param>
/// <param name="BodyShape">Optional shape descriptor (e.g. <c>"47 items"</c>).</param>
/// <param name="RedactedAuthorization">Redacted bearer token, never the raw key.</param>
public sealed record RequestSummary(
    string Method,
    string Path,
    IReadOnlyDictionary<string, string> Query,
    string? BodyShape,
    string RedactedAuthorization)
{
    /// <summary>Construct with a defensively-copied query map; null becomes empty.</summary>
    public static RequestSummary Create(
        string method,
        string path,
        IReadOnlyDictionary<string, string>? query,
        string? bodyShape,
        string redactedAuthorization)
    {
        var copy = query is null
            ? new ReadOnlyDictionary<string, string>(new Dictionary<string, string>())
            : new ReadOnlyDictionary<string, string>(new Dictionary<string, string>(query));
        return new RequestSummary(method, path, copy, bodyShape, redactedAuthorization);
    }
}
