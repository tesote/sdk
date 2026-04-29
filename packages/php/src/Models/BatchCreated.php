<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from POST /v2/accounts/{id}/batches. */
final readonly class BatchCreated
{
    /**
     * @param list<TransactionOrder>     $orders
     * @param list<array<string, mixed>> $errors
     */
    public function __construct(
        public string $batchId,
        public array $orders,
        public array $errors,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $orders = [];
        foreach ((is_array($data['orders'] ?? null) ? $data['orders'] : []) as $entry) {
            if (is_array($entry)) {
                $orders[] = TransactionOrder::fromArray($entry);
            }
        }
        $errors = [];
        foreach ((is_array($data['errors'] ?? null) ? $data['errors'] : []) as $entry) {
            if (is_array($entry)) {
                /** @var array<string, mixed> $entry */
                $errors[] = $entry;
            }
        }

        return new self(
            batchId: (string) ($data['batch_id'] ?? ''),
            orders: $orders,
            errors: $errors,
        );
    }
}
