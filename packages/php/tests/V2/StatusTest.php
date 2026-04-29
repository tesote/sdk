<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V2;

use Tesote\Sdk\Errors\UnauthorizedException;
use Tesote\Sdk\Models\StatusInfo;
use Tesote\Sdk\Models\WhoAmI;
use Tesote\Sdk\Tests\Support\TestCaseBase;
use Tesote\Sdk\V2\Client;

final class StatusTest extends TestCaseBase
{
    public function testStatus(): void
    {
        $this->enqueueOk(['status' => 'ok', 'authenticated' => false]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $info = $client->status->check();

        self::assertInstanceOf(StatusInfo::class, $info);
        self::assertStringContainsString('/v2/status', $this->lastUrl());
    }

    public function testWhoami(): void
    {
        $this->enqueueOk(['client' => ['id' => 'u-1', 'name' => 'User', 'type' => 'user']]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $who = $client->status->whoami();

        self::assertInstanceOf(WhoAmI::class, $who);
        self::assertSame('user', $who->type);
    }

    public function testWhoamiUnauthorized(): void
    {
        $this->enqueueError(401, ['error' => 'no', 'error_code' => 'UNAUTHORIZED']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(UnauthorizedException::class);
        $client->status->whoami();
    }
}
