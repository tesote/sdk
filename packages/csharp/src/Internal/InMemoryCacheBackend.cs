using System;
using System.Collections.Concurrent;

namespace Tesote.Sdk.Internal;

/// <summary>
/// Default thread-safe in-memory backend with TTL expiry. Not LRU-bounded —
/// users with real cache pressure plug in their own backend.
/// </summary>
public sealed class InMemoryCacheBackend : ICacheBackend
{
    private readonly ConcurrentDictionary<string, Entry> _store = new();

    /// <inheritdoc />
    public byte[]? Get(string key)
    {
        ArgumentNullException.ThrowIfNull(key);
        if (!_store.TryGetValue(key, out var entry))
        {
            return null;
        }
        if (DateTimeOffset.UtcNow > entry.ExpiresAt)
        {
            _store.TryRemove(key, out _);
            return null;
        }
        return entry.Value;
    }

    /// <inheritdoc />
    public void Put(string key, byte[] value, TimeSpan ttl)
    {
        ArgumentNullException.ThrowIfNull(key);
        ArgumentNullException.ThrowIfNull(value);
        _store[key] = new Entry(value, DateTimeOffset.UtcNow.Add(ttl));
    }

    /// <inheritdoc />
    public void Invalidate(string key)
    {
        ArgumentNullException.ThrowIfNull(key);
        _store.TryRemove(key, out _);
    }

    /// <inheritdoc />
    public void InvalidatePrefix(string prefix)
    {
        ArgumentNullException.ThrowIfNull(prefix);
        foreach (var k in _store.Keys)
        {
            if (k.StartsWith(prefix, StringComparison.Ordinal))
            {
                _store.TryRemove(k, out _);
            }
        }
    }

    private readonly record struct Entry(byte[] Value, DateTimeOffset ExpiresAt);
}
