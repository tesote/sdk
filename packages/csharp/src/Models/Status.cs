using System.Text.Json.Serialization;

namespace Tesote.Sdk.Models;

/// <summary>Response envelope for GET /status and /v2/status.</summary>
public sealed record StatusResponse(
    [property: JsonPropertyName("status")] string Status,
    [property: JsonPropertyName("authenticated")] bool Authenticated);

/// <summary>Identification details returned from GET /whoami and /v2/whoami.</summary>
public sealed record WhoamiClient(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("name")] string Name,
    [property: JsonPropertyName("type")] string Type);

/// <summary>Response envelope for GET /whoami and /v2/whoami.</summary>
public sealed record WhoamiResponse(
    [property: JsonPropertyName("client")] WhoamiClient Client);
