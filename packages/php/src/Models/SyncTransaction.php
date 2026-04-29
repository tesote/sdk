<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/**
 * SyncTransaction — flattened, Plaid-compatible shape used by the
 * /v2/.../transactions/sync endpoints. Distinct from Transaction.
 */
final class SyncTransaction
{
    /**
     * @param list<string> $category
     */
    public function __construct(
        public readonly string $transactionId,
        public readonly string $accountId,
        public readonly float $amount,
        public readonly string $isoCurrencyCode,
        public readonly ?string $unofficialCurrencyCode,
        public readonly string $date,
        public readonly ?string $datetime,
        public readonly string $name,
        public readonly ?string $merchantName,
        public readonly bool $pending,
        public readonly array $category,
        public readonly ?int $runningBalanceCents,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $rawCategory = is_array($data['category'] ?? null) ? $data['category'] : [];
        $category = [];
        foreach ($rawCategory as $entry) {
            $category[] = (string) $entry;
        }

        return new self(
            transactionId: (string) ($data['transaction_id'] ?? ''),
            accountId: (string) ($data['account_id'] ?? ''),
            amount: (float) ($data['amount'] ?? 0.0),
            isoCurrencyCode: (string) ($data['iso_currency_code'] ?? ''),
            unofficialCurrencyCode: isset($data['unofficial_currency_code']) ? (string) $data['unofficial_currency_code'] : null,
            date: (string) ($data['date'] ?? ''),
            datetime: isset($data['datetime']) ? (string) $data['datetime'] : null,
            name: (string) ($data['name'] ?? ''),
            merchantName: isset($data['merchant_name']) ? (string) $data['merchant_name'] : null,
            pending: (bool) ($data['pending'] ?? false),
            category: $category,
            runningBalanceCents: isset($data['running_balance_cents']) ? (int) $data['running_balance_cents'] : null,
        );
    }
}
