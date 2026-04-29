<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/**
 * Counterparty descriptor.
 *
 * Transaction.counterparty is name-only; PaymentMethod.counterparty also
 * carries the id. Both shapes parse via fromArray() — id stays null when
 * absent.
 */
final class Counterparty
{
    public function __construct(
        public readonly ?string $id,
        public readonly string $name,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            id: isset($data['id']) ? (string) $data['id'] : null,
            name: (string) ($data['name'] ?? ''),
        );
    }
}
