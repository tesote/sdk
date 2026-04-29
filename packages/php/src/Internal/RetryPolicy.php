<?php

declare(strict_types=1);

namespace Tesote\Sdk\Internal;

use Tesote\Sdk\Http\CurlErrorClassifier;
use Tesote\Sdk\Http\CurlResult;

/**
 * Decides whether a given response or transport failure should be retried.
 */
final class RetryPolicy
{
    private const MUTATING_METHODS = ['POST', 'PUT', 'PATCH', 'DELETE'];
    private const RETRIABLE_STATUSES = [429, 502, 503, 504];

    public static function isMutating(string $method): bool
    {
        return in_array($method, self::MUTATING_METHODS, true);
    }

    public static function isRetriableStatus(int $status): bool
    {
        return in_array($status, self::RETRIABLE_STATUSES, true);
    }

    public static function shouldRetryTransport(CurlResult $result, string $method, ?string $idempotencyKey): bool
    {
        // why: docs/architecture/transport.md — never retry POST without idempotency key on read timeout.
        $isReadTimeout = $result->errno === CURLE_OPERATION_TIMEOUTED;
        if ($isReadTimeout && self::isMutating($method) && $idempotencyKey === null) {
            return false;
        }
        return CurlErrorClassifier::isRetriableErrno($result->errno);
    }

    /**
     * @param array<string, string> $headers
     */
    public static function parseRetryAfter(array $headers): ?int
    {
        if (isset($headers['retry-after']) && is_numeric($headers['retry-after'])) {
            return (int) $headers['retry-after'];
        }
        return null;
    }
}
