using System;
using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V1;

/// <summary>v1 Transactions resource. Lookup-by-id only; list lives on the account client.</summary>
public sealed class TransactionsClient
{
    private const string BasePath = "/v1/transactions";

    private readonly Transport _transport;

    internal TransactionsClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>Fetch a single transaction by id.</summary>
    public async Task<Transaction> GetAsync(string transactionId, CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(transactionId);
        var opts = RequestOptions.Get(BasePath + "/" + Uri.EscapeDataString(transactionId));
        opts.CacheTtl = TimeSpan.FromMinutes(5);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<Transaction>(node);
    }
}
