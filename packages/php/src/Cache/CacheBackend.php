<?php

declare(strict_types=1);

namespace Tesote\Sdk\Cache;

/**
 * Pluggable cache backend used by Transport for opt-in TTL caching of GETs.
 *
 * Implementations must be safe to share between requests in the same process.
 * Keys are opaque strings; values are the raw decoded JSON arrays returned by
 * the API plus minimal envelope metadata. Implementations should treat values
 * as opaque.
 */
interface CacheBackend
{
    /**
     * Fetch a previously stored value.
     *
     * @return array{body: mixed, headers: array<string, string>, status: int}|null
     */
    public function get(string $key): ?array;

    /**
     * Store a value with a TTL in seconds. A TTL of 0 means do not store.
     *
     * @param array{body: mixed, headers: array<string, string>, status: int} $value
     */
    public function set(string $key, array $value, int $ttlSeconds): void;

    /**
     * Drop a single key. Used when a mutating call invalidates a path.
     */
    public function delete(string $key): void;
}
