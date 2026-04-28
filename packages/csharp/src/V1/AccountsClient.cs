using System;
using System.Collections.Generic;
using System.Text.Json.Nodes;
using System.Threading;
using System.Threading.Tasks;

namespace Tesote.Sdk.V1;

/// <summary>
/// v1 Accounts resource. Read-only listing and lookup.
/// </summary>
public sealed class AccountsClient
{
    private const string BasePath = "/v1/accounts";

    private readonly Transport _transport;

    internal AccountsClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>List accounts. Returns the raw response envelope until typed models land.</summary>
    public Task<JsonNode?> ListAsync(IReadOnlyDictionary<string, string>? query = null, CancellationToken cancellationToken = default)
    {
        var opts = RequestOptions.Get(BasePath);
        opts.Query = query;
        return _transport.RequestAsync(opts, cancellationToken);
    }

    /// <summary>Fetch a single account by id.</summary>
    public Task<JsonNode?> GetAsync(string accountId, CancellationToken cancellationToken = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        var opts = RequestOptions.Get(BasePath + "/" + Uri.EscapeDataString(accountId));
        return _transport.RequestAsync(opts, cancellationToken);
    }
}
