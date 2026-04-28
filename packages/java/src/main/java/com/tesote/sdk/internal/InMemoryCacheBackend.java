package com.tesote.sdk.internal;

import com.tesote.sdk.CacheBackend;

import java.time.Duration;
import java.time.Instant;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Default thread-safe in-memory backend with TTL expiry. Not LRU-bounded — the
 * design doc suggests LRU but the default is fine for the bootstrap; users with
 * real cache pressure plug in their own backend.
 */
public final class InMemoryCacheBackend implements CacheBackend {
    private record Entry(byte[] value, Instant expiresAt) {}

    private final Map<String, Entry> store = new ConcurrentHashMap<>();

    @Override
    public Optional<byte[]> get(String key) {
        Entry e = store.get(key);
        if (e == null) return Optional.empty();
        if (Instant.now().isAfter(e.expiresAt())) {
            store.remove(key, e);
            return Optional.empty();
        }
        return Optional.of(e.value());
    }

    @Override
    public void put(String key, byte[] value, Duration ttl) {
        store.put(key, new Entry(value, Instant.now().plus(ttl)));
    }

    @Override
    public void invalidate(String key) {
        store.remove(key);
    }

    @Override
    public void invalidatePrefix(String prefix) {
        store.keySet().removeIf(k -> k.startsWith(prefix));
    }
}
