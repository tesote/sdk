<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** TransactionOrder model (v2). */
final readonly class TransactionOrder
{
    public function __construct(
        public string $id,
        public string $status,
        public string $amount,
        public string $currency,
        public string $description,
        public ?string $reference,
        public ?string $externalReference,
        public ?string $idempotencyKey,
        public ?string $batchId,
        public ?string $scheduledFor,
        public ?string $approvedAt,
        public ?string $submittedAt,
        public ?string $completedAt,
        public ?string $failedAt,
        public ?string $cancelledAt,
        public SourceAccount $sourceAccount,
        public Destination $destination,
        public ?Money $fee,
        public ?string $executionStrategy,
        public ?TesoteTransactionRef $tesoteTransaction,
        public ?TransactionOrderAttempt $latestAttempt,
        public string $createdAt,
        public string $updatedAt,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $source = is_array($data['source_account'] ?? null) ? $data['source_account'] : [];
        $dest = is_array($data['destination'] ?? null) ? $data['destination'] : [];
        $fee = is_array($data['fee'] ?? null) ? $data['fee'] : null;
        $tesoteTxn = is_array($data['tesote_transaction'] ?? null) ? $data['tesote_transaction'] : null;
        $attempt = is_array($data['latest_attempt'] ?? null) ? $data['latest_attempt'] : null;

        return new self(
            id: (string) ($data['id'] ?? ''),
            status: (string) ($data['status'] ?? ''),
            amount: (string) ($data['amount'] ?? '0'),
            currency: (string) ($data['currency'] ?? ''),
            description: (string) ($data['description'] ?? ''),
            reference: isset($data['reference']) ? (string) $data['reference'] : null,
            externalReference: isset($data['external_reference']) ? (string) $data['external_reference'] : null,
            idempotencyKey: isset($data['idempotency_key']) ? (string) $data['idempotency_key'] : null,
            batchId: isset($data['batch_id']) ? (string) $data['batch_id'] : null,
            scheduledFor: isset($data['scheduled_for']) ? (string) $data['scheduled_for'] : null,
            approvedAt: isset($data['approved_at']) ? (string) $data['approved_at'] : null,
            submittedAt: isset($data['submitted_at']) ? (string) $data['submitted_at'] : null,
            completedAt: isset($data['completed_at']) ? (string) $data['completed_at'] : null,
            failedAt: isset($data['failed_at']) ? (string) $data['failed_at'] : null,
            cancelledAt: isset($data['cancelled_at']) ? (string) $data['cancelled_at'] : null,
            sourceAccount: SourceAccount::fromArray($source),
            destination: Destination::fromArray($dest),
            fee: $fee !== null ? Money::fromArray($fee) : null,
            executionStrategy: isset($data['execution_strategy']) ? (string) $data['execution_strategy'] : null,
            tesoteTransaction: $tesoteTxn !== null ? TesoteTransactionRef::fromArray($tesoteTxn) : null,
            latestAttempt: $attempt !== null ? TransactionOrderAttempt::fromArray($attempt) : null,
            createdAt: (string) ($data['created_at'] ?? ''),
            updatedAt: (string) ($data['updated_at'] ?? ''),
        );
    }
}
