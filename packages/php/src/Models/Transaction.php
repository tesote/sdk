<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Transaction (v1 schema — also returned by GET /v2/transactions/{id}). */
final class Transaction
{
    /**
     * @param list<TransactionCategory> $transactionCategories
     */
    public function __construct(
        public readonly string $id,
        public readonly string $status,
        public readonly TransactionData $data,
        public readonly string $tesoteImportedAt,
        public readonly string $tesoteUpdatedAt,
        public readonly array $transactionCategories,
        public readonly ?Counterparty $counterparty,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $rawData = is_array($data['data'] ?? null) ? $data['data'] : [];
        $rawCategories = is_array($data['transaction_categories'] ?? null) ? $data['transaction_categories'] : [];
        $rawCounterparty = is_array($data['counterparty'] ?? null) ? $data['counterparty'] : null;

        $categories = [];
        foreach ($rawCategories as $entry) {
            if (is_array($entry)) {
                $categories[] = TransactionCategory::fromArray($entry);
            }
        }

        return new self(
            id: (string) ($data['id'] ?? ''),
            status: (string) ($data['status'] ?? ''),
            data: TransactionData::fromArray($rawData),
            tesoteImportedAt: (string) ($data['tesote_imported_at'] ?? ''),
            tesoteUpdatedAt: (string) ($data['tesote_updated_at'] ?? ''),
            transactionCategories: $categories,
            counterparty: $rawCounterparty !== null ? Counterparty::fromArray($rawCounterparty) : null,
        );
    }
}
