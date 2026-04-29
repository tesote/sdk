using System;
using System.Buffers;
using System.Collections.Generic;
using System.Globalization;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Net.Sockets;
using System.Security.Authentication;
using System.Text;
using System.Text.Json.Nodes;
using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Errors;
using Tesote.Sdk.Internal;

namespace Tesote.Sdk;

/// <summary>
/// Single HTTP entry point for every resource client.
///
/// Owns: bearer injection, retries with exponential backoff + jitter,
/// rate-limit header capture, idempotency-key generation, request-id
/// propagation into thrown exceptions, and the opt-in TTL response cache.
///
/// Resource clients call <see cref="RequestAsync"/> and receive a parsed
/// <see cref="JsonNode"/> (or <c>null</c> for empty bodies) or a typed
/// <see cref="TesoteException"/>.
/// </summary>
public sealed class Transport : IAsyncDisposable, IDisposable
{
    /// <summary>Default base URL used when none is supplied via options or env.</summary>
    public const string DefaultBaseUrl = "https://equipo.tesote.com/api";

    /// <summary>Current SDK version, exposed for User-Agent.</summary>
    public const string SdkVersion = "0.2.0";

    private static readonly HashSet<string> MutatingMethods =
        new(StringComparer.OrdinalIgnoreCase) { "POST", "PUT", "PATCH", "DELETE" };

    private readonly HttpClient _httpClient;
    private readonly bool _ownsHttpClient;
    private readonly string _apiKey;
    private readonly string _baseUrl;
    private readonly string _userAgent;
    private readonly RetryPolicy _retryPolicy;
    private readonly ICacheBackend? _cacheBackend;
    private readonly Action<LogEvent>? _logger;
    private RateLimitSnapshot _lastRateLimit = RateLimitSnapshot.Empty;
    private bool _disposed;

    /// <summary>Construct a transport from <see cref="ClientOptions"/>; resolves env-var fallbacks.</summary>
    public Transport(ClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(options);

        var resolvedKey = options.ApiKey ?? Environment.GetEnvironmentVariable("TESOTE_SDK_API_KEY");
        if (string.IsNullOrWhiteSpace(resolvedKey))
        {
            throw new ConfigException("apiKey is required");
        }
        _apiKey = resolvedKey;

        var resolvedBase = options.BaseUrl
            ?? Environment.GetEnvironmentVariable("TESOTE_SDK_API_URL")
            ?? DefaultBaseUrl;
        if (!Uri.TryCreate(resolvedBase, UriKind.Absolute, out _))
        {
            throw new ConfigException($"baseUrl is not a valid absolute URL: {resolvedBase}");
        }
        _baseUrl = TrimTrailingSlash(resolvedBase);

        _userAgent = options.UserAgent ?? DefaultUserAgent();
        _retryPolicy = options.RetryPolicy ?? RetryPolicy.Defaults;
        _cacheBackend = options.CacheBackend;
        _logger = options.Logger;

        if (options.HttpHandler is not null)
        {
            _httpClient = new HttpClient(options.HttpHandler, disposeHandler: false);
            _ownsHttpClient = true;
        }
        else
        {
            // why: SocketsHttpHandler reuses the connection pool per Transport instance,
            // matching the "one pool per Client" contract from transport.md.
            var handler = new SocketsHttpHandler
            {
                ConnectTimeout = TimeSpan.FromSeconds(5),
                PooledConnectionLifetime = TimeSpan.FromMinutes(5),
            };
            _httpClient = new HttpClient(handler, disposeHandler: true);
            _ownsHttpClient = true;
        }
        _httpClient.Timeout = options.RequestTimeout ?? TimeSpan.FromSeconds(30);
    }

    /// <summary>Last captured rate-limit snapshot. Empty until the first request lands.</summary>
    public RateLimitSnapshot LastRateLimit => Volatile.Read(ref _lastRateLimit);

    /// <summary>
    /// Send a request with retries, rate-limit awareness, and opt-in caching.
    /// Returns the parsed response, or <c>null</c> for empty bodies.
    /// </summary>
    public async Task<JsonNode?> RequestAsync(RequestOptions options, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(options);
        ThrowIfDisposed();

        var summary = RequestSummary.Create(
            options.Method, options.Path, options.Query,
            options.BodyShape, RedactBearer(_apiKey));

        var cacheable = string.Equals(options.Method, "GET", StringComparison.OrdinalIgnoreCase)
            && _cacheBackend is not null
            && options.CacheTtl is { Ticks: > 0 };

        string? cacheKey = cacheable ? BuildCacheKey(options) : null;
        if (cacheable && cacheKey is not null)
        {
            var hit = _cacheBackend!.Get(cacheKey);
            if (hit is not null)
            {
                return Json.Parse(hit);
            }
        }

        var idempotencyKey = options.IdempotencyKey;
        if (idempotencyKey is null && MutatingMethods.Contains(options.Method))
        {
            idempotencyKey = Guid.NewGuid().ToString();
        }

        Exception? lastTransport = null;
        ApiException? lastApi = null;

        for (var attempt = 1; attempt <= _retryPolicy.MaxAttempts; attempt++)
        {
            using var request = BuildRequest(options, idempotencyKey);
            HttpResponseMessage? response = null;
            byte[]? body = null;
            try
            {
                response = await _httpClient.SendAsync(request, HttpCompletionOption.ResponseHeadersRead, cancellationToken)
                    .ConfigureAwait(false);
                body = await response.Content.ReadAsByteArrayAsync(cancellationToken).ConfigureAwait(false);
            }
            catch (TaskCanceledException ex) when (!cancellationToken.IsCancellationRequested)
            {
                response?.Dispose();
                lastTransport = ex;
                _logger?.Invoke(new LogEvent(summary, attempt, -1, ex));
                if (!_retryPolicy.RetryOnNetwork || attempt == _retryPolicy.MaxAttempts)
                {
                    throw new TesoteTimeoutException("request timed out", summary, attempt, ex);
                }
                await Task.Delay(Backoff(attempt, null), cancellationToken).ConfigureAwait(false);
                continue;
            }
            catch (OperationCanceledException) when (cancellationToken.IsCancellationRequested)
            {
                response?.Dispose();
                throw;
            }
            catch (AuthenticationException ex)
            {
                response?.Dispose();
                throw new TlsException("TLS error: " + ex.Message, summary, attempt, ex);
            }
            catch (HttpRequestException ex) when (IsTlsError(ex))
            {
                response?.Dispose();
                throw new TlsException("TLS error: " + ex.Message, summary, attempt, ex);
            }
            catch (HttpRequestException ex)
            {
                response?.Dispose();
                lastTransport = ex;
                _logger?.Invoke(new LogEvent(summary, attempt, -1, ex));
                if (!_retryPolicy.RetryOnNetwork || attempt == _retryPolicy.MaxAttempts)
                {
                    throw new NetworkException(ex.Message, summary, attempt, ex);
                }
                await Task.Delay(Backoff(attempt, null), cancellationToken).ConfigureAwait(false);
                continue;
            }
            catch (IOException ex)
            {
                response?.Dispose();
                lastTransport = ex;
                _logger?.Invoke(new LogEvent(summary, attempt, -1, ex));
                if (!_retryPolicy.RetryOnNetwork || attempt == _retryPolicy.MaxAttempts)
                {
                    throw new NetworkException(ex.Message, summary, attempt, ex);
                }
                await Task.Delay(Backoff(attempt, null), cancellationToken).ConfigureAwait(false);
                continue;
            }
            catch (SocketException ex)
            {
                response?.Dispose();
                lastTransport = ex;
                _logger?.Invoke(new LogEvent(summary, attempt, -1, ex));
                if (!_retryPolicy.RetryOnNetwork || attempt == _retryPolicy.MaxAttempts)
                {
                    throw new NetworkException(ex.Message, summary, attempt, ex);
                }
                await Task.Delay(Backoff(attempt, null), cancellationToken).ConfigureAwait(false);
                continue;
            }

            using (response)
            {
                CaptureRateLimit(response);
                var status = (int)response.StatusCode;
                var requestId = FirstHeader(response, "X-Request-Id");

                _logger?.Invoke(new LogEvent(summary, attempt, status, null));

                if (status >= 200 && status < 300)
                {
                    var parsed = body is null || body.Length == 0
                        ? null
                        : Json.Parse(body);
                    if (cacheable && cacheKey is not null && _cacheBackend is not null)
                    {
                        _cacheBackend.Put(cacheKey, body ?? Array.Empty<byte>(), options.CacheTtl!.Value);
                    }
                    if (MutatingMethods.Contains(options.Method) && _cacheBackend is not null)
                    {
                        // why: any mutation invalidates GET caches under the same path prefix.
                        _cacheBackend.InvalidatePrefix("GET " + options.Path);
                    }
                    return parsed;
                }

                var api = BuildApiException(summary, response, body, requestId, attempt);
                lastApi = api;

                if (ShouldRetry(status, attempt))
                {
                    var sleepFor = Backoff(attempt, RetryAfterSeconds(response));
                    await Task.Delay(sleepFor, cancellationToken).ConfigureAwait(false);
                    continue;
                }
                throw api;
            }
        }

        // why: loop only exits via throw on terminal failure or return on success.
        if (lastApi is not null)
        {
            throw lastApi;
        }
        if (lastTransport is not null)
        {
            throw new NetworkException("retries exhausted", summary, _retryPolicy.MaxAttempts, lastTransport);
        }
        throw new NetworkException("unexpected transport state", summary, _retryPolicy.MaxAttempts, null);
    }

    /// <summary>
    /// Send a request and return the raw response body alongside content-type.
    /// Used for file-download endpoints (CSV / JSON export) where the body is not a JSON envelope.
    /// Bypasses caching; still applies retries, rate-limit awareness, idempotency, and error mapping.
    /// </summary>
    public async Task<RawResponse> RequestRawAsync(RequestOptions options, CancellationToken cancellationToken = default)
    {
        ArgumentNullException.ThrowIfNull(options);
        ThrowIfDisposed();

        var summary = RequestSummary.Create(
            options.Method, options.Path, options.Query,
            options.BodyShape, RedactBearer(_apiKey));

        var idempotencyKey = options.IdempotencyKey;
        if (idempotencyKey is null && MutatingMethods.Contains(options.Method))
        {
            idempotencyKey = Guid.NewGuid().ToString();
        }

        Exception? lastTransport = null;
        ApiException? lastApi = null;

        for (var attempt = 1; attempt <= _retryPolicy.MaxAttempts; attempt++)
        {
            using var request = BuildRequest(options, idempotencyKey);
            HttpResponseMessage? response = null;
            byte[]? body = null;
            try
            {
                response = await _httpClient.SendAsync(request, HttpCompletionOption.ResponseHeadersRead, cancellationToken)
                    .ConfigureAwait(false);
                body = await response.Content.ReadAsByteArrayAsync(cancellationToken).ConfigureAwait(false);
            }
            catch (TaskCanceledException ex) when (!cancellationToken.IsCancellationRequested)
            {
                response?.Dispose();
                lastTransport = ex;
                _logger?.Invoke(new LogEvent(summary, attempt, -1, ex));
                if (!_retryPolicy.RetryOnNetwork || attempt == _retryPolicy.MaxAttempts)
                {
                    throw new TesoteTimeoutException("request timed out", summary, attempt, ex);
                }
                await Task.Delay(Backoff(attempt, null), cancellationToken).ConfigureAwait(false);
                continue;
            }
            catch (OperationCanceledException) when (cancellationToken.IsCancellationRequested)
            {
                response?.Dispose();
                throw;
            }
            catch (HttpRequestException ex) when (IsTlsError(ex))
            {
                response?.Dispose();
                throw new TlsException("TLS error: " + ex.Message, summary, attempt, ex);
            }
            catch (HttpRequestException ex)
            {
                response?.Dispose();
                lastTransport = ex;
                _logger?.Invoke(new LogEvent(summary, attempt, -1, ex));
                if (!_retryPolicy.RetryOnNetwork || attempt == _retryPolicy.MaxAttempts)
                {
                    throw new NetworkException(ex.Message, summary, attempt, ex);
                }
                await Task.Delay(Backoff(attempt, null), cancellationToken).ConfigureAwait(false);
                continue;
            }
            catch (IOException ex)
            {
                response?.Dispose();
                lastTransport = ex;
                _logger?.Invoke(new LogEvent(summary, attempt, -1, ex));
                if (!_retryPolicy.RetryOnNetwork || attempt == _retryPolicy.MaxAttempts)
                {
                    throw new NetworkException(ex.Message, summary, attempt, ex);
                }
                await Task.Delay(Backoff(attempt, null), cancellationToken).ConfigureAwait(false);
                continue;
            }
            catch (SocketException ex)
            {
                response?.Dispose();
                lastTransport = ex;
                _logger?.Invoke(new LogEvent(summary, attempt, -1, ex));
                if (!_retryPolicy.RetryOnNetwork || attempt == _retryPolicy.MaxAttempts)
                {
                    throw new NetworkException(ex.Message, summary, attempt, ex);
                }
                await Task.Delay(Backoff(attempt, null), cancellationToken).ConfigureAwait(false);
                continue;
            }

            using (response)
            {
                CaptureRateLimit(response);
                var status = (int)response.StatusCode;
                var requestId = FirstHeader(response, "X-Request-Id");
                _logger?.Invoke(new LogEvent(summary, attempt, status, null));

                if (status >= 200 && status < 300)
                {
                    var contentType = response.Content.Headers.ContentType?.ToString() ?? "application/octet-stream";
                    var contentDisposition = response.Content.Headers.ContentDisposition?.ToString();
                    return new RawResponse(body ?? Array.Empty<byte>(), contentType, contentDisposition, requestId, status);
                }

                var api = BuildApiException(summary, response, body, requestId, attempt);
                lastApi = api;

                if (ShouldRetry(status, attempt))
                {
                    var sleepFor = Backoff(attempt, RetryAfterSeconds(response));
                    await Task.Delay(sleepFor, cancellationToken).ConfigureAwait(false);
                    continue;
                }
                throw api;
            }
        }

        if (lastApi is not null)
        {
            throw lastApi;
        }
        if (lastTransport is not null)
        {
            throw new NetworkException("retries exhausted", summary, _retryPolicy.MaxAttempts, lastTransport);
        }
        throw new NetworkException("unexpected transport state", summary, _retryPolicy.MaxAttempts, null);
    }

    private HttpRequestMessage BuildRequest(RequestOptions options, string? idempotencyKey)
    {
        var uri = BuildUri(options);
        var method = new HttpMethod(options.Method);
        var request = new HttpRequestMessage(method, uri);
        request.Headers.TryAddWithoutValidation("Authorization", "Bearer " + _apiKey);
        request.Headers.TryAddWithoutValidation("Accept", "application/json");
        request.Headers.TryAddWithoutValidation("User-Agent", _userAgent);

        if (idempotencyKey is not null)
        {
            request.Headers.TryAddWithoutValidation("Idempotency-Key", idempotencyKey);
        }
        if (options.ExtraHeaders is not null)
        {
            foreach (var kv in options.ExtraHeaders)
            {
                request.Headers.TryAddWithoutValidation(kv.Key, kv.Value);
            }
        }

        if (options.Body is not null)
        {
            var content = new ByteArrayContent(options.Body);
            content.Headers.TryAddWithoutValidation("Content-Type", "application/json");
            request.Content = content;
        }
        return request;
    }

    private Uri BuildUri(RequestOptions options)
    {
        var sb = new StringBuilder(_baseUrl);
        if (!options.Path.StartsWith('/'))
        {
            sb.Append('/');
        }
        sb.Append(options.Path);
        if (options.Query is { Count: > 0 })
        {
            sb.Append('?');
            var first = true;
            foreach (var kv in options.Query)
            {
                if (!first)
                {
                    sb.Append('&');
                }
                first = false;
                sb.Append(WebUtility.UrlEncode(kv.Key));
                sb.Append('=');
                sb.Append(WebUtility.UrlEncode(kv.Value));
            }
        }
        return new Uri(sb.ToString(), UriKind.Absolute);
    }

    private bool ShouldRetry(int status, int attempt)
    {
        if (attempt >= _retryPolicy.MaxAttempts)
        {
            return false;
        }
        return status is 429 or 502 or 503 or 504;
    }

    private TimeSpan Backoff(int attempt, int? retryAfterSeconds)
    {
        if (retryAfterSeconds is > 0)
        {
            return TimeSpan.FromSeconds(retryAfterSeconds.Value);
        }
        var baseMs = _retryPolicy.BaseDelay.TotalMilliseconds;
        var capMs = _retryPolicy.MaxDelay.TotalMilliseconds;
        var exp = baseMs * Math.Pow(2, attempt - 1);
        var capped = Math.Min(capMs, exp);
        var jitter = Random.Shared.NextDouble() * (capped / 4);
        return TimeSpan.FromMilliseconds(capped + jitter);
    }

    private void CaptureRateLimit(HttpResponseMessage response)
    {
        var limit = ParseInt(FirstHeader(response, "X-RateLimit-Limit"), -1);
        var remaining = ParseInt(FirstHeader(response, "X-RateLimit-Remaining"), -1);
        var resetAt = ParseResetAt(FirstHeader(response, "X-RateLimit-Reset"));
        Volatile.Write(ref _lastRateLimit, new RateLimitSnapshot(limit, remaining, resetAt));
    }

    private static DateTimeOffset? ParseResetAt(string? raw)
    {
        if (string.IsNullOrWhiteSpace(raw))
        {
            return null;
        }
        var trimmed = raw.Trim();
        if (long.TryParse(trimmed, NumberStyles.Integer, CultureInfo.InvariantCulture, out var epoch))
        {
            return DateTimeOffset.FromUnixTimeSeconds(epoch);
        }
        if (DateTimeOffset.TryParse(trimmed, CultureInfo.InvariantCulture, DateTimeStyles.AssumeUniversal, out var iso))
        {
            return iso;
        }
        return null;
    }

    private static int? RetryAfterSeconds(HttpResponseMessage response)
    {
        var raw = FirstHeader(response, "Retry-After");
        if (string.IsNullOrWhiteSpace(raw))
        {
            return null;
        }
        return int.TryParse(raw.Trim(), NumberStyles.Integer, CultureInfo.InvariantCulture, out var v) ? v : null;
    }

    private static int ParseInt(string? raw, int dflt)
    {
        if (string.IsNullOrWhiteSpace(raw))
        {
            return dflt;
        }
        return int.TryParse(raw.Trim(), NumberStyles.Integer, CultureInfo.InvariantCulture, out var v) ? v : dflt;
    }

    private static string? FirstHeader(HttpResponseMessage r, string name)
    {
        if (r.Headers.TryGetValues(name, out var values))
        {
            foreach (var v in values)
            {
                return v;
            }
        }
        if (r.Content.Headers.TryGetValues(name, out var cvalues))
        {
            foreach (var v in cvalues)
            {
                return v;
            }
        }
        return null;
    }

    private static ApiException BuildApiException(
        RequestSummary summary,
        HttpResponseMessage response,
        byte[]? body,
        string? requestId,
        int attempt)
    {
        var bytes = body ?? Array.Empty<byte>();
        var bodyStr = bytes.Length == 0 ? string.Empty : Encoding.UTF8.GetString(bytes);
        var envelope = Json.Parse(bytes) as JsonObject;

        var message = StringField(envelope, "error") ?? "HTTP " + (int)response.StatusCode;
        var errorCode = StringField(envelope, "error_code");
        var errorId = StringField(envelope, "error_id");
        int? retryAfter = RetryAfterSeconds(response);
        if (retryAfter is null && envelope is not null && envelope.TryGetPropertyValue("retry_after", out var ra) && ra is not null)
        {
            if (ra.GetValueKind() == System.Text.Json.JsonValueKind.Number)
            {
                retryAfter = ra.GetValue<int>();
            }
        }

        return ErrorDispatcher.Dispatch(
            message, errorCode, (int)response.StatusCode, requestId, errorId,
            retryAfter, bodyStr, summary, attempt, null);
    }

    private static string? StringField(JsonObject? obj, string key)
    {
        if (obj is null)
        {
            return null;
        }
        if (!obj.TryGetPropertyValue(key, out var node) || node is null)
        {
            return null;
        }
        if (node.GetValueKind() == System.Text.Json.JsonValueKind.String)
        {
            return node.GetValue<string>();
        }
        return node.ToJsonString();
    }

    private string BuildCacheKey(RequestOptions opts)
    {
        var sb = new StringBuilder("GET ").Append(opts.Path);
        if (opts.Query is { Count: > 0 })
        {
            sb.Append('?');
            var sorted = new SortedDictionary<string, string>(new Dictionary<string, string>(opts.Query), StringComparer.Ordinal);
            foreach (var kv in sorted)
            {
                sb.Append(kv.Key).Append('=').Append(kv.Value).Append('&');
            }
        }
        sb.Append("|key=").Append(ApiKeyHash());
        return sb.ToString();
    }

    private string ApiKeyHash()
    {
        // why: scope cache by API-key hash to avoid cross-tenant bleed without
        // logging the bearer token itself.
        return _apiKey.GetHashCode(StringComparison.Ordinal).ToString("x", CultureInfo.InvariantCulture);
    }

    /// <summary>Redact a bearer token to <c>Bearer ****&lt;last4&gt;</c>.</summary>
    public static string RedactBearer(string? key)
    {
        if (string.IsNullOrEmpty(key) || key.Length < 4)
        {
            return "Bearer ****";
        }
        return "Bearer ****" + key.Substring(key.Length - 4);
    }

    private static string TrimTrailingSlash(string s)
    {
        return s.EndsWith('/') ? s[..^1] : s;
    }

    private static string DefaultUserAgent()
    {
        var runtime = Environment.Version.ToString();
        return $"tesote-sdk-csharp/{SdkVersion} (dotnet/{runtime})";
    }

    private static bool IsTlsError(HttpRequestException ex)
    {
        return ex.InnerException is AuthenticationException
            || (ex.Message?.Contains("SSL", StringComparison.OrdinalIgnoreCase) ?? false)
            || (ex.Message?.Contains("TLS", StringComparison.OrdinalIgnoreCase) ?? false);
    }

    private void ThrowIfDisposed()
    {
        if (_disposed)
        {
            throw new ConfigException("Transport has been disposed");
        }
    }

    /// <summary>Async dispose. Releases the pooled <see cref="HttpClient"/>.</summary>
    public ValueTask DisposeAsync()
    {
        Dispose();
        return ValueTask.CompletedTask;
    }

    /// <summary>Synchronous dispose. Releases the pooled <see cref="HttpClient"/>.</summary>
    public void Dispose()
    {
        if (_disposed)
        {
            return;
        }
        _disposed = true;
        if (_ownsHttpClient)
        {
            _httpClient.Dispose();
        }
    }

    // why: convenience for tests that want to inspect the mutating-methods set.
    internal static IReadOnlyCollection<string> MutatingMethodsView => MutatingMethods;
}
