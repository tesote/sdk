<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** PaymentMethod.tesote_account — back-reference to a source account. */
final readonly class TesoteAccountRef
{
    public function __construct(
        public string $id,
        public string $name,
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
        );
    }
}
