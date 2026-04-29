<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /v2/payment_methods. */
final readonly class PaymentMethodList
{
    /**
     * @param list<PaymentMethod> $items
     */
    public function __construct(
        public array $items,
        public bool $hasMore,
        public int $limit,
        public int $offset,
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
