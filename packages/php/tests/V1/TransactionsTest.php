<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V1;

use Tesote\Sdk\Errors\TransactionNotFoundException;
use Tesote\Sdk\Models\Transaction;
use Tesote\Sdk\Tests\Support\TestCaseBase;
use Tesote\Sdk\V1\Client;

final class TransactionsTest extends TestCaseBase
{
    public function testGetReturnsTransaction(): void
    {
        $this->enqueueOk([
            'id' => 'txn-1',
            'status' => 'posted',
            'data' => [
                'amount_cents' => 12500,
                'currency' => 'VES',
                'description' => 'Salary',
                'transaction_date' => '2026-04-15',
                'created_at' => '2026-04-15T08:00:00Z',
                'note' => null,
                'running_balance_cents' => 50000,
            ],
            'tesote_imported_at' => '2026-04-15T08:00:00Z',
            'tesote_updated_at' => '2026-04-15T08:00:00Z',
            'transaction_categories' => [
                ['name' => 'Income', 'external_category_code' => 'INC', 'created_at' => 't', 'updated_at' => 't'],
            ],
            'counterparty' => ['name' => 'Employer'],
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $txn = $client->transactions->get('txn-1');

        self::assertInstanceOf(Transaction::class, $txn);
        self::assertSame(12500, $txn->data->amountCents);
        self::assertSame(50000, $txn->data->runningBalanceCents);
        self::assertCount(1, $txn->transactionCategories);
        self::assertSame('Income', $txn->transactionCategories[0]->name);
    }

    public function testGetMapsTransactionNotFound(): void
    {
        $this->enqueueError(404, ['error' => 'no', 'error_code' => 'TRANSACTION_NOT_FOUND']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(TransactionNotFoundException::class);
        $client->transactions->get('missing');
    }
}
