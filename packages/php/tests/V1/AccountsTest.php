<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V1;

use Tesote\Sdk\Errors\AccountNotFoundException;
use Tesote\Sdk\Errors\InvalidDateRangeException;
use Tesote\Sdk\Errors\UnauthorizedException;
use Tesote\Sdk\Models\Account;
use Tesote\Sdk\Models\AccountList;
use Tesote\Sdk\Models\TransactionList;
use Tesote\Sdk\Tests\Support\TestCaseBase;
use Tesote\Sdk\V1\Client;

final class AccountsTest extends TestCaseBase
{
    public function testListReturnsTypedAccountList(): void
    {
        $this->enqueueOk([
            'total' => 1,
            'accounts' => [[
                'id' => 'acct-1',
                'name' => 'Checking',
                'data' => [
                    'masked_account_number' => '****1234',
                    'currency' => 'VES',
                    'transactions_data_current_as_of' => '2026-04-28T00:00:00Z',
                    'balance_data_current_as_of' => '2026-04-28T00:00:00Z',
                    'custom_user_provided_identifier' => 'main',
                    'balance_cents' => '12345',
                ],
                'bank' => ['name' => 'Mercantil'],
                'legal_entity' => ['id' => 'le-1', 'legal_name' => 'ACME C.A.'],
                'tesote_created_at' => '2026-01-01T00:00:00Z',
                'tesote_updated_at' => '2026-04-01T00:00:00Z',
            ]],
            'pagination' => [
                'current_page' => 1,
                'per_page' => 50,
                'total_pages' => 1,
                'total_count' => 1,
            ],
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $list = $client->accounts->list(['page' => 1, 'per_page' => 50]);

        self::assertInstanceOf(AccountList::class, $list);
        self::assertSame(1, $list->total);
        self::assertSame('acct-1', $list->accounts[0]->id);
        self::assertSame('Mercantil', $list->accounts[0]->bank->name);
        self::assertSame('12345', $list->accounts[0]->data->balanceCents);
        self::assertStringContainsString('/v1/accounts?page=1&per_page=50', $this->lastUrl());
        self::assertSame('GET', $this->lastMethod());
    }

    public function testGetReturnsAccount(): void
    {
        $this->enqueueOk([
            'id' => 'acct-1',
            'name' => 'Savings',
            'data' => ['currency' => 'VES'],
            'bank' => ['name' => 'BNC'],
            'legal_entity' => null,
            'tesote_created_at' => '2026-01-01T00:00:00Z',
            'tesote_updated_at' => '2026-04-01T00:00:00Z',
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $account = $client->accounts->get('acct-1');

        self::assertInstanceOf(Account::class, $account);
        self::assertNull($account->legalEntity);
        self::assertStringContainsString('/v1/accounts/acct-1', $this->lastUrl());
    }

    public function testGetMapsAccountNotFound(): void
    {
        $this->enqueueError(404, ['error' => 'gone', 'error_code' => 'ACCOUNT_NOT_FOUND']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(AccountNotFoundException::class);
        $client->accounts->get('missing');
    }

    public function testListMapsUnauthorized(): void
    {
        $this->enqueueError(401, ['error' => 'no key', 'error_code' => 'UNAUTHORIZED']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(UnauthorizedException::class);
        $client->accounts->list();
    }

    public function testListTransactionsCursorPagination(): void
    {
        $this->enqueueOk([
            'total' => 2,
            'transactions' => [
                [
                    'id' => 'txn-1',
                    'status' => 'posted',
                    'data' => [
                        'amount_cents' => 1000,
                        'currency' => 'VES',
                        'description' => 'Coffee',
                        'transaction_date' => '2026-04-20',
                    ],
                    'tesote_imported_at' => '2026-04-20T01:00:00Z',
                    'tesote_updated_at' => '2026-04-20T01:00:00Z',
                    'transaction_categories' => [],
                    'counterparty' => ['name' => 'Cafe'],
                ],
                [
                    'id' => 'txn-2',
                    'status' => 'posted',
                    'data' => [
                        'amount_cents' => 500,
                        'currency' => 'VES',
                        'description' => 'Bread',
                        'transaction_date' => '2026-04-21',
                    ],
                    'tesote_imported_at' => '2026-04-21T01:00:00Z',
                    'tesote_updated_at' => '2026-04-21T01:00:00Z',
                    'transaction_categories' => [],
                    'counterparty' => null,
                ],
            ],
            'pagination' => [
                'has_more' => true,
                'per_page' => 50,
                'after_id' => 'txn-2',
                'before_id' => 'txn-1',
            ],
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $list = $client->accounts->listTransactions('acct-1', [
            'transactions_after_id' => 'cursor-prev',
            'per_page' => 50,
        ]);

        self::assertInstanceOf(TransactionList::class, $list);
        self::assertCount(2, $list->transactions);
        self::assertTrue($list->pagination->hasMore);
        self::assertSame('txn-2', $list->pagination->afterId);
        self::assertStringContainsString('transactions_after_id=cursor-prev', $this->lastUrl());
    }

    public function testListTransactionsMapsInvalidDateRange(): void
    {
        $this->enqueueError(422, ['error' => 'bad range', 'error_code' => 'INVALID_DATE_RANGE']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(InvalidDateRangeException::class);
        $client->accounts->listTransactions('acct-1', ['start_date' => '2099-01-01', 'end_date' => '2025-01-01']);
    }
}
