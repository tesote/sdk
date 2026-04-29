<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\Support;

use PHPUnit\Framework\TestCase;
use Tesote\Sdk\Http\CurlResult;
use Tesote\Sdk\Transport;

/**
 * Shared boilerplate for resource-client tests.
 *
 * Each test gets a fresh FakeCurl + Transport + recorded sleeps. Helpers below
 * cover the success / typed-error / 415 / cursor / idempotency cases without
 * forcing every per-resource test to repeat the wiring.
 */
abstract class TestCaseBase extends TestCase
{
    protected FakeCurl $curl;

    /** @var list<int> */
    protected array $sleeps = [];

    protected function setUp(): void
    {
        $this->curl = new FakeCurl();
        $this->sleeps = [];
    }

    /**
     * @param array<string, mixed> $overrides
     */
    protected function makeTransport(array $overrides = []): Transport
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

    /**
     * @param array<string, mixed> $body
     * @param array<string, string> $headers
     */
    protected function enqueueOk(array $body, int $status = 200, array $headers = []): void
    {
        $payload = $body === [] ? '{}' : json_encode($body, JSON_THROW_ON_ERROR);
        $this->curl->enqueue(new CurlResult(
            status: $status,
            body: $payload,
            headers: array_merge(['x-request-id' => 'req-ok'], $headers),
            errno: 0,
            errorMessage: '',
        ));
    }

    /**
     * @param array<string, string> $headers
     */
    protected function enqueueRaw(string $body, int $status = 200, array $headers = []): void
    {
        $this->curl->enqueue(new CurlResult(
            status: $status,
            body: $body,
            headers: array_merge(['x-request-id' => 'req-ok'], $headers),
            errno: 0,
            errorMessage: '',
        ));
    }

    /**
     * @param array<string, mixed>  $envelope
     * @param array<string, string> $headers
     */
    protected function enqueueError(int $status, array $envelope, array $headers = []): void
    {
        $this->curl->enqueue(new CurlResult(
            status: $status,
            body: json_encode($envelope, JSON_THROW_ON_ERROR),
            headers: array_merge(['x-request-id' => 'req-err'], $headers),
            errno: 0,
            errorMessage: '',
        ));
    }

    protected function lastUrl(): string
    {
        $call = end($this->curl->calls);
        if ($call === false) {
            self::fail('No requests captured.');
        }
        return (string) $call['options'][CURLOPT_URL];
    }

    protected function lastMethod(): string
    {
        $call = end($this->curl->calls);
        if ($call === false) {
            self::fail('No requests captured.');
        }
        return (string) $call['options'][CURLOPT_CUSTOMREQUEST];
    }

    /**
     * @return array<string, string>
     */
    protected function lastHeaders(): array
    {
        $call = end($this->curl->calls);
        if ($call === false) {
            self::fail('No requests captured.');
        }
        return $call['headers'];
    }

    protected function lastBody(): ?string
    {
        $call = end($this->curl->calls);
        if ($call === false) {
            self::fail('No requests captured.');
        }
        return isset($call['options'][CURLOPT_POSTFIELDS]) ? (string) $call['options'][CURLOPT_POSTFIELDS] : null;
    }
}
