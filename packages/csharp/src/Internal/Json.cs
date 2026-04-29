using System;
using System.Text.Json;
using System.Text.Json.Nodes;
using System.Text.Json.Serialization;

namespace Tesote.Sdk.Internal;

/// <summary>
/// Thin wrapper around shared <see cref="JsonSerializerOptions"/>. Default
/// options keep snake_case-on-the-wire by relying on per-property
/// <see cref="JsonPropertyNameAttribute"/> markers in the model records.
/// </summary>
public static class Json
{
    /// <summary>Shared default options. Treat as immutable.</summary>
    public static readonly JsonSerializerOptions DefaultOptions = new()
    {
        PropertyNameCaseInsensitive = true,
        WriteIndented = false,
        DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull,
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

    /// <summary>Deserialize a parsed <see cref="JsonNode"/> into a typed model record.</summary>
    public static T Deserialize<T>(JsonNode? node)
    {
        if (node is null)
        {
            throw new InvalidOperationException("cannot deserialize null JSON node");
        }
        var raw = node.ToJsonString(DefaultOptions);
        var result = JsonSerializer.Deserialize<T>(raw, DefaultOptions);
        if (result is null)
        {
            throw new InvalidOperationException("JSON deserialized to null for " + typeof(T).FullName);
        }
        return result;
    }
}
