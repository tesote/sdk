<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V1;

use Tesote\Sdk\Models\StatusInfo;
use Tesote\Sdk\Models\WhoAmI;
use Tesote\Sdk\Tests\Support\TestCaseBase;
use Tesote\Sdk\V1\Client;

final class StatusTest extends TestCaseBase
{
    public function testStatusReturnsModel(): void
    {
        $this->enqueueOk(['status' => 'ok', 'authenticated' => false]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $info = $client->status->check();

        self::assertInstanceOf(StatusInfo::class, $info);
        self::assertSame('ok', $info->status);
        self::assertFalse($info->authenticated);
        self::assertStringContainsString('/status', $this->lastUrl());
    }

    public function testWhoamiReturnsModel(): void
    {
        $this->enqueueOk([
            'client' => ['id' => 'ws-1', 'name' => 'ACME', 'type' => 'workspace'],
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $who = $client->status->whoami();

        self::assertInstanceOf(WhoAmI::class, $who);
        self::assertSame('workspace', $who->type);
        self::assertStringContainsString('/whoami', $this->lastUrl());
    }
}
