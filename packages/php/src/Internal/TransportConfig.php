<?php

declare(strict_types=1);

namespace Tesote\Sdk\Internal;

use Tesote\Sdk\Cache\CacheBackend;
use Tesote\Sdk\Errors\ConfigException;
use Tesote\Sdk\Http\CurlInterface;
use Tesote\Sdk\Http\ExtCurl;
use Tesote\Sdk\Transport;

/**
 * Validated, defaulted view of the Transport constructor config array.
 *
 * Pulls the array-shape juggling out of Transport so the orchestrator only
 * sees typed properties.
 */
final class TransportConfig
{
    public readonly string $apiKey;
    public readonly string $baseUrl;
    public readonly string $userAgent;
    public readonly int $maxAttempts;
    public readonly int $baseDelayMs;
    public readonly int $maxDelayMs;
    public readonly int $connectTimeoutMs;
    public readonly int $timeoutMs;
    public readonly ?CacheBackend $cache;
    public readonly CurlInterface $curl;
    /** @var (callable(string, string, array<string, mixed>): void)|null */
    public $logger;
    /** @var (callable(int): void)|null */
    public $sleeper;

    /**
     * @param array{
     *     apiKey?: string,
     *     baseUrl?: string,
     *     userAgent?: string,
     *     maxAttempts?: int,
     *     baseDelayMs?: int,
     *     maxDelayMs?: int,
     *     connectTimeoutMs?: int,
     *     timeoutMs?: int,
     *     cache?: CacheBackend|null,
     *     curl?: CurlInterface|null,
     *     logger?: callable|null,
     *     sleeper?: callable|null,
     * } $config
     */
    public function __construct(array $config)
    {
        $apiKey = $config['apiKey'] ?? '';
        if (!is_string($apiKey) || $apiKey === '') {
            throw new ConfigException('apiKey is required and must be a non-empty string.');
        }
        $this->apiKey = $apiKey;
        $this->baseUrl = rtrim((string) ($config['baseUrl'] ?? Transport::DEFAULT_BASE_URL), '/');
        $this->userAgent = (string) ($config['userAgent'] ?? sprintf(
            'tesote-sdk-php/%s (php/%s)',
            Transport::VERSION,
            PHP_VERSION,
        ));
        $this->maxAttempts = max(1, (int) ($config['maxAttempts'] ?? 3));
        $this->baseDelayMs = max(1, (int) ($config['baseDelayMs'] ?? 250));
        $this->maxDelayMs = max($this->baseDelayMs, (int) ($config['maxDelayMs'] ?? 8000));
        $this->connectTimeoutMs = max(1, (int) ($config['connectTimeoutMs'] ?? 5000));
        $this->timeoutMs = max(1, (int) ($config['timeoutMs'] ?? 30000));
        $this->cache = $config['cache'] ?? null;
        $this->curl = $config['curl'] ?? new ExtCurl();
        $this->logger = $config['logger'] ?? null;
        $this->sleeper = $config['sleeper'] ?? null;
    }
}
