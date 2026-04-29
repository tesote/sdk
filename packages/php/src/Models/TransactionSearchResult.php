<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /v2/transactions/search. */
final readonly class TransactionSearchResult
{
    /**
     * @param list<Transaction> $transactions
     */
    public function __construct(
        public array $transactions,
        public int $total,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $transactions = [];
        foreach ((is_array($data['transactions'] ?? null) ? $data['transactions'] : []) as $entry) {
            if (is_array($entry)) {
                $transactions[] = Transaction::fromArray($entry);
            }
        }

        return new self(
            transactions: $transactions,
            total: (int) ($data['total'] ?? 0),
        );
    }
}
