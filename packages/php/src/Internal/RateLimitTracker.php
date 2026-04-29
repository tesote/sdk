<?php

declare(strict_types=1);

namespace Tesote\Sdk\Internal;

/**
 * Holds the most recently observed rate-limit + request-id values.
 *
 * Lives outside Transport so that the orchestration loop can stay focused
 * on retry/cache logic; this is plain mutable state.
 */
final class RateLimitTracker
{
    /** @var array{limit: ?string, remaining: ?string, reset: ?string} */
    private array $lastRateLimit = [
        'limit' => null,
        'remaining' => null,
        'reset' => null,
    ];

    private ?string $lastRequestId = null;

    /**
     * @param array<string, string> $headers
     */
    public function record(array $headers): void
    {
        $this->lastRateLimit = [
            'limit' => $headers['x-ratelimit-limit'] ?? null,
            'remaining' => $headers['x-ratelimit-remaining'] ?? null,
            'reset' => $headers['x-ratelimit-reset'] ?? null,
        ];
        $this->lastRequestId = $headers['x-request-id'] ?? null;
    }

    /**
     * @return array{limit: ?string, remaining: ?string, reset: ?string}
     */
    public function lastRateLimit(): array
    {
        return $this->lastRateLimit;
    }

    public function lastRequestId(): ?string
    {
        return $this->lastRequestId;
    }
}
