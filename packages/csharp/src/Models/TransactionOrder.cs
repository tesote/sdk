using System.Collections.Generic;
using System.Text.Json.Nodes;
using System.Text.Json.Serialization;

namespace Tesote.Sdk.Models;

/// <summary>Source-account stub embedded in a <see cref="TransactionOrder"/>.</summary>
public sealed record TransactionOrderSourceAccount(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("name")] string Name,
    [property: JsonPropertyName("payment_method_id")] string PaymentMethodId);

/// <summary>Destination metadata embedded in a <see cref="TransactionOrder"/>.</summary>
public sealed record TransactionOrderDestination(
    [property: JsonPropertyName("payment_method_id")] string PaymentMethodId,
    [property: JsonPropertyName("counterparty_id")] string? CounterpartyId,
    [property: JsonPropertyName("counterparty_name")] string? CounterpartyName);

/// <summary>Optional fee charged on a <see cref="TransactionOrder"/>.</summary>
public sealed record TransactionOrderFee(
    [property: JsonPropertyName("amount")] decimal Amount,
    [property: JsonPropertyName("currency")] string Currency);

/// <summary>Underlying bank transaction created from a successful order.</summary>
public sealed record TransactionOrderTesoteTransaction(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("status")] string Status);

/// <summary>Most recent submission attempt for a <see cref="TransactionOrder"/>.</summary>
public sealed record TransactionOrderLatestAttempt(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("status")] string Status,
    [property: JsonPropertyName("attempt_number")] int AttemptNumber,
    [property: JsonPropertyName("external_reference")] string? ExternalReference,
    [property: JsonPropertyName("submitted_at")] string? SubmittedAt,
    [property: JsonPropertyName("completed_at")] string? CompletedAt,
    [property: JsonPropertyName("error_code")] string? ErrorCode,
    [property: JsonPropertyName("error_message")] string? ErrorMessage);

/// <summary>Beneficiary payload accepted by POST /v2/.../transaction_orders when no payment-method id is supplied.</summary>
public sealed record Beneficiary(
    [property: JsonPropertyName("name")] string Name,
    [property: JsonPropertyName("bank_code")] string? BankCode,
    [property: JsonPropertyName("account_number")] string? AccountNumber,
    [property: JsonPropertyName("identification_type")] string? IdentificationType,
    [property: JsonPropertyName("identification_number")] string? IdentificationNumber);

/// <summary>Single transfer order — the unit of payment in v2.</summary>
public sealed record TransactionOrder(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("status")] string Status,
    [property: JsonPropertyName("amount")] decimal Amount,
    [property: JsonPropertyName("currency")] string Currency,
    [property: JsonPropertyName("description")] string Description,
    [property: JsonPropertyName("reference")] string? Reference,
    [property: JsonPropertyName("external_reference")] string? ExternalReference,
    [property: JsonPropertyName("idempotency_key")] string? IdempotencyKey,
    [property: JsonPropertyName("batch_id")] string? BatchId,
    [property: JsonPropertyName("scheduled_for")] string? ScheduledFor,
    [property: JsonPropertyName("approved_at")] string? ApprovedAt,
    [property: JsonPropertyName("submitted_at")] string? SubmittedAt,
    [property: JsonPropertyName("completed_at")] string? CompletedAt,
    [property: JsonPropertyName("failed_at")] string? FailedAt,
    [property: JsonPropertyName("cancelled_at")] string? CancelledAt,
    [property: JsonPropertyName("source_account")] TransactionOrderSourceAccount SourceAccount,
    [property: JsonPropertyName("destination")] TransactionOrderDestination Destination,
    [property: JsonPropertyName("fee")] TransactionOrderFee? Fee,
    [property: JsonPropertyName("execution_strategy")] string? ExecutionStrategy,
    [property: JsonPropertyName("tesote_transaction")] TransactionOrderTesoteTransaction? TesoteTransaction,
    [property: JsonPropertyName("latest_attempt")] TransactionOrderLatestAttempt? LatestAttempt,
    [property: JsonPropertyName("metadata")] JsonObject? Metadata,
    [property: JsonPropertyName("created_at")] string CreatedAt,
    [property: JsonPropertyName("updated_at")] string UpdatedAt);

/// <summary>Response envelope for GET /v2/.../transaction_orders.</summary>
public sealed record TransactionOrderListResponse(
    [property: JsonPropertyName("items")] IReadOnlyList<TransactionOrder> Items,
    [property: JsonPropertyName("has_more")] bool HasMore,
    [property: JsonPropertyName("limit")] int Limit,
    [property: JsonPropertyName("offset")] int Offset);

/// <summary>Single error entry returned from POST /v2/.../batches.</summary>
public sealed record BatchOrderError(
    [property: JsonPropertyName("index")] int? Index,
    [property: JsonPropertyName("error")] string? Error,
    [property: JsonPropertyName("error_code")] string? ErrorCode);

/// <summary>Response envelope for POST /v2/.../batches.</summary>
public sealed record BatchCreateResponse(
    [property: JsonPropertyName("batch_id")] string BatchId,
    [property: JsonPropertyName("orders")] IReadOnlyList<TransactionOrder> Orders,
    [property: JsonPropertyName("errors")] IReadOnlyList<BatchOrderError> Errors);

/// <summary>Per-status order counts inside a batch.</summary>
public sealed record BatchStatusCounts(
    [property: JsonPropertyName("draft")] int Draft,
    [property: JsonPropertyName("pending_approval")] int PendingApproval,
    [property: JsonPropertyName("approved")] int Approved,
    [property: JsonPropertyName("processing")] int Processing,
    [property: JsonPropertyName("completed")] int Completed,
    [property: JsonPropertyName("failed")] int Failed,
    [property: JsonPropertyName("cancelled")] int Cancelled);

/// <summary>Response envelope for GET /v2/.../batches/{batch_id}.</summary>
public sealed record BatchSummary(
    [property: JsonPropertyName("batch_id")] string BatchId,
    [property: JsonPropertyName("total_orders")] int TotalOrders,
    [property: JsonPropertyName("total_amount_cents")] long TotalAmountCents,
    [property: JsonPropertyName("amount_currency")] string AmountCurrency,
    [property: JsonPropertyName("statuses")] BatchStatusCounts Statuses,
    [property: JsonPropertyName("batch_status")] string BatchStatus,
    [property: JsonPropertyName("created_at")] string CreatedAt,
    [property: JsonPropertyName("orders")] IReadOnlyList<TransactionOrder> Orders);

/// <summary>Response envelope for POST /v2/.../batches/{batch_id}/approve.</summary>
public sealed record BatchApproveResponse(
    [property: JsonPropertyName("approved")] int Approved,
    [property: JsonPropertyName("failed")] int Failed);

/// <summary>Response envelope for POST /v2/.../batches/{batch_id}/submit.</summary>
public sealed record BatchSubmitResponse(
    [property: JsonPropertyName("enqueued")] int Enqueued,
    [property: JsonPropertyName("failed")] int Failed);

/// <summary>Response envelope for POST /v2/.../batches/{batch_id}/cancel.</summary>
public sealed record BatchCancelResponse(
    [property: JsonPropertyName("cancelled")] int Cancelled,
    [property: JsonPropertyName("skipped")] int Skipped,
    [property: JsonPropertyName("errors")] IReadOnlyList<BatchOrderError> Errors);
