<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V2;

use Tesote\Sdk\Errors\InvalidQueryException;
use Tesote\Sdk\Errors\TransactionNotFoundException;
use Tesote\Sdk\Errors\UnprocessableContentException;
use Tesote\Sdk\Models\BulkResult;
use Tesote\Sdk\Models\Transaction;
use Tesote\Sdk\Models\TransactionSearchResult;
use Tesote\Sdk\Tests\Support\TestCaseBase;
use Tesote\Sdk\V2\Client;

final class TransactionsTest extends TestCaseBase
{
    public function testGet(): void
    {
        $this->enqueueOk([
            'id' => 'txn-1',
            'status' => 'posted',
            'data' => [
                'amount_cents' => 100,
                'currency' => 'VES',
                'description' => 'x',
                'transaction_date' => '2026-04-20',
            ],
            'tesote_imported_at' => 't',
            'tesote_updated_at' => 't',
            'transaction_categories' => [],
            'counterparty' => null,
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $txn = $client->transactions->get('txn-1');
        self::assertInstanceOf(Transaction::class, $txn);
        self::assertSame('txn-1', $txn->id);
        self::assertStringEndsWith('/v2/transactions/txn-1', $this->lastUrl());
    }

    public function testGetMapsNotFound(): void
    {
        $this->enqueueError(404, ['error' => 'gone', 'error_code' => 'TRANSACTION_NOT_FOUND']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(TransactionNotFoundException::class);
        $client->transactions->get('missing');
    }

    public function testLegacySync(): void
    {
        $this->enqueueOk([
            'added' => [], 'modified' => [], 'removed' => [], 'next_cursor' => null, 'has_more' => false,
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $client->transactions->sync(['count' => 50, 'cursor' => null]);
        self::assertSame('POST', $this->lastMethod());
        self::assertStringEndsWith('/v2/transactions/sync', $this->lastUrl());
        // why: SDK auto-generates an idempotency key for mutations even without an explicit caller value.
        self::assertNotEmpty($this->lastHeaders()['Idempotency-Key']);
    }

    public function testBulkRespectsBody(): void
    {
        $this->enqueueOk([
            'bulk_results' => [
                [
                    'account_id' => 'a-1',
                    'transactions' => [],
                    'pagination' => ['has_more' => false, 'per_page' => 50, 'after_id' => null, 'before_id' => null],
                ],
            ],
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $bulk = $client->transactions->bulk(['account_ids' => ['a-1', 'a-2'], 'limit' => 50]);

        self::assertInstanceOf(BulkResult::class, $bulk);
        self::assertSame('a-1', $bulk->bulkResults[0]->accountId);

        $body = json_decode($this->lastBody() ?? '', true);
        self::assertSame(['a-1', 'a-2'], $body['account_ids']);
    }

    public function testBulkMapsUnprocessable(): void
    {
        $this->enqueueError(422, ['error' => 'too many', 'error_code' => 'UNPROCESSABLE_CONTENT']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(UnprocessableContentException::class);
        $client->transactions->bulk(['account_ids' => array_fill(0, 200, 'x')]);
    }

    public function testSearchBuildsQuery(): void
    {
        $this->enqueueOk(['transactions' => [], 'total' => 0]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $result = $client->transactions->search(['q' => 'coffee', 'limit' => 25]);

        self::assertInstanceOf(TransactionSearchResult::class, $result);
        self::assertStringContainsString('q=coffee', $this->lastUrl());
        self::assertStringContainsString('limit=25', $this->lastUrl());
    }

    public function testSearchMapsInvalidQuery(): void
    {
        $this->enqueueError(422, ['error' => 'no q', 'error_code' => 'INVALID_QUERY']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(InvalidQueryException::class);
        $client->transactions->search(['q' => '']);
    }
}
