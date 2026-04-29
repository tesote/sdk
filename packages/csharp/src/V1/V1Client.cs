using System;
using System.Threading.Tasks;

namespace Tesote.Sdk.V1;

/// <summary>
/// v1 client. Read-only foundation: status, accounts, transactions.
/// </summary>
public sealed class V1Client : IAsyncDisposable, IDisposable
{
    /// <summary>The transport instance, exposed for advanced users.</summary>
    public Transport Transport { get; }

    /// <summary>Status + whoami resource client.</summary>
    public StatusClient Status { get; }

    /// <summary>Accounts resource client.</summary>
    public AccountsClient Accounts { get; }

    /// <summary>Transactions resource client.</summary>
    public TransactionsClient Transactions { get; }

    /// <summary>Construct a v1 client from <see cref="ClientOptions"/>.</summary>
    public V1Client(ClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(options);
        Transport = new Transport(options);
        Status = new StatusClient(Transport);
        Accounts = new AccountsClient(Transport);
        Transactions = new TransactionsClient(Transport);
    }

    /// <summary>Last captured rate-limit snapshot.</summary>
    public RateLimitSnapshot LastRateLimit => Transport.LastRateLimit;

    /// <summary>Async dispose.</summary>
    public ValueTask DisposeAsync()
    {
        Dispose();
        return ValueTask.CompletedTask;
    }

    /// <summary>Dispose the underlying <see cref="Transport"/>.</summary>
    public void Dispose() => Transport.Dispose();
}
