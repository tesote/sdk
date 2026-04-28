using System.Collections.Generic;
using Tesote.Sdk;
using Tesote.Sdk.Errors;
using Xunit;

namespace Tesote.Sdk.Tests;

public sealed class ErrorsTests
{
    private static RequestSummary Summary() =>
        RequestSummary.Create("GET", "/v3/accounts", new Dictionary<string, string>(), null, "Bearer ****1234");

    private static ApiException Dispatch(string code, int status) =>
        ErrorDispatcher.Dispatch("msg", code, status, "req_1", "err_1", null, "{}", Summary(), 1, null);

    [Fact]
    public void UnauthorizedMaps() =>
        Assert.IsType<UnauthorizedException>(Dispatch("UNAUTHORIZED", 401));

    [Fact]
    public void ApiKeyRevokedMaps() =>
        Assert.IsType<ApiKeyRevokedException>(Dispatch("API_KEY_REVOKED", 401));

    [Fact]
    public void WorkspaceSuspendedMaps() =>
        Assert.IsType<WorkspaceSuspendedException>(Dispatch("WORKSPACE_SUSPENDED", 403));

    [Fact]
    public void AccountDisabledMaps() =>
        Assert.IsType<AccountDisabledException>(Dispatch("ACCOUNT_DISABLED", 403));

    [Fact]
    public void HistorySyncForbiddenMaps() =>
        Assert.IsType<HistorySyncForbiddenException>(Dispatch("HISTORY_SYNC_FORBIDDEN", 403));

    [Fact]
    public void MutationConflictMaps() =>
        Assert.IsType<MutationDuringPaginationException>(Dispatch("MUTATION_CONFLICT", 409));

    [Fact]
    public void UnprocessableContentMaps() =>
        Assert.IsType<UnprocessableContentException>(Dispatch("UNPROCESSABLE_CONTENT", 422));

    [Fact]
    public void InvalidDateRangeMaps()
    {
        var ex = Dispatch("INVALID_DATE_RANGE", 422);
        Assert.IsType<InvalidDateRangeException>(ex);
        // why: subclass of UnprocessableContentException so callers can catch the parent.
        Assert.IsAssignableFrom<UnprocessableContentException>(ex);
    }

    [Fact]
    public void RateLimitExceededMaps() =>
        Assert.IsType<RateLimitExceededException>(Dispatch("RATE_LIMIT_EXCEEDED", 429));

    [Fact]
    public void ServiceUnavailableMapsByStatusWithEmptyCode()
    {
        var ex = Dispatch("", 503);
        Assert.IsType<ServiceUnavailableException>(ex);
    }

    [Fact]
    public void UnknownCodeFallsBackToApiException()
    {
        var ex = Dispatch("MYSTERY_CODE", 418);
        Assert.Equal(typeof(ApiException), ex.GetType());
        Assert.Equal("MYSTERY_CODE", ex.ErrorCode);
    }

    [Fact]
    public void RequiredFieldsPopulated()
    {
        var ex = Dispatch("UNAUTHORIZED", 401);
        Assert.Equal(401, ex.HttpStatus);
        Assert.Equal("UNAUTHORIZED", ex.ErrorCode);
        Assert.Equal("req_1", ex.RequestId);
        Assert.Equal("err_1", ex.ErrorId);
        Assert.Equal(1, ex.Attempts);
        Assert.NotNull(ex.RequestSummary);
    }

    [Fact]
    public void TransportErrorsExtendTransportException()
    {
        var net = new NetworkException("boom", Summary(), 1, null);
        var to = new TesoteTimeoutException("slow", Summary(), 1, null);
        var tls = new TlsException("cert", Summary(), 1, null);
        Assert.IsAssignableFrom<TransportException>(net);
        Assert.IsAssignableFrom<TransportException>(to);
        Assert.IsAssignableFrom<TransportException>(tls);
        Assert.Equal("NETWORK_ERROR", net.ErrorCode);
        Assert.Equal("TIMEOUT", to.ErrorCode);
        Assert.Equal("TLS_ERROR", tls.ErrorCode);
    }

    [Fact]
    public void BearerRedactionUtility()
    {
        var redacted = Transport.RedactBearer("sk_test_abcd1234");
        Assert.StartsWith("Bearer ****", redacted);
        Assert.EndsWith("1234", redacted);
    }

    [Fact]
    public void ShortKeyStillRedacted()
    {
        var redacted = Transport.RedactBearer("ab");
        Assert.Equal("Bearer ****", redacted);
    }

    [Fact]
    public void RequestSummaryNeverContainsRawKey()
    {
        // why: bearer redaction must hold regardless of how RequestSummary is constructed.
        const string apiKey = "sk_live_supersecret_abcd1234";
        var redacted = Transport.RedactBearer(apiKey);
        var summary = RequestSummary.Create("POST", "/v2/accounts", null, "1 item", redacted);
        Assert.DoesNotContain("supersecret", summary.RedactedAuthorization);
    }
}
