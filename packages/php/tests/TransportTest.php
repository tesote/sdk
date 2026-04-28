<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests;

use PHPUnit\Framework\TestCase;
use Tesote\Sdk\Cache\InMemoryCache;
use Tesote\Sdk\Errors\ConfigException;
use Tesote\Sdk\Errors\NetworkException;
use Tesote\Sdk\Errors\RateLimitExceededException;
use Tesote\Sdk\Errors\UnauthorizedException;
use Tesote\Sdk\Http\CurlResult;
use Tesote\Sdk\Tests\Support\FakeCurl;
use Tesote\Sdk\Transport;

final class TransportTest extends TestCase
{
    private FakeCurl $curl;
    /** @var list<int> */
    private array $sleeps = [];

    protected function setUp(): void
    {
        $this->curl = new FakeCurl();
        $this->sleeps = [];
    }

    public function testRequiresApiKey(): void
    {
        $this->expectException(ConfigException::class);
        new Transport(['apiKey' => '']);
    }

    public function testInjectsBearerAndDefaultHeaders(): void
    {
        $this->curl->enqueue(self::ok('{"ok":true}'));
        $transport = $this->makeTransport();
        $transport->request('GET', '/v3/accounts');

        $headers = $this->curl->calls[0]['headers'];
        self::assertSame('Bearer secret-key-1234', $headers['Authorization']);
        self::assertSame('application/json', $headers['Accept']);
        self::assertSame('application/json', $headers['Content-Type']);
        self::assertStringStartsWith('tesote-sdk-php/', $headers['User-Agent']);
    }

    public function testBuildsUrlWithBaseAndQuery(): void
    {
        $this->curl->enqueue(self::ok('{}'));
        $transport = $this->makeTransport();
        $transport->request('GET', '/v3/accounts', ['limit' => 25, 'cursor' => 'abc']);

        self::assertSame(
            'https://equipo.tesote.com/api/v3/accounts?limit=25&cursor=abc',
            $this->curl->calls[0]['options'][CURLOPT_URL],
        );
    }

    public function testCapturesRateLimitHeaders(): void
    {
        $this->curl->enqueue(new CurlResult(
            status: 200,
            body: '{"data":[]}',
            headers: [
                'x-ratelimit-limit' => '200',
                'x-ratelimit-remaining' => '199',
                'x-ratelimit-reset' => '60',
                'x-request-id' => 'req-1',
            ],
            errno: 0,
            errorMessage: '',
        ));
        $transport = $this->makeTransport();
        $transport->request('GET', '/v3/accounts');

        self::assertSame(['limit' => '200', 'remaining' => '199', 'reset' => '60'], $transport->getLastRateLimit());
        self::assertSame('req-1', $transport->getLastRequestId());
    }

    public function testRetriesOn503ThenSucceeds(): void
    {
        $this->curl->enqueue(new CurlResult(503, '{"error":"down","error_code":"SERVICE"}', [], 0, ''));
        $this->curl->enqueue(new CurlResult(503, '{"error":"down","error_code":"SERVICE"}', [], 0, ''));
        $this->curl->enqueue(self::ok('{"ok":true}'));

        $transport = $this->makeTransport();
        $body = $transport->request('GET', '/v3/accounts');

        self::assertSame(['ok' => true], $body);
        self::assertCount(3, $this->curl->calls);
        self::assertCount(2, $this->sleeps);
    }

    public function testHonoursRetryAfterOn429(): void
    {
        $this->curl->enqueue(new CurlResult(429, '{"error":"slow","error_code":"RATE_LIMIT_EXCEEDED"}', ['retry-after' => '2'], 0, ''));
        $this->curl->enqueue(self::ok('{"ok":true}'));

        $this->makeTransport(['maxAttempts' => 2])->request('GET', '/v3/accounts');

        self::assertSame([2000], $this->sleeps);
    }

    public function testRaisesRateLimitAfterRetriesExhausted(): void
    {
        for ($i = 0; $i < 3; $i++) {
            $this->curl->enqueue(new CurlResult(
                429,
                '{"error":"slow","error_code":"RATE_LIMIT_EXCEEDED","retry_after":42}',
                ['retry-after' => '42', 'x-request-id' => 'req-z'],
                0,
                '',
            ));
        }
        try {
            $this->makeTransport()->request('GET', '/v3/accounts');
            self::fail('expected RateLimitExceededException');
        } catch (RateLimitExceededException $e) {
            self::assertSame('RATE_LIMIT_EXCEEDED', $e->errorCode);
            self::assertSame(429, $e->httpStatus);
            self::assertSame(42, $e->retryAfter);
            self::assertSame('req-z', $e->requestId);
            self::assertSame(3, $e->attempts);
        }
    }

    public function testRaisesUnauthorizedOn401(): void
    {
        $this->curl->enqueue(new CurlResult(
            401,
            '{"error":"bad token","error_code":"UNAUTHORIZED"}',
            ['x-request-id' => 'req-u'],
            0,
            '',
        ));
        try {
            $this->makeTransport()->request('GET', '/v3/accounts');
            self::fail('expected UnauthorizedException');
        } catch (UnauthorizedException $e) {
            self::assertSame('UNAUTHORIZED', $e->errorCode);
            self::assertSame(1, $e->attempts);
            self::assertSame('req-u', $e->requestId);
        }
    }

    public function testDoesNotRetry4xxOtherThan429(): void
    {
        $this->curl->enqueue(new CurlResult(401, '{"error":"x","error_code":"UNAUTHORIZED"}', [], 0, ''));
        try {
            $this->makeTransport()->request('GET', '/v3/accounts');
            self::fail('expected exception');
        } catch (UnauthorizedException) {
            self::assertCount(1, $this->curl->calls);
        }
    }

    public function testGeneratesIdempotencyKeyForMutations(): void
    {
        $this->curl->enqueue(self::ok('{"ok":true}'));
        $this->makeTransport()->request('POST', '/v3/accounts/abc/sync', null, []);

        $key = $this->curl->calls[0]['headers']['Idempotency-Key'] ?? null;
        self::assertNotNull($key);
        self::assertMatchesRegularExpression(
            '/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/',
            $key,
        );
    }

    public function testPropagatesUserSuppliedIdempotencyKey(): void
    {
        $this->curl->enqueue(self::ok('{"ok":true}'));
        $this->makeTransport()->request('POST', '/v3/accounts/abc/sync', null, [], ['idempotencyKey' => 'caller-supplied-key']);
        self::assertSame('caller-supplied-key', $this->curl->calls[0]['headers']['Idempotency-Key']);
    }

    public function testNetworkErrorRaisedAfterRetries(): void
    {
        for ($i = 0; $i < 3; $i++) {
            $this->curl->enqueue(new CurlResult(0, '', [], CURLE_COULDNT_RESOLVE_HOST, 'dns failure'));
        }
        try {
            $this->makeTransport()->request('GET', '/v3/accounts');
            self::fail('expected NetworkException');
        } catch (NetworkException $e) {
            self::assertSame(3, $e->attempts);
            self::assertSame('dns failure', $e->getMessage());
        }
    }

    public function testTimeoutRetriedWhenIdempotencyKeyPresent(): void
    {
        // why: with an idempotency key the request is replay-safe, so timeouts retry up to maxAttempts.
        $this->curl->enqueue(new CurlResult(0, '', [], CURLE_OPERATION_TIMEOUTED, 'timed out'));
        $this->curl->enqueue(self::ok('{"ok":true}'));
        $body = $this->makeTransport(['maxAttempts' => 2])
            ->request('POST', '/v3/accounts/abc/sync', null, [], ['idempotencyKey' => 'caller-key']);
        self::assertSame(['ok' => true], $body);
        self::assertCount(2, $this->curl->calls);
    }

    public function testRequestSummaryRedactsBearer(): void
    {
        $this->curl->enqueue(new CurlResult(
            422,
            '{"error":"bad","error_code":"UNPROCESSABLE_CONTENT"}',
            [],
            0,
            '',
        ));
        try {
            $this->makeTransport()->request('POST', '/v3/accounts', ['x' => '1'], ['name' => 'x']);
            self::fail('expected exception');
        } catch (\Tesote\Sdk\Errors\UnprocessableContentException $e) {
            $summary = $e->requestSummary;
            self::assertNotNull($summary);
            self::assertSame('Bearer 1234', $summary['auth']);
            self::assertStringNotContainsString('secret-key-1234', json_encode($summary, JSON_THROW_ON_ERROR) ?: '');
        }
    }

    public function testCachesGetResponses(): void
    {
        $this->curl->enqueue(self::ok('{"data":[1,2,3]}'));
        $cache = new InMemoryCache();
        $transport = $this->makeTransport(['cache' => $cache]);

        $first = $transport->request('GET', '/v3/accounts', null, null, ['cacheTtl' => 30]);
        $second = $transport->request('GET', '/v3/accounts', null, null, ['cacheTtl' => 30]);

        self::assertSame($first, $second);
        self::assertCount(1, $this->curl->calls);
    }

    public function testCacheSkippedForMutations(): void
    {
        $this->curl->enqueue(self::ok('{"ok":true}'));
        $this->curl->enqueue(self::ok('{"ok":true}'));
        $cache = new InMemoryCache();
        $transport = $this->makeTransport(['cache' => $cache]);

        $transport->request('POST', '/v3/accounts/abc/sync', null, [], ['cacheTtl' => 30]);
        $transport->request('POST', '/v3/accounts/abc/sync', null, [], ['cacheTtl' => 30]);
        self::assertCount(2, $this->curl->calls);
    }

    public function testGenerateUuidV4Format(): void
    {
        for ($i = 0; $i < 25; $i++) {
            $uuid = Transport::generateUuidV4();
            self::assertMatchesRegularExpression(
                '/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/',
                $uuid,
            );
        }
    }

    public function testTimeoutsConfigured(): void
    {
        $this->curl->enqueue(self::ok('{}'));
        $this->makeTransport(['connectTimeoutMs' => 1234, 'timeoutMs' => 5678])->request('GET', '/v3/accounts');
        $opts = $this->curl->calls[0]['options'];
        self::assertSame(1234, $opts[CURLOPT_CONNECTTIMEOUT_MS]);
        self::assertSame(5678, $opts[CURLOPT_TIMEOUT_MS]);
    }

    /**
     * @param array<string, mixed> $overrides
     */
    private function makeTransport(array $overrides = []): Transport
    {
        $config = array_merge([
            'apiKey' => 'secret-key-1234',
            'curl' => $this->curl,
            'sleeper' => function (int $ms): void {
                $this->sleeps[] = $ms;
            },
        ], $overrides);
        return new Transport($config);
    }

    private static function ok(string $body): CurlResult
    {
        return new CurlResult(200, $body, ['x-request-id' => 'req-ok'], 0, '');
    }
}
