<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /v2/accounts/{id}/batches/{batch_id}. */
final class BatchSummary
{
    /**
     * @param array<string, int>     $statuses
     * @param list<TransactionOrder> $orders
     */
    public function __construct(
        public readonly string $batchId,
        public readonly int $totalOrders,
        public readonly int $totalAmountCents,
        public readonly string $amountCurrency,
        public readonly array $statuses,
        public readonly string $batchStatus,
        public readonly string $createdAt,
        public readonly array $orders,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $statuses = [];
        foreach ((is_array($data['statuses'] ?? null) ? $data['statuses'] : []) as $status => $count) {
            $statuses[(string) $status] = (int) $count;
        }
        $orders = [];
        foreach ((is_array($data['orders'] ?? null) ? $data['orders'] : []) as $entry) {
            if (is_array($entry)) {
                $orders[] = TransactionOrder::fromArray($entry);
            }
        }

        return new self(
            batchId: (string) ($data['batch_id'] ?? ''),
            totalOrders: (int) ($data['total_orders'] ?? 0),
            totalAmountCents: (int) ($data['total_amount_cents'] ?? 0),
            amountCurrency: (string) ($data['amount_currency'] ?? ''),
            statuses: $statuses,
            batchStatus: (string) ($data['batch_status'] ?? ''),
            createdAt: (string) ($data['created_at'] ?? ''),
            orders: $orders,
        );
    }
}
