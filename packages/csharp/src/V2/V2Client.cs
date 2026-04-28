using System;
using System.Threading.Tasks;

namespace Tesote.Sdk.V2;

/// <summary>
/// v2 client. Adds writes for payments + sync orchestration on top of v1.
///
/// 0.1.0 ships transport plumbing plus <see cref="Accounts"/> list/get;
/// other resources stub with <see cref="NotImplementedException"/> until wired.
/// </summary>
public sealed class V2Client : IAsyncDisposable, IDisposable
{
    /// <summary>The transport instance, exposed for advanced users.</summary>
    public Transport Transport { get; }

    /// <summary>Accounts resource client.</summary>
    public AccountsClient Accounts { get; }

    /// <summary>Construct a v2 client from <see cref="ClientOptions"/>.</summary>
    public V2Client(ClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(options);
        Transport = new Transport(options);
        Accounts = new AccountsClient(Transport);
    }

    /// <summary>Status endpoint — not implemented in 0.1.0.</summary>
    public Task StatusAsync() => throw new NotImplementedException("not implemented");

    /// <summary>Transactions endpoint — not implemented in 0.1.0.</summary>
    public Task TransactionsAsync() => throw new NotImplementedException("not implemented");

    /// <summary>Sync sessions endpoint — not implemented in 0.1.0.</summary>
    public Task SyncSessionsAsync() => throw new NotImplementedException("not implemented");

    /// <summary>Transaction orders endpoint — not implemented in 0.1.0.</summary>
    public Task TransactionOrdersAsync() => throw new NotImplementedException("not implemented");

    /// <summary>Batches endpoint — not implemented in 0.1.0.</summary>
    public Task BatchesAsync() => throw new NotImplementedException("not implemented");

    /// <summary>Payment methods endpoint — not implemented in 0.1.0.</summary>
    public Task PaymentMethodsAsync() => throw new NotImplementedException("not implemented");

    /// <summary>Async dispose.</summary>
    public ValueTask DisposeAsync()
    {
        Dispose();
        return ValueTask.CompletedTask;
    }

    /// <summary>Dispose the underlying <see cref="Transport"/>.</summary>
    public void Dispose() => Transport.Dispose();
}
