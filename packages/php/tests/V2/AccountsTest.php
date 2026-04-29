<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V2;

use Tesote\Sdk\Errors\AccountNotFoundException;
use Tesote\Sdk\Errors\BankConnectionNotFoundException;
use Tesote\Sdk\Errors\BankUnderMaintenanceException;
use Tesote\Sdk\Errors\HistorySyncForbiddenException;
use Tesote\Sdk\Errors\InvalidCountException;
use Tesote\Sdk\Errors\InvalidCursorException;
use Tesote\Sdk\Errors\SyncInProgressException;
use Tesote\Sdk\Errors\SyncRateLimitExceededException;
use Tesote\Sdk\Models\AccountList;
use Tesote\Sdk\Models\ExportFile;
use Tesote\Sdk\Models\SyncResult;
use Tesote\Sdk\Models\SyncStarted;
use Tesote\Sdk\Tests\Support\TestCaseBase;
use Tesote\Sdk\V2\Client;

final class AccountsTest extends TestCaseBase
{
    public function testListReturnsModel(): void
    {
        $this->enqueueOk([
            'total' => 0,
            'accounts' => [],
            'pagination' => ['current_page' => 1, 'per_page' => 50, 'total_pages' => 0, 'total_count' => 0],
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $list = $client->accounts->list(['per_page' => 50]);

        self::assertInstanceOf(AccountList::class, $list);
        self::assertSame(0, $list->total);
        self::assertStringContainsString('/v2/accounts?per_page=50', $this->lastUrl());
    }

    public function testSyncReturnsSyncStartedAndSendsIdempotencyKey(): void
    {
        $this->enqueueOk([
            'message' => 'Sync started',
            'sync_session_id' => 'ss-1',
            'status' => 'pending',
            'started_at' => '2026-04-28T10:00:00Z',
        ], 202);
        $client = new Client(['transport' => $this->makeTransport()]);
        $started = $client->accounts->sync('acct-1', 'caller-key-1');

        self::assertInstanceOf(SyncStarted::class, $started);
        self::assertSame('ss-1', $started->syncSessionId);
        self::assertSame('caller-key-1', $this->lastHeaders()['Idempotency-Key']);
        self::assertSame('POST', $this->lastMethod());
        self::assertStringEndsWith('/v2/accounts/acct-1/sync', $this->lastUrl());
    }

    public function testSyncMapsAllSyncErrorCodes(): void
    {
        $cases = [
            ['ACCOUNT_NOT_FOUND', 404, AccountNotFoundException::class],
            ['BANK_CONNECTION_NOT_FOUND', 404, BankConnectionNotFoundException::class],
            ['SYNC_IN_PROGRESS', 409, SyncInProgressException::class],
            ['SYNC_RATE_LIMIT_EXCEEDED', 429, SyncRateLimitExceededException::class],
            ['BANK_UNDER_MAINTENANCE', 503, BankUnderMaintenanceException::class],
        ];
        foreach ($cases as [$code, $status, $class]) {
            $this->setUp();
            // why: 429/503 retry, so we need 3 entries to surface the exception within maxAttempts=3.
            $needRetries = in_array($code, ['SYNC_RATE_LIMIT_EXCEEDED', 'BANK_UNDER_MAINTENANCE'], true);
            $count = $needRetries ? 3 : 1;
            for ($i = 0; $i < $count; $i++) {
                $this->enqueueError($status, ['error' => 'x', 'error_code' => $code]);
            }
            $client = new Client(['transport' => $this->makeTransport()]);
            try {
                $client->accounts->sync('acct-1');
                self::fail("expected $class for $code");
            } catch (\Throwable $e) {
                self::assertInstanceOf($class, $e, "wrong class for $code");
            }
        }
    }

    public function testSyncTransactionsCursorAndCountErrors(): void
    {
        $this->enqueueError(422, ['error' => 'bad', 'error_code' => 'INVALID_CURSOR']);
        $client = new Client(['transport' => $this->makeTransport()]);
        try {
            $client->accounts->syncTransactions('acct-1', ['cursor' => 'broken']);
            self::fail('expected InvalidCursorException');
        } catch (InvalidCursorException $e) {
            self::assertSame('INVALID_CURSOR', $e->errorCode);
        }

        $this->setUp();
        $this->enqueueError(422, ['error' => 'too many', 'error_code' => 'INVALID_COUNT']);
        $client = new Client(['transport' => $this->makeTransport()]);
        try {
            $client->accounts->syncTransactions('acct-1', ['count' => 9999]);
            self::fail('expected InvalidCountException');
        } catch (InvalidCountException $e) {
            self::assertSame('INVALID_COUNT', $e->errorCode);
        }

        $this->setUp();
        $this->enqueueError(403, ['error' => 'pre-cutoff', 'error_code' => 'HISTORY_SYNC_FORBIDDEN']);
        $client = new Client(['transport' => $this->makeTransport()]);
        try {
            $client->accounts->syncTransactions('acct-1', ['cursor' => 'old']);
            self::fail('expected HistorySyncForbiddenException');
        } catch (HistorySyncForbiddenException $e) {
            self::assertSame('HISTORY_SYNC_FORBIDDEN', $e->errorCode);
        }
    }

    public function testSyncTransactionsHappyPath(): void
    {
        $this->enqueueOk([
            'added' => [[
                'transaction_id' => 't-1',
                'account_id' => 'acct-1',
                'amount' => 99.5,
                'iso_currency_code' => 'VES',
                'date' => '2026-04-20',
                'name' => 'Coffee',
                'pending' => false,
                'category' => ['food'],
            ]],
            'modified' => [],
            'removed' => [['transaction_id' => 't-x', 'account_id' => 'acct-1']],
            'next_cursor' => 'cursor-next',
            'has_more' => false,
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $result = $client->accounts->syncTransactions(
            'acct-1',
            ['count' => 100, 'cursor' => 'now', 'options' => ['include_running_balance' => true]],
            'idem-1',
        );

        self::assertInstanceOf(SyncResult::class, $result);
        self::assertCount(1, $result->added);
        self::assertSame('t-1', $result->added[0]->transactionId);
        self::assertSame('cursor-next', $result->nextCursor);
        self::assertCount(1, $result->removed);
        self::assertSame('idem-1', $this->lastHeaders()['Idempotency-Key']);

        $body = json_decode($this->lastBody() ?? '', true);
        self::assertSame(100, $body['count']);
        self::assertSame('now', $body['cursor']);
    }

    public function testExportTransactionsReturnsRawCsv(): void
    {
        $csv = "Transaction ID,Date,Description\nt-1,2026-04-20,Coffee\n";
        $this->enqueueRaw($csv, 200, [
            'content-disposition' => 'attachment; filename="transactions_acct-1_2026-04-28.csv"',
            'content-type' => 'text/csv',
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $file = $client->accounts->exportTransactions('acct-1', ['format' => 'csv', 'start_date' => '2026-04-01']);

        self::assertInstanceOf(ExportFile::class, $file);
        self::assertSame('csv', $file->format);
        self::assertSame('transactions_acct-1_2026-04-28.csv', $file->filename);
        self::assertSame($csv, $file->body);
        self::assertStringContainsString('format=csv', $this->lastUrl());
    }

    public function testRequest415WhenContentTypeMissing(): void
    {
        // why: emulate the server's 415 response when Content-Type is missing on a mutation.
        $this->enqueueError(415, ['error' => 'need json', 'error_code' => 'UNSUPPORTED_MEDIA_TYPE']);
        $client = new Client(['transport' => $this->makeTransport()]);
        try {
            $client->accounts->syncTransactions('acct-1', ['count' => 50]);
            self::fail('expected ApiException');
        } catch (\Tesote\Sdk\Errors\ApiException $e) {
            self::assertSame(415, $e->httpStatus);
        }
    }
}
