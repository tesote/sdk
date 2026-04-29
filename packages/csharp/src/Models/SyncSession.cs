using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Tesote.Sdk.Models;

/// <summary>Failure information attached to a failed <see cref="SyncSession"/>.</summary>
public sealed record SyncSessionError(
    [property: JsonPropertyName("type")] string Type,
    [property: JsonPropertyName("message")] string Message);

/// <summary>Performance metrics for a completed <see cref="SyncSession"/>.</summary>
public sealed record SyncSessionPerformance(
    [property: JsonPropertyName("total_duration")] double TotalDuration,
    [property: JsonPropertyName("complexity_score")] double ComplexityScore,
    [property: JsonPropertyName("sync_speed_score")] double SyncSpeedScore);

/// <summary>Bank-sync session record.</summary>
public sealed record SyncSession(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("status")] string Status,
    [property: JsonPropertyName("started_at")] string StartedAt,
    [property: JsonPropertyName("completed_at")] string? CompletedAt,
    [property: JsonPropertyName("transactions_synced")] int TransactionsSynced,
    [property: JsonPropertyName("accounts_count")] int AccountsCount,
    [property: JsonPropertyName("error")] SyncSessionError? Error,
    [property: JsonPropertyName("performance")] SyncSessionPerformance? Performance);

/// <summary>Acceptance envelope for POST /v2/accounts/{id}/sync.</summary>
public sealed record SyncStartResponse(
    [property: JsonPropertyName("message")] string Message,
    [property: JsonPropertyName("sync_session_id")] string SyncSessionId,
    [property: JsonPropertyName("status")] string Status,
    [property: JsonPropertyName("started_at")] string StartedAt);

/// <summary>Response envelope for GET /v2/accounts/{id}/sync_sessions.</summary>
public sealed record SyncSessionListResponse(
    [property: JsonPropertyName("sync_sessions")] IReadOnlyList<SyncSession> SyncSessions,
    [property: JsonPropertyName("limit")] int Limit,
    [property: JsonPropertyName("offset")] int Offset,
    [property: JsonPropertyName("has_more")] bool HasMore);
