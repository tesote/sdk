<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Transaction.transaction_categories[] entry. */
final class TransactionCategory
{
    public function __construct(
        public readonly string $name,
        public readonly ?string $externalCategoryCode,
        public readonly string $createdAt,
        public readonly string $updatedAt,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            name: (string) ($data['name'] ?? ''),
            externalCategoryCode: isset($data['external_category_code']) ? (string) $data['external_category_code'] : null,
            createdAt: (string) ($data['created_at'] ?? ''),
            updatedAt: (string) ($data['updated_at'] ?? ''),
        );
    }
}
