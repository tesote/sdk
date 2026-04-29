<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Account.legal_entity — owning legal entity (id and legal_name may be null). */
final class LegalEntity
{
    public function __construct(
        public readonly ?string $id,
        public readonly ?string $legalName,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            id: isset($data['id']) ? (string) $data['id'] : null,
            legalName: isset($data['legal_name']) ? (string) $data['legal_name'] : null,
        );
    }
}
