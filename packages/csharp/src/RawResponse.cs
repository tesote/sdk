namespace Tesote.Sdk;

/// <summary>
/// Non-JSON HTTP response payload, returned from <see cref="Transport.RequestRawAsync"/>.
/// Exposes the raw bytes plus the headers callers need to interpret them.
/// </summary>
/// <param name="Body">Raw response bytes (may be empty, never null).</param>
/// <param name="ContentType">Value of <c>Content-Type</c>, or <c>application/octet-stream</c>.</param>
/// <param name="ContentDisposition">Value of <c>Content-Disposition</c> (filename hints).</param>
/// <param name="RequestId">Value of <c>X-Request-Id</c>.</param>
/// <param name="HttpStatus">HTTP status code (always 2xx since errors throw).</param>
public sealed record RawResponse(
    byte[] Body,
    string ContentType,
    string? ContentDisposition,
    string? RequestId,
    int HttpStatus);
