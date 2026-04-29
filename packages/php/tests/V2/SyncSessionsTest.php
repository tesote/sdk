<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V2;

use Tesote\Sdk\Errors\SyncSessionNotFoundException;
use Tesote\Sdk\Models\SyncSession;
use Tesote\Sdk\Models\SyncSessionList;
use Tesote\Sdk\Tests\Support\TestCaseBase;
use Tesote\Sdk\V2\Client;

final class SyncSessionsTest extends TestCaseBase
{
    public function testListForAccount(): void
    {
        $this->enqueueOk([
            'sync_sessions' => [
                [
                    'id' => 'ss-1',
                    'status' => 'completed',
                    'started_at' => '2026-04-28T10:00:00Z',
                    'completed_at' => '2026-04-28T10:00:30Z',
                    'transactions_synced' => 5,
                    'accounts_count' => 1,
                    'error' => null,
                    'performance' => [
                        'total_duration' => 30.5,
                        'complexity_score' => 1.0,
                        'sync_speed_score' => 0.5,
                    ],
                ],
            ],
            'limit' => 50,
            'offset' => 0,
            'has_more' => false,
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $list = $client->syncSessions->listForAccount('acct-1', ['limit' => 50, 'status' => 'completed']);

        self::assertInstanceOf(SyncSessionList::class, $list);
        self::assertCount(1, $list->syncSessions);
        self::assertNotNull($list->syncSessions[0]->performance);
        self::assertSame(30.5, $list->syncSessions[0]->performance['total_duration']);
        self::assertStringContainsString('limit=50&status=completed', $this->lastUrl());
    }

    public function testGet(): void
    {
        $this->enqueueOk([
            'id' => 'ss-1',
            'status' => 'failed',
            'started_at' => 't',
            'completed_at' => null,
            'transactions_synced' => 0,
            'accounts_count' => 1,
            'error' => ['type' => 'BankError', 'message' => 'down'],
            'performance' => null,
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $session = $client->syncSessions->get('acct-1', 'ss-1');

        self::assertInstanceOf(SyncSession::class, $session);
        self::assertSame('failed', $session->status);
        self::assertNotNull($session->error);
        self::assertSame('BankError', $session->error['type']);
    }

    public function testGetMapsNotFound(): void
    {
        $this->enqueueError(404, ['error' => 'gone', 'error_code' => 'SYNC_SESSION_NOT_FOUND']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(SyncSessionNotFoundException::class);
        $client->syncSessions->get('acct-1', 'missing');
    }
}
