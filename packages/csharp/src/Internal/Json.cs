using System.Text.Json;
using System.Text.Json.Nodes;

namespace Tesote.Sdk.Internal;

/// <summary>
/// Thin wrapper around shared <see cref="JsonSerializerOptions"/>. Default
/// options use snake_case property naming so SDK model records can declare
/// PascalCase properties and still match the wire shape.
/// </summary>
public static class Json
{
    /// <summary>Shared default options. Treat as immutable.</summary>
    public static readonly JsonSerializerOptions DefaultOptions = new()
    {
        PropertyNamingPolicy = JsonNamingPolicy.SnakeCaseLower,
        PropertyNameCaseInsensitive = true,
        WriteIndented = false,
    };

    /// <summary>Parse a UTF-8 byte buffer into a tolerant <see cref="JsonNode"/>; null on parse failure.</summary>
    public static JsonNode? Parse(byte[]? bytes)
    {
        if (bytes is null || bytes.Length == 0)
        {
            return null;
        }
        try
        {
            return JsonNode.Parse(bytes);
        }
        catch (JsonException)
        {
            return null;
        }
    }

    /// <summary>Serialize an arbitrary value to UTF-8 JSON bytes using the shared options.</summary>
    public static byte[] SerializeToUtf8Bytes<T>(T value)
    {
        return JsonSerializer.SerializeToUtf8Bytes(value, DefaultOptions);
    }
}
