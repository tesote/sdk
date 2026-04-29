using System;
using System.Text;

namespace Tesote.Sdk.Internal;

/// <summary>Internal helpers for resource clients to construct <see cref="RequestOptions"/>.</summary>
internal static class Requests
{
    /// <summary>Build a JSON-encoded mutating request with idempotency-key support.</summary>
    public static RequestOptions Json(string method, string path, object? body, string? idempotencyKey = null)
    {
        var opts = new RequestOptions
        {
            Method = method,
            Path = path,
            IdempotencyKey = idempotencyKey,
        };

        var requiresBody = method.Equals("POST", StringComparison.OrdinalIgnoreCase)
            || method.Equals("PUT", StringComparison.OrdinalIgnoreCase)
            || method.Equals("PATCH", StringComparison.OrdinalIgnoreCase);

        if (body is not null)
        {
            opts.Body = Internal.Json.SerializeToUtf8Bytes(body);
            opts.BodyShape = opts.Body.Length + " bytes";
        }
        else if (requiresBody)
        {
            // why: server requires Content-Type: application/json on every POST/PUT/PATCH.
            opts.Body = Encoding.UTF8.GetBytes("{}");
            opts.BodyShape = "0 bytes";
        }
        return opts;
    }
}
