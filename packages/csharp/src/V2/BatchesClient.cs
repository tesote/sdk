using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V2;

/// <summary>v2 batches resource — atomic creation of multiple orders + bulk lifecycle.</summary>
public sealed class BatchesClient
{
    private readonly Transport _transport;

    internal BatchesClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>Single order entry inside a batch create payload.</summary>
    public sealed class BatchOrderInput
    {
        /// <summary>Existing payment-method id (null when supplying a beneficiary).</summary>
        public string? DestinationPaymentMethodId { get; set; }
        /// <summary>Inline beneficiary; mutually exclusive with the payment-method id.</summary>
        public Beneficiary? Beneficiary { get; set; }
        /// <summary>Order amount as a decimal string.</summary>
        public string Amount { get; set; } = "0";
        /// <summary>ISO currency code.</summary>
        public string Currency { get; set; } = "VES";
        /// <summary>Free-text description.</summary>
        public string Description { get; set; } = string.Empty;
        /// <summary>Optional schedule.</summary>
        public string? ScheduledFor { get; set; }
        /// <summary>Optional metadata bag.</summary>
        public IDictionary<string, object?>? Metadata { get; set; }
    }

    /// <summary>Create a new batch of transaction orders for an account.</summary>
    public async Task<BatchCreateResponse> CreateAsync(
        string accountId,
        IReadOnlyList<BatchOrderInput> orders,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentNullException.ThrowIfNull(orders);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId) + "/batches";
        var serialized = new List<Dictionary<string, object?>>(orders.Count);
        foreach (var o in orders)
        {
            serialized.Add(new Dictionary<string, object?>
            {
                ["destination_payment_method_id"] = o.DestinationPaymentMethodId,
                ["beneficiary"] = o.Beneficiary,
                ["amount"] = o.Amount,
                ["currency"] = o.Currency,
                ["description"] = o.Description,
                ["scheduled_for"] = o.ScheduledFor,
                ["metadata"] = o.Metadata,
            });
        }
        var body = new Dictionary<string, object?> { ["orders"] = serialized };
        var opts = Requests.Json("POST", path, body, idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<BatchCreateResponse>(node);
    }

    /// <summary>Fetch the summary view for a batch.</summary>
    public async Task<BatchSummary> GetAsync(string accountId, string batchId, CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentException.ThrowIfNullOrEmpty(batchId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId)
            + "/batches/" + Uri.EscapeDataString(batchId);
        var opts = RequestOptions.Get(path);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<BatchSummary>(node);
    }

    /// <summary>Approve every draft order in a batch.</summary>
    public async Task<BatchApproveResponse> ApproveAsync(
        string accountId,
        string batchId,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentException.ThrowIfNullOrEmpty(batchId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId)
            + "/batches/" + Uri.EscapeDataString(batchId) + "/approve";
        var opts = Requests.Json("POST", path, body: null, idempotencyKey: idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<BatchApproveResponse>(node);
    }

    /// <summary>Submit every approved order in a batch.</summary>
    public async Task<BatchSubmitResponse> SubmitAsync(
        string accountId,
        string batchId,
        string? token = null,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentException.ThrowIfNullOrEmpty(batchId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId)
            + "/batches/" + Uri.EscapeDataString(batchId) + "/submit";
        var body = new Dictionary<string, object?> { ["token"] = token };
        var opts = Requests.Json("POST", path, body, idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<BatchSubmitResponse>(node);
    }

    /// <summary>Cancel every cancellable order in a batch.</summary>
    public async Task<BatchCancelResponse> CancelAsync(
        string accountId,
        string batchId,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(accountId);
        ArgumentException.ThrowIfNullOrEmpty(batchId);
        var path = "/v2/accounts/" + Uri.EscapeDataString(accountId)
            + "/batches/" + Uri.EscapeDataString(batchId) + "/cancel";
        var opts = Requests.Json("POST", path, body: null, idempotencyKey: idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<BatchCancelResponse>(node);
    }
}
