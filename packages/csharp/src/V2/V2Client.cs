using System;
using System.Threading.Tasks;

namespace Tesote.Sdk.V2;

/// <summary>
/// v2 client. Adds writes for payments + sync orchestration on top of v1.
/// </summary>
public sealed class V2Client : IAsyncDisposable, IDisposable
{
    /// <summary>The transport instance, exposed for advanced users.</summary>
    public Transport Transport { get; }

    /// <summary>Status + whoami resource client.</summary>
    public StatusClient Status { get; }

    /// <summary>Accounts resource client (list, get, sync).</summary>
    public AccountsClient Accounts { get; }

    /// <summary>Transactions resource client (list, export, sync, bulk, search, get).</summary>
    public TransactionsClient Transactions { get; }

    /// <summary>Sync sessions resource client.</summary>
    public SyncSessionsClient SyncSessions { get; }

    /// <summary>Transaction orders resource client.</summary>
    public TransactionOrdersClient TransactionOrders { get; }

    /// <summary>Batches resource client.</summary>
    public BatchesClient Batches { get; }

    /// <summary>Payment methods resource client.</summary>
    public PaymentMethodsClient PaymentMethods { get; }

    /// <summary>Construct a v2 client from <see cref="ClientOptions"/>.</summary>
    public V2Client(ClientOptions options)
    {
        ArgumentNullException.ThrowIfNull(options);
        Transport = new Transport(options);
        Status = new StatusClient(Transport);
        Accounts = new AccountsClient(Transport);
        Transactions = new TransactionsClient(Transport);
        SyncSessions = new SyncSessionsClient(Transport);
        TransactionOrders = new TransactionOrdersClient(Transport);
        Batches = new BatchesClient(Transport);
        PaymentMethods = new PaymentMethodsClient(Transport);
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
