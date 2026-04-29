<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Account.bank — minimal bank descriptor (name only on the wire). */
final readonly class Bank
{
    public function __construct(
        public string $name,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(name: (string) ($data['name'] ?? ''));
    }
}
