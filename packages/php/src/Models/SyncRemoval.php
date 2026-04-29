<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Entry in the `removed` array from the /v2/transactions/sync response. */
final class SyncRemoval
{
    public function __construct(
        public readonly string $transactionId,
        public readonly string $accountId,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            transactionId: (string) ($data['transaction_id'] ?? ''),
            accountId: (string) ($data['account_id'] ?? ''),
        );
    }
}
