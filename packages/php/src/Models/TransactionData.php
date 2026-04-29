<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/**
 * Transaction.data subobject (v1 schema).
 *
 * runningBalanceCents only present when the workspace has running-balance
 * display enabled and the caller opted in.
 */
final readonly class TransactionData
{
    public function __construct(
        public int $amountCents,
        public string $currency,
        public string $description,
        public string $transactionDate,
        public ?string $createdAt,
        public ?string $createdAtDate,
        public ?string $note,
        public ?string $externalServiceId,
        public ?int $runningBalanceCents,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            amountCents: (int) ($data['amount_cents'] ?? 0),
            currency: (string) ($data['currency'] ?? ''),
            description: (string) ($data['description'] ?? ''),
            transactionDate: (string) ($data['transaction_date'] ?? ''),
            createdAt: isset($data['created_at']) ? (string) $data['created_at'] : null,
            createdAtDate: isset($data['created_at_date']) ? (string) $data['created_at_date'] : null,
            note: isset($data['note']) ? (string) $data['note'] : null,
            externalServiceId: isset($data['external_service_id']) ? (string) $data['external_service_id'] : null,
            runningBalanceCents: isset($data['running_balance_cents']) ? (int) $data['running_balance_cents'] : null,
        );
    }
}
