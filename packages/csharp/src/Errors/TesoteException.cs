using System;
using System.Text;

namespace Tesote.Sdk.Errors;

/// <summary>
/// Root of the SDK exception hierarchy.
///
/// Underlying causes are preserved through <see cref="Exception.InnerException"/>;
/// never lose the chain.
/// </summary>
public class TesoteException : Exception
{
    /// <summary>API <c>error_code</c> value, or a synthetic code for transport errors.</summary>
    public string? ErrorCode { get; }

    /// <summary>HTTP status from the failed response, or 0 for transport errors.</summary>
    public int HttpStatus { get; }

    /// <summary>Value of the <c>X-Request-Id</c> response header.</summary>
    public string? RequestId { get; }

    /// <summary>Value of the envelope <c>error_id</c> field.</summary>
    public string? ErrorId { get; }

    /// <summary>Suggested backoff seconds, from <c>Retry-After</c> or envelope.</summary>
    public int? RetryAfter { get; }

    /// <summary>Raw response body bytes decoded as UTF-8.</summary>
    public string? ResponseBody { get; }

    /// <summary>Redacted summary of the originating request (no bearer token).</summary>
    public RequestSummary? RequestSummary { get; }

    /// <summary>Number of attempts made before this exception was raised.</summary>
    public int Attempts { get; }

    /// <summary>Construct with the full required-field set.</summary>
    public TesoteException(
        string? message,
        string? errorCode,
        int httpStatus,
        string? requestId,
        string? errorId,
        int? retryAfter,
        string? responseBody,
        RequestSummary? requestSummary,
        int attempts,
        Exception? cause)
        : base(message, cause)
    {
        ErrorCode = errorCode;
        HttpStatus = httpStatus;
        RequestId = requestId;
        ErrorId = errorId;
        RetryAfter = retryAfter;
        ResponseBody = responseBody;
        RequestSummary = requestSummary;
        Attempts = attempts;
    }

    /// <summary>Greppable single-error multi-line string with all required fields.</summary>
    public override string ToString()
    {
        var sb = new StringBuilder();
        sb.Append(GetType().Name).Append(": ");
        if (HttpStatus > 0)
        {
            sb.Append(HttpStatus).Append(' ');
        }
        sb.Append(base.Message ?? string.Empty);
        if (ErrorCode is not null)
        {
            sb.Append("\n  error_code: ").Append(ErrorCode);
        }
        if (RequestId is not null)
        {
            sb.Append("\n  request_id: ").Append(RequestId);
        }
        if (ErrorId is not null)
        {
            sb.Append("\n  error_id: ").Append(ErrorId);
        }
        if (RetryAfter is not null)
        {
            sb.Append("\n  retry_after: ").Append(RetryAfter.Value).Append('s');
        }
        if (Attempts > 0)
        {
            sb.Append("\n  attempts: ").Append(Attempts);
        }
        if (RequestSummary is not null)
        {
            sb.Append("\n  request: ").Append(RequestSummary.Method)
              .Append(' ').Append(RequestSummary.Path);
            if (RequestSummary.BodyShape is not null)
            {
                sb.Append(" (body: ").Append(RequestSummary.BodyShape).Append(')');
            }
        }
        if (!string.IsNullOrEmpty(ResponseBody))
        {
            sb.Append("\n  response: ").Append(ResponseBody);
        }
        return sb.ToString();
    }
}
