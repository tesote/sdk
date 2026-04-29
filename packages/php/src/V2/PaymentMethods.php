<?php

declare(strict_types=1);

namespace Tesote\Sdk\V2;

use Tesote\Sdk\Models\PaymentMethod;
use Tesote\Sdk\Models\PaymentMethodList;
use Tesote\Sdk\Transport;

/** v2 payment methods: CRUD. */
final class PaymentMethods
{
    public function __construct(private readonly Transport $transport)
    {
    }

    /**
     * @param array{
     *     limit?: int,
     *     offset?: int,
     *     method_type?: string,
     *     currency?: string,
     *     counterparty_id?: string,
     *     verified?: 'true'|'false',
     * } $query
     */
    public function list(array $query = []): PaymentMethodList
    {
        $body = $this->transport->request('GET', '/v2/payment_methods', $query) ?? [];
        return PaymentMethodList::fromArray($body);
    }

    public function get(string $id): PaymentMethod
    {
        $body = $this->transport->request('GET', '/v2/payment_methods/' . rawurlencode($id)) ?? [];
        return PaymentMethod::fromArray($body);
    }

    /**
     * @param array{
     *     method_type: string,
     *     currency: string,
     *     label?: string|null,
     *     counterparty_id?: string|null,
     *     counterparty?: array<string, mixed>,
     *     details: array<string, mixed>,
     * } $paymentMethod
     */
    public function create(array $paymentMethod, ?string $idempotencyKey = null): PaymentMethod
    {
        $body = $this->transport->request(
            'POST',
            '/v2/payment_methods',
            null,
            ['payment_method' => $paymentMethod],
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return PaymentMethod::fromArray($body);
    }

    /**
     * @param array<string, mixed> $changes
     */
    public function update(string $id, array $changes, ?string $idempotencyKey = null): PaymentMethod
    {
        $body = $this->transport->request(
            'PATCH',
            '/v2/payment_methods/' . rawurlencode($id),
            null,
            ['payment_method' => $changes],
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return PaymentMethod::fromArray($body);
    }

    public function delete(string $id, ?string $idempotencyKey = null): void
    {
        $this->transport->request(
            'DELETE',
            '/v2/payment_methods/' . rawurlencode($id),
            null,
            null,
            ['idempotencyKey' => $idempotencyKey],
        );
    }
}
