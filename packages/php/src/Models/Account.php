<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Account model (v1 + v2 — identical wire shape). */
final readonly class Account
{
    public function __construct(
        public string $id,
        public string $name,
        public AccountData $data,
        public Bank $bank,
        public ?LegalEntity $legalEntity,
        public string $tesoteCreatedAt,
        public string $tesoteUpdatedAt,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $rawData = is_array($data['data'] ?? null) ? $data['data'] : [];
        $rawBank = is_array($data['bank'] ?? null) ? $data['bank'] : [];
        $rawLegal = is_array($data['legal_entity'] ?? null) ? $data['legal_entity'] : null;

        return new self(
            id: (string) ($data['id'] ?? ''),
            name: (string) ($data['name'] ?? ''),
            data: AccountData::fromArray($rawData),
            bank: Bank::fromArray($rawBank),
            legalEntity: $rawLegal !== null ? LegalEntity::fromArray($rawLegal) : null,
            tesoteCreatedAt: (string) ($data['tesote_created_at'] ?? ''),
            tesoteUpdatedAt: (string) ($data['tesote_updated_at'] ?? ''),
        );
    }
}
