<?php

declare(strict_types=1);

namespace Tesote\Sdk\Cache;

/**
 * Tiny in-process TTL cache. Not LRU-bounded — caller manages lifecycle.
 *
 * Suitable as a default for short-lived processes (CLI scripts, single
 * request handlers). For long-lived workers, inject a Redis-backed
 * CacheBackend instead.
 */
final class InMemoryCache implements CacheBackend
{
    /** @var array<string, array{value: array{body: mixed, headers: array<string, string>, status: int}, expires_at: int}> */
    private array $store = [];

    public function get(string $key): ?array
    {
        if (!isset($this->store[$key])) {
            return null;
        }

        $entry = $this->store[$key];
        if ($entry['expires_at'] <= time()) {
            unset($this->store[$key]);
            return null;
        }

        return $entry['value'];
    }

    public function set(string $key, array $value, int $ttlSeconds): void
    {
        if ($ttlSeconds <= 0) {
            return;
        }

        $this->store[$key] = [
            'value' => $value,
            'expires_at' => time() + $ttlSeconds,
        ];
    }

    public function delete(string $key): void
    {
        unset($this->store[$key]);
    }
}
