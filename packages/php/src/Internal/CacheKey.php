<?php

declare(strict_types=1);

namespace Tesote\Sdk\Internal;

/**
 * Derives a stable cache key for a (method, url, api-key-suffix) tuple.
 *
 * The api-key suffix is included so concurrent processes sharing one cache
 * backend with different keys cannot read each other's cached responses.
 */
final class CacheKey
{
    public static function for(string $method, string $url, string $apiKey): string
    {
        return hash('sha256', $method . ' ' . $url . ' ' . substr($apiKey, -4));
    }
}
