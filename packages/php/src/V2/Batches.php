<?php

declare(strict_types=1);

namespace Tesote\Sdk\V2;

use Tesote\Sdk\Models\BatchActionResult;
use Tesote\Sdk\Models\BatchCreated;
use Tesote\Sdk\Models\BatchSummary;
use Tesote\Sdk\Transport;

/** v2 batches: create, get, approve, submit, cancel. */
final class Batches
{
    public function __construct(private readonly Transport $transport)
    {
    }

    /**
     * @param list<array<string, mixed>> $orders
     */
    public function create(string $accountId, array $orders, ?string $idempotencyKey = null): BatchCreated
    {
        $body = $this->transport->request(
            'POST',
            '/v2/accounts/' . rawurlencode($accountId) . '/batches',
            null,
            ['orders' => $orders],
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return BatchCreated::fromArray($body);
    }

    public function get(string $accountId, string $batchId): BatchSummary
    {
        $body = $this->transport->request(
            'GET',
            '/v2/accounts/' . rawurlencode($accountId) . '/batches/' . rawurlencode($batchId),
        ) ?? [];
        return BatchSummary::fromArray($body);
    }

    public function approve(string $accountId, string $batchId, ?string $idempotencyKey = null): BatchActionResult
    {
        $body = $this->transport->request(
            'POST',
            '/v2/accounts/' . rawurlencode($accountId) . '/batches/' . rawurlencode($batchId) . '/approve',
            null,
            [],
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return BatchActionResult::fromArray($body);
    }

    public function submit(string $accountId, string $batchId, ?string $token = null, ?string $idempotencyKey = null): BatchActionResult
    {
        $payload = $token !== null ? ['token' => $token] : (object) [];
        $body = $this->transport->request(
            'POST',
            '/v2/accounts/' . rawurlencode($accountId) . '/batches/' . rawurlencode($batchId) . '/submit',
            null,
            (array) $payload,
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return BatchActionResult::fromArray($body);
    }

    public function cancel(string $accountId, string $batchId, ?string $idempotencyKey = null): BatchActionResult
    {
        $body = $this->transport->request(
            'POST',
            '/v2/accounts/' . rawurlencode($accountId) . '/batches/' . rawurlencode($batchId) . '/cancel',
            null,
            [],
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return BatchActionResult::fromArray($body);
    }
}
