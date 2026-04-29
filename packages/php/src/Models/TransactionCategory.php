<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Transaction.transaction_categories[] entry. */
final readonly class TransactionCategory
{
    public function __construct(
        public string $name,
        public ?string $externalCategoryCode,
        public string $createdAt,
        public string $updatedAt,
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
