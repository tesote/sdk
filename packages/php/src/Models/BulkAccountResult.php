<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** One row from the bulk_results array of POST /v2/transactions/bulk. */
final class BulkAccountResult
{
    /**
     * @param list<Transaction> $transactions
     */
    public function __construct(
        public readonly string $accountId,
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
            accountId: (string) ($data['account_id'] ?? ''),
            transactions: $transactions,
            pagination: CursorPagination::fromArray($pagination),
        );
    }
}
