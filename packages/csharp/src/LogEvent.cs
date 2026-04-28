using System;
using Tesote.Sdk.Errors;

namespace Tesote.Sdk;

/// <summary>Single-event log payload. <see cref="Error"/> is null on success.</summary>
/// <param name="Request">Redacted request summary.</param>
/// <param name="Attempt">1-based attempt index.</param>
/// <param name="Status">HTTP status received, or -1 for transport errors.</param>
/// <param name="Error">Underlying exception when a transport error occurred.</param>
public sealed record LogEvent(RequestSummary Request, int Attempt, int Status, Exception? Error);
