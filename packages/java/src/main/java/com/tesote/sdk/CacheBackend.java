package com.tesote.sdk;

import java.time.Duration;
import java.util.Optional;

/**
 * Pluggable cache for the opt-in TTL response cache. Default in-memory
 * implementation lives in {@link com.tesote.sdk.internal.InMemoryCacheBackend};
 * users can drop in Redis/memcached by implementing this interface.
 */
public interface CacheBackend {
    Optional<byte[]> get(String key);

    void put(String key, byte[] value, Duration ttl);

    void invalidate(String key);

    void invalidatePrefix(String prefix);
}
