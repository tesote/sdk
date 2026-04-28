<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

use RuntimeException;
use Throwable;

/**
 * Root of the SDK exception hierarchy.
 *
 * Every error thrown by the SDK extends this. Catch this only as a last
 * resort — real call sites should catch the typed subclass that matches
 * the error_code or transport failure they care about.
 *
 * The required-fields contract from docs/architecture/errors.md is
 * implemented here so subclasses don't redefine the same properties.
 */
class TesoteException extends RuntimeException
{
    /** @var array<string, mixed>|null */
    public readonly ?array $requestSummary;

    /**
     * @param array<string, mixed>|null $requestSummary
     */
    public function __construct(
        string $message,
        public readonly string $errorCode,
        public readonly int $httpStatus,
        public readonly ?string $requestId,
        public readonly ?string $errorId,
        public readonly ?int $retryAfter,
        public readonly ?string $responseBody,
        ?array $requestSummary,
        public readonly int $attempts,
        ?Throwable $previous = null,
    ) {
        parent::__construct($message, 0, $previous);
        $this->requestSummary = $requestSummary;
    }

    /**
     * Greppable single-line summary suitable for logs.
     */
    public function summary(): string
    {
        $parts = [
            static::class . ': ' . $this->httpStatus . ' ' . $this->message,
            'error_code=' . $this->errorCode,
        ];
        if ($this->requestId !== null) {
            $parts[] = 'request_id=' . $this->requestId;
        }
        if ($this->retryAfter !== null) {
            $parts[] = 'retry_after=' . $this->retryAfter . 's';
        }
        $parts[] = 'attempts=' . $this->attempts;
        return implode(' ', $parts);
    }
}
