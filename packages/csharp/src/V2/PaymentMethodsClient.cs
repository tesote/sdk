using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using Tesote.Sdk.Internal;
using Tesote.Sdk.Models;

namespace Tesote.Sdk.V2;

/// <summary>v2 payment-methods resource — list, get, create, update, soft-delete.</summary>
public sealed class PaymentMethodsClient
{
    private const string BasePath = "/v2/payment_methods";

    private readonly Transport _transport;

    internal PaymentMethodsClient(Transport transport)
    {
        _transport = transport;
    }

    /// <summary>List payment methods (offset paginated). Soft-deleted entries are excluded.</summary>
    public async Task<PaymentMethodListResponse> ListAsync(
        int? limit = null,
        int? offset = null,
        string? methodType = null,
        string? currency = null,
        string? counterpartyId = null,
        bool? verified = null,
        CancellationToken ct = default)
    {
        var query = new QueryBuilder()
            .Add("limit", limit)
            .Add("offset", offset)
            .Add("method_type", methodType)
            .Add("currency", currency)
            .Add("counterparty_id", counterpartyId)
            .Add("verified", verified)
            .BuildOrNull();
        var opts = RequestOptions.Get(BasePath);
        opts.Query = query;
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<PaymentMethodListResponse>(node);
    }

    /// <summary>Fetch a single payment method by id.</summary>
    public async Task<PaymentMethod> GetAsync(string id, CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(id);
        var opts = RequestOptions.Get(BasePath + "/" + Uri.EscapeDataString(id));
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<PaymentMethod>(node);
    }

    /// <summary>Counterparty stub accepted by create when no <c>counterparty_id</c> is provided.</summary>
    public sealed class CounterpartyInput
    {
        /// <summary>Counterparty display name.</summary>
        public string Name { get; set; } = string.Empty;
    }

    /// <summary>Type-specific payload for a payment method.</summary>
    public sealed class DetailsInput
    {
        /// <summary>Bank routing/clabe code.</summary>
        public string? BankCode { get; set; }
        /// <summary>Account number.</summary>
        public string? AccountNumber { get; set; }
        /// <summary>Holder full name.</summary>
        public string? HolderName { get; set; }
        /// <summary>Identification document type.</summary>
        public string? IdentificationType { get; set; }
        /// <summary>Identification document number.</summary>
        public string? IdentificationNumber { get; set; }
    }

    /// <summary>Body for POST/PATCH /v2/payment_methods.</summary>
    public sealed class WriteRequest
    {
        /// <summary>Method type (e.g. <c>bank_account</c>, <c>pago_movil</c>, <c>wire</c>).</summary>
        public string? MethodType { get; set; }
        /// <summary>ISO currency code.</summary>
        public string? Currency { get; set; }
        /// <summary>Optional human label.</summary>
        public string? Label { get; set; }
        /// <summary>Existing counterparty id; mutually exclusive with <see cref="Counterparty"/>.</summary>
        public string? CounterpartyId { get; set; }
        /// <summary>Inline counterparty stub.</summary>
        public CounterpartyInput? Counterparty { get; set; }
        /// <summary>Type-specific details payload.</summary>
        public DetailsInput? Details { get; set; }
    }

    /// <summary>Create a new payment method.</summary>
    public Task<PaymentMethod> CreateAsync(
        WriteRequest request,
        string? idempotencyKey = null,
        CancellationToken ct = default) => Write("POST", BasePath, request, idempotencyKey, ct);

    /// <summary>Patch an existing payment method.</summary>
    public Task<PaymentMethod> UpdateAsync(
        string id,
        WriteRequest request,
        string? idempotencyKey = null,
        CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(id);
        return Write("PATCH", BasePath + "/" + Uri.EscapeDataString(id), request, idempotencyKey, ct);
    }

    /// <summary>Soft-delete a payment method. 204 No Content on success.</summary>
    public async Task DeleteAsync(string id, string? idempotencyKey = null, CancellationToken ct = default)
    {
        ArgumentException.ThrowIfNullOrEmpty(id);
        var opts = Requests.Json("DELETE", BasePath + "/" + Uri.EscapeDataString(id), body: null, idempotencyKey: idempotencyKey);
        await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
    }

    private async Task<PaymentMethod> Write(string method, string path, WriteRequest request, string? idempotencyKey, CancellationToken ct)
    {
        ArgumentNullException.ThrowIfNull(request);
        var inner = new Dictionary<string, object?>
        {
            ["method_type"] = request.MethodType,
            ["currency"] = request.Currency,
            ["label"] = request.Label,
            ["counterparty_id"] = request.CounterpartyId,
            ["counterparty"] = request.Counterparty is null ? null : new Dictionary<string, object?> { ["name"] = request.Counterparty.Name },
            ["details"] = request.Details is null ? null : new Dictionary<string, object?>
            {
                ["bank_code"] = request.Details.BankCode,
                ["account_number"] = request.Details.AccountNumber,
                ["holder_name"] = request.Details.HolderName,
                ["identification_type"] = request.Details.IdentificationType,
                ["identification_number"] = request.Details.IdentificationNumber,
            },
        };
        var body = new Dictionary<string, object?> { ["payment_method"] = inner };
        var opts = Requests.Json(method, path, body, idempotencyKey);
        var node = await _transport.RequestAsync(opts, ct).ConfigureAwait(false);
        return Json.Deserialize<PaymentMethod>(node);
    }
}
