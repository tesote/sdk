<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** PaymentMethod (v2). details is type-specific so kept as a generic map. */
final readonly class PaymentMethod
{
    /**
     * @param array<string, mixed> $details
     */
    public function __construct(
        public string $id,
        public string $methodType,
        public string $currency,
        public ?string $label,
        public array $details,
        public bool $verified,
        public ?string $verifiedAt,
        public ?string $lastUsedAt,
        public ?Counterparty $counterparty,
        public ?TesoteAccountRef $tesoteAccount,
        public string $createdAt,
        public string $updatedAt,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $rawDetails = is_array($data['details'] ?? null) ? $data['details'] : [];
        $rawCounterparty = is_array($data['counterparty'] ?? null) ? $data['counterparty'] : null;
        $rawAccount = is_array($data['tesote_account'] ?? null) ? $data['tesote_account'] : null;

        return new self(
            id: (string) ($data['id'] ?? ''),
            methodType: (string) ($data['method_type'] ?? ''),
            currency: (string) ($data['currency'] ?? ''),
            label: isset($data['label']) ? (string) $data['label'] : null,
            details: $rawDetails,
            verified: (bool) ($data['verified'] ?? false),
            verifiedAt: isset($data['verified_at']) ? (string) $data['verified_at'] : null,
            lastUsedAt: isset($data['last_used_at']) ? (string) $data['last_used_at'] : null,
            counterparty: $rawCounterparty !== null ? Counterparty::fromArray($rawCounterparty) : null,
            tesoteAccount: $rawAccount !== null ? TesoteAccountRef::fromArray($rawAccount) : null,
            createdAt: (string) ($data['created_at'] ?? ''),
            updatedAt: (string) ($data['updated_at'] ?? ''),
        );
    }
}
