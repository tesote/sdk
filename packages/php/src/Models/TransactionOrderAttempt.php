<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** TransactionOrder.latest_attempt — most-recent execution attempt summary. */
final readonly class TransactionOrderAttempt
{
    public function __construct(
        public string $id,
        public string $status,
        public int $attemptNumber,
        public ?string $externalReference,
        public ?string $submittedAt,
        public ?string $completedAt,
        public ?string $errorCode,
        public ?string $errorMessage,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            id: (string) ($data['id'] ?? ''),
            status: (string) ($data['status'] ?? ''),
            attemptNumber: (int) ($data['attempt_number'] ?? 0),
            externalReference: isset($data['external_reference']) ? (string) $data['external_reference'] : null,
            submittedAt: isset($data['submitted_at']) ? (string) $data['submitted_at'] : null,
            completedAt: isset($data['completed_at']) ? (string) $data['completed_at'] : null,
            errorCode: isset($data['error_code']) ? (string) $data['error_code'] : null,
            errorMessage: isset($data['error_message']) ? (string) $data['error_message'] : null,
        );
    }
}
