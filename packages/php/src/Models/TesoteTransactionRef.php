<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** TransactionOrder.tesote_transaction — reference to the underlying ledger transaction. */
final readonly class TesoteTransactionRef
{
    public function __construct(
        public string $id,
        public string $status,
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
        );
    }
}
