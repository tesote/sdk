<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/**
 * Generic money triple — used by TransactionOrder.fee.
 *
 * Amount stored as string for decimal-safety (matches the wire). Callers that
 * need numeric comparison should parse with bcmath / ext-decimal.
 */
final readonly class Money
{
    public function __construct(
        public string $amount,
        public string $currency,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            amount: (string) ($data['amount'] ?? '0'),
            currency: (string) ($data['currency'] ?? ''),
        );
    }
}
