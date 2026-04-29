<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** TransactionOrder.destination — beneficiary identifiers. */
final class Destination
{
    public function __construct(
        public readonly string $paymentMethodId,
        public readonly ?string $counterpartyId,
        public readonly ?string $counterpartyName,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            paymentMethodId: (string) ($data['payment_method_id'] ?? ''),
            counterpartyId: isset($data['counterparty_id']) ? (string) $data['counterparty_id'] : null,
            counterpartyName: isset($data['counterparty_name']) ? (string) $data['counterparty_name'] : null,
        );
    }
}
