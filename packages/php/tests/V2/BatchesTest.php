<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V2;

use Tesote\Sdk\Errors\BatchNotFoundException;
use Tesote\Sdk\Errors\BatchValidationException;
use Tesote\Sdk\Errors\InvalidOrderStateException;
use Tesote\Sdk\Models\BatchActionResult;
use Tesote\Sdk\Models\BatchCreated;
use Tesote\Sdk\Models\BatchSummary;
use Tesote\Sdk\Tests\Support\TestCaseBase;
use Tesote\Sdk\V2\Client;

final class BatchesTest extends TestCaseBase
{
    public function testCreate(): void
    {
        $this->enqueueOk(['batch_id' => 'b-1', 'orders' => [], 'errors' => []], 201);
        $client = new Client(['transport' => $this->makeTransport()]);
        $batch = $client->batches->create('acct-1', [
            ['amount' => '500.00', 'currency' => 'VES', 'description' => 'A'],
            ['amount' => '500.00', 'currency' => 'VES', 'description' => 'B'],
        ], 'idem-batch-1');

        self::assertInstanceOf(BatchCreated::class, $batch);
        self::assertSame('b-1', $batch->batchId);
        self::assertSame('idem-batch-1', $this->lastHeaders()['Idempotency-Key']);
        $body = json_decode($this->lastBody() ?? '', true);
        self::assertCount(2, $body['orders']);
    }

    public function testCreateMapsBatchValidationError(): void
    {
        $this->enqueueError(400, ['error' => 'bad batch', 'error_code' => 'BATCH_VALIDATION_ERROR']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(BatchValidationException::class);
        $client->batches->create('acct-1', []);
    }

    public function testGetSummary(): void
    {
        $this->enqueueOk([
            'batch_id' => 'b-1',
            'total_orders' => 2,
            'total_amount_cents' => 100000,
            'amount_currency' => 'VES',
            'statuses' => ['draft' => 2],
            'batch_status' => 'draft',
            'created_at' => '2026-04-28T10:00:00Z',
            'orders' => [],
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $summary = $client->batches->get('acct-1', 'b-1');

        self::assertInstanceOf(BatchSummary::class, $summary);
        self::assertSame(['draft' => 2], $summary->statuses);
    }

    public function testGetMapsNotFound(): void
    {
        $this->enqueueError(404, ['error' => 'gone', 'error_code' => 'BATCH_NOT_FOUND']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(BatchNotFoundException::class);
        $client->batches->get('acct-1', 'missing');
    }

    public function testApproveMutates(): void
    {
        $this->enqueueOk(['approved' => 5, 'failed' => 0]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $result = $client->batches->approve('acct-1', 'b-1');

        self::assertInstanceOf(BatchActionResult::class, $result);
        self::assertSame(5, $result->approved);
        self::assertNotEmpty($this->lastHeaders()['Idempotency-Key']);
        self::assertStringEndsWith('/approve', $this->lastUrl());
    }

    public function testSubmitWithToken(): void
    {
        $this->enqueueOk(['enqueued' => 5, 'failed' => 0]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $client->batches->submit('acct-1', 'b-1', 'OTP-1');

        $body = json_decode($this->lastBody() ?? '', true);
        self::assertSame('OTP-1', $body['token']);
    }

    public function testCancelMapsInvalidOrderState(): void
    {
        $this->enqueueError(409, ['error' => 'cant', 'error_code' => 'INVALID_ORDER_STATE']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(InvalidOrderStateException::class);
        $client->batches->cancel('acct-1', 'b-1');
    }
}
