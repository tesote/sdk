<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /v2/accounts/{id}/transaction_orders. */
final class TransactionOrderList
{
    /**
     * @param list<TransactionOrder> $items
     */
    public function __construct(
        public readonly array $items,
        public readonly bool $hasMore,
        public readonly int $limit,
        public readonly int $offset,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $items = [];
        foreach ((is_array($data['items'] ?? null) ? $data['items'] : []) as $entry) {
            if (is_array($entry)) {
                $items[] = TransactionOrder::fromArray($entry);
            }
        }

        return new self(
            items: $items,
            hasMore: (bool) ($data['has_more'] ?? false),
            limit: (int) ($data['limit'] ?? 0),
            offset: (int) ($data['offset'] ?? 0),
        );
    }
}
