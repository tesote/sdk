using System;

namespace Tesote.Sdk;

/// <summary>
/// Pluggable cache for the opt-in TTL response cache. Default in-memory
/// implementation lives in <see cref="Internal.InMemoryCacheBackend"/>;
/// users can drop in Redis/memcached by implementing this interface.
/// </summary>
public interface ICacheBackend
{
    /// <summary>Return cached bytes for the key, or null when missing or expired.</summary>
    byte[]? Get(string key);

    /// <summary>Store bytes under the key with a TTL.</summary>
    void Put(string key, byte[] value, TimeSpan ttl);

    /// <summary>Remove a single key.</summary>
    void Invalidate(string key);

    /// <summary>Remove every key starting with the given prefix.</summary>
    void InvalidatePrefix(string prefix);
}
