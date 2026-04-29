<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /v1/accounts/{id}/transactions and GET /v2/accounts/{id}/transactions. */
final class TransactionList
{
    /**
     * @param list<Transaction> $transactions
     */
    public function __construct(
        public readonly int $total,
        public readonly array $transactions,
        public readonly CursorPagination $pagination,
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
        $pagination = is_array($data['pagination'] ?? null) ? $data['pagination'] : [];

        return new self(
            total: (int) ($data['total'] ?? 0),
            transactions: $transactions,
            pagination: CursorPagination::fromArray($pagination),
        );
    }
}
