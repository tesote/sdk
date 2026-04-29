<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/**
 * Transaction.data subobject (v1 schema).
 *
 * runningBalanceCents only present when the workspace has running-balance
 * display enabled and the caller opted in.
 */
final class TransactionData
{
    public function __construct(
        public readonly int $amountCents,
        public readonly string $currency,
        public readonly string $description,
        public readonly string $transactionDate,
        public readonly ?string $createdAt,
        public readonly ?string $createdAtDate,
        public readonly ?string $note,
        public readonly ?string $externalServiceId,
        public readonly ?int $runningBalanceCents,
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
