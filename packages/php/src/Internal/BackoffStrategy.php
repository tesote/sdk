<?php

declare(strict_types=1);

namespace Tesote\Sdk\Internal;

/**
 * Computes retry sleep durations and dispatches the actual sleep call.
 *
 * The sleeper callable is injectable so tests can record sleeps without
 * blocking the test process.
 */
final class BackoffStrategy
{
    /** @var (callable(int): void)|null */
    private $sleeper;

    /**
     * @param (callable(int): void)|null $sleeper
     */
    public function __construct(
        private readonly int $baseDelayMs,
        private readonly int $maxDelayMs,
        ?callable $sleeper = null,
    ) {
        $this->sleeper = $sleeper;
    }

    public function delayMs(int $attempt, ?int $retryAfterSeconds): int
    {
        if ($retryAfterSeconds !== null) {
            return min($this->maxDelayMs, $retryAfterSeconds * 1000);
        }
        $expo = $this->baseDelayMs * (2 ** ($attempt - 1));
        $capped = min($this->maxDelayMs, $expo);
        // why: full jitter — distributes retries to avoid thundering herd.
        return random_int(0, max(1, (int) $capped));
    }

    public function sleep(int $milliseconds): void
    {
        if ($milliseconds <= 0) {
            return;
        }
        if ($this->sleeper !== null) {
            ($this->sleeper)($milliseconds);
            return;
        }
        usleep($milliseconds * 1000);
    }
}
