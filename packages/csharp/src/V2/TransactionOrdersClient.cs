using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V2;

/// <summary>v2 transaction-orders resource — list, get, create, submit, cancel.</summary>
public sealed class TransactionOrdersClient
{
    private readonly Transport _transport;

    internal TransactionOrdersClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>List transaction orders for an account.</summary>
    public async Task<TransactionOrderListResponse> ListAsync(
        string accountId,
        int? limit = null,
        int? offset = null,
        string? status = null,
        string? createdAfter = null,
        string? createdBefore = null,
        string? batchId = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId) + "/transaction_orders";
        var query = new QueryBuilder()
            .Add("limit", limit)
            .Add("offset", offset)
            .Add("status", status)
            .Add("created_after", createdAfter)
            .Add("created_before", createdBefore)
            .Add("batch_id", batchId)
            .BuildOrNull();
        var opts = RequestOptions.Get(path);
        opts.Query = query;
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<TransactionOrderListResponse>(node);
    }

    /// <summary>Fetch a single transaction order.</summary>
    public async Task<TransactionOrder> GetAsync(string accountId, string orderId, CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentException.ThrowIfNullOrEmpty(orderId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId)
            + "/transaction_orders/" + Uri.EscapeDataString(orderId);
        var opts = RequestOptions.Get(path);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<TransactionOrder>(node);
    }

    /// <summary>Request body for POST /v2/.../transaction_orders.</summary>
    public sealed class CreateRequest
    {
        /// <summary>Existing payment-method id, or null when supplying a beneficiary.</summary>
        public string? DestinationPaymentMethodId { get; set; }
        /// <summary>Inline beneficiary; mutually exclusive with <see cref="DestinationPaymentMethodId"/>.</summary>
        public Beneficiary? Beneficiary { get; set; }
        /// <summary>Order amount as a decimal string.</summary>
        public string Amount { get; set; } = "0";
        /// <summary>ISO currency code, e.g. <c>VES</c>.</summary>
        public string Currency { get; set; } = "VES";
        /// <summary>Free-text description.</summary>
        public string Description { get; set; } = string.Empty;
        /// <summary>Optional ISO8601 schedule timestamp.</summary>
        public string? ScheduledFor { get; set; }
        /// <summary>Server-side idempotency key (separate from the transport idempotency header).</summary>
        public string? IdempotencyKey { get; set; }
        /// <summary>Optional metadata bag.</summary>
        public IDictionary<string, object?>? Metadata { get; set; }
    }

    /// <summary>Create a draft transaction order.</summary>
    public async Task<TransactionOrder> CreateAsync(
        string accountId,
        CreateRequest request,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentNullException.ThrowIfNull(request);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId) + "/transaction_orders";
        var inner = new Dictionary<string, object?>
        {
            ["destination_payment_method_id"] = request.DestinationPaymentMethodId,
            ["beneficiary"] = request.Beneficiary,
            ["amount"] = request.Amount,
            ["currency"] = request.Currency,
            ["description"] = request.Description,
            ["scheduled_for"] = request.ScheduledFor,
            ["idempotency_key"] = request.IdempotencyKey,
            ["metadata"] = request.Metadata,
        };
        var body = new Dictionary<string, object?> { ["transaction_order"] = inner };
        var opts = Requests.Json("POST", path, body, idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<TransactionOrder>(node);
    }

    /// <summary>Submit a draft (or pending-approval) order for processing.</summary>
    public async Task<TransactionOrder> SubmitAsync(
        string accountId,
        string orderId,
        string? token = null,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentException.ThrowIfNullOrEmpty(orderId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId)
            + "/transaction_orders/" + Uri.EscapeDataString(orderId) + "/submit";
        var body = new Dictionary<string, object?> { ["token"] = token };
        var opts = Requests.Json("POST", path, body, idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<TransactionOrder>(node);
    }

    /// <summary>Cancel an order; transitions to <c>cancelled</c>.</summary>
    public async Task<TransactionOrder> CancelAsync(
        string accountId,
        string orderId,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentException.ThrowIfNullOrEmpty(orderId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId)
            + "/transaction_orders/" + Uri.EscapeDataString(orderId) + "/cancel";
        var opts = Requests.Json("POST", path, body: null, idempotencyKey: idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<TransactionOrder>(node);
    }
}
