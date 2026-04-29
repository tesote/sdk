<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** TransactionOrder.source_account — the account funding the order. */
final readonly class SourceAccount
{
    public function __construct(
        public string $id,
        public string $name,
        public string $paymentMethodId,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            id: (string) ($data['id'] ?? ''),
            name: (string) ($data['name'] ?? ''),
            paymentMethodId: (string) ($data['payment_method_id'] ?? ''),
        );
    }
}
