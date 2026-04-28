<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V3;

use PHPUnit\Framework\TestCase;
use Tesote\Sdk\Http\CurlResult;
use Tesote\Sdk\Tests\Support\FakeCurl;
use Tesote\Sdk\Transport;
use Tesote\Sdk\V3\Client;

final class AccountsTest extends TestCase
{
    public function testListIssuesGetToVersionedPath(): void
    {
        $curl = new FakeCurl();
        $curl->enqueue(self::ok('{"data":[{"id":"acct_1"}]}'));

        $client = new Client(['apiKey' => 'k1234', 'curl' => $curl]);
        $body = $client->accounts->list(['limit' => 10]);

        self::assertSame(['data' => [['id' => 'acct_1']]], $body);
        self::assertSame('GET', $curl->calls[0]['options'][CURLOPT_CUSTOMREQUEST]);
        self::assertSame(
            'https://equipo.tesote.com/api/v3/accounts?limit=10',
            $curl->calls[0]['options'][CURLOPT_URL],
        );
    }

    public function testGetEncodesId(): void
    {
        $curl = new FakeCurl();
        $curl->enqueue(self::ok('{"id":"acct/with/slash"}'));

        $client = new Client(['apiKey' => 'k1234', 'curl' => $curl]);
        $body = $client->accounts->get('acct/with/slash');

        self::assertSame(['id' => 'acct/with/slash'], $body);
        self::assertSame(
            'https://equipo.tesote.com/api/v3/accounts/acct%2Fwith%2Fslash',
            $curl->calls[0]['options'][CURLOPT_URL],
        );
    }

    public function testStubResourceThrowsLogicException(): void
    {
        $client = new Client(['apiKey' => 'k1234', 'curl' => new FakeCurl()]);
        $this->expectException(\LogicException::class);
        $this->expectExceptionMessage('not implemented: webhooks.list()');
        /** @phpstan-ignore-next-line method.notFound */
        $client->webhooks->list();
    }

    public function testTransportIsShared(): void
    {
        $client = new Client(['apiKey' => 'k1234', 'curl' => new FakeCurl()]);
        self::assertInstanceOf(Transport::class, $client->transport);
    }

    private static function ok(string $body): CurlResult
    {
        return new CurlResult(200, $body, ['x-request-id' => 'req-acct'], 0, '');
    }
}
