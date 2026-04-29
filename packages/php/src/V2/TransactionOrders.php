<?php

declare(strict_types=1);

namespace Tesote\Sdk\V2;

use Tesote\Sdk\Models\TransactionOrder;
use Tesote\Sdk\Models\TransactionOrderList;
use Tesote\Sdk\Transport;

/** v2 transaction orders: list, get, create, submit, cancel. */
final class TransactionOrders
{
    public function __construct(private readonly Transport $transport)
    {
    }

    /**
     * @param array{
     *     limit?: int,
     *     offset?: int,
     *     status?: string,
     *     created_after?: string,
     *     created_before?: string,
     *     batch_id?: string,
     * } $query
     */
    public function listForAccount(string $accountId, array $query = []): TransactionOrderList
    {
        $body = $this->transport->request(
            'GET',
            '/v2/accounts/' . rawurlencode($accountId) . '/transaction_orders',
            $query,
        ) ?? [];
        return TransactionOrderList::fromArray($body);
    }

    public function get(string $accountId, string $orderId): TransactionOrder
    {
        $body = $this->transport->request(
            'GET',
            '/v2/accounts/' . rawurlencode($accountId) . '/transaction_orders/' . rawurlencode($orderId),
        ) ?? [];
        return TransactionOrder::fromArray($body);
    }

    /**
     * @param array{
     *     destination_payment_method_id?: string|null,
     *     beneficiary?: array<string, mixed>,
     *     amount: string,
     *     currency: string,
     *     description: string,
     *     scheduled_for?: string|null,
     *     idempotency_key?: string|null,
     *     metadata?: array<string, mixed>,
     * } $order
     */
    public function create(string $accountId, array $order, ?string $idempotencyKey = null): TransactionOrder
    {
        $body = $this->transport->request(
            'POST',
            '/v2/accounts/' . rawurlencode($accountId) . '/transaction_orders',
            null,
            ['transaction_order' => $order],
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return TransactionOrder::fromArray($body);
    }

    public function submit(string $accountId, string $orderId, ?string $token = null, ?string $idempotencyKey = null): TransactionOrder
    {
        $payload = $token !== null ? ['token' => $token] : (object) [];
        $body = $this->transport->request(
            'POST',
            '/v2/accounts/' . rawurlencode($accountId) . '/transaction_orders/' . rawurlencode($orderId) . '/submit',
            null,
            (array) $payload,
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return TransactionOrder::fromArray($body);
    }

    public function cancel(string $accountId, string $orderId, ?string $idempotencyKey = null): TransactionOrder
    {
        $body = $this->transport->request(
            'POST',
            '/v2/accounts/' . rawurlencode($accountId) . '/transaction_orders/' . rawurlencode($orderId) . '/cancel',
            null,
            [],
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return TransactionOrder::fromArray($body);
    }
}
