<?php

declare(strict_types=1);

namespace Tesote\Sdk\Http;

/**
 * Plain value object describing what came back from a single cURL exec.
 *
 * Headers are normalised to lowercase keys for predictable lookup; the
 * original casing from the wire is irrelevant to callers.
 */
final class CurlResult
{
    /**
     * @param array<string, string> $headers Lowercased header name => last value seen.
     */
    public function __construct(
        public readonly int $status,
        public readonly string $body,
        public readonly array $headers,
        public readonly int $errno,
        public readonly string $errorMessage,
    ) {
    }
}
