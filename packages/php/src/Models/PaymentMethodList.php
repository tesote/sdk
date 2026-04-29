<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /v2/payment_methods. */
final class PaymentMethodList
{
    /**
     * @param list<PaymentMethod> $items
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
                $items[] = PaymentMethod::fromArray($entry);
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
