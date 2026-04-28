<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\Support;

use Tesote\Sdk\Http\CurlInterface;
use Tesote\Sdk\Http\CurlResult;

/**
 * In-memory CurlInterface used by transport / resource tests.
 *
 * Pre-load a queue of responses with enqueue() and inspect what the
 * Transport sent via the recorded calls list. Throws if the queue runs dry,
 * which catches "transport made unexpected extra request" bugs.
 */
final class FakeCurl implements CurlInterface
{
    /** @var list<CurlResult> */
    private array $queue = [];

    /** @var list<array{options: array<int, mixed>, headers: array<string, string>}> */
    public array $calls = [];

    public function enqueue(CurlResult $result): void
    {
        $this->queue[] = $result;
    }

    public function execute(array $options): CurlResult
    {
        if ($this->queue === []) {
            throw new \RuntimeException('FakeCurl queue exhausted; transport sent an unexpected request.');
        }
        $headers = [];
        /** @var list<string> $rawHeaders */
        $rawHeaders = $options[CURLOPT_HTTPHEADER] ?? [];
        foreach ($rawHeaders as $line) {
            $colon = strpos($line, ':');
            if ($colon === false) {
                continue;
            }
            $headers[trim(substr($line, 0, $colon))] = trim(substr($line, $colon + 1));
        }
        $this->calls[] = [
            'options' => $options,
            'headers' => $headers,
        ];
        return array_shift($this->queue);
    }
}
