<?php

declare(strict_types=1);

namespace Tesote\Sdk;

use Tesote\Sdk\Cache\CacheBackend;
use Tesote\Sdk\Errors\ApiException;
use Tesote\Sdk\Errors\ConfigException;
use Tesote\Sdk\Http\CurlErrorClassifier;
use Tesote\Sdk\Http\CurlInterface;
use Tesote\Sdk\Http\CurlResult;
use Tesote\Sdk\Http\ExtCurl;

/**
 * Single HTTP client for the SDK. Resource clients call request() and get
 * back a parsed JSON array (or null for 204), or a typed exception.
 *
 * Owns: bearer injection, retries with exp-backoff + jitter, rate-limit
 * header capture, idempotency-key auto-gen for mutations, request-id
 * propagation into exceptions, optional TTL cache.
 *
 * Configuration is passed as a single associative array to mirror the
 * shape used by the language-specific Client constructors.
 */
final class Transport
{
    public const VERSION = '0.1.0';
    public const DEFAULT_BASE_URL = 'https://equipo.tesote.com/api';

    private readonly string $apiKey;
    private readonly string $baseUrl;
    private readonly string $userAgent;
    private readonly int $maxAttempts;
    private readonly int $baseDelayMs;
    private readonly int $maxDelayMs;
    private readonly int $connectTimeoutMs;
    private readonly int $timeoutMs;
    private readonly ?CacheBackend $cache;
    private readonly CurlInterface $curl;
    /** @var (callable(string, string, array<string, mixed>): void)|null */
    private $logger;
    /** @var (callable(int): void)|null */
    private $sleeper;

    /** @var array{limit: ?string, remaining: ?string, reset: ?string} */
    private array $lastRateLimit = [
        'limit' => null,
        'remaining' => null,
        'reset' => null,
    ];

    private ?string $lastRequestId = null;

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
        $this->baseUrl = rtrim((string) ($config['baseUrl'] ?? self::DEFAULT_BASE_URL), '/');
        $this->userAgent = (string) ($config['userAgent'] ?? sprintf(
            'tesote-sdk-php/%s (php/%s)',
            self::VERSION,
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

    /**
     * Issue a request and return decoded JSON (or null for 204 / empty body).
     *
     * @param string                                   $method  Uppercase HTTP verb.
     * @param string                                   $path    Path relative to base URL, including version segment.
     * @param array<string, scalar|array<int|string, scalar>>|null $query
     * @param array<mixed>|null                        $body    JSON-encodable structure.
     * @param array{
     *     idempotencyKey?: string|null,
     *     cacheTtl?: int|null,
     *     headers?: array<string, string>,
     * } $opts
     *
     * @return array<mixed>|null
     */
    public function request(
        string $method,
        string $path,
        ?array $query = null,
        ?array $body = null,
        array $opts = [],
    ): ?array {
        $method = strtoupper($method);
        $url = $this->buildUrl($path, $query);
        $isMutation = in_array($method, ['POST', 'PUT', 'PATCH', 'DELETE'], true);
        $cacheKey = $this->cacheKey($method, $url);
        $cacheTtl = isset($opts['cacheTtl']) ? (int) $opts['cacheTtl'] : 0;

        if (!$isMutation && $this->cache !== null && $cacheTtl > 0) {
            $hit = $this->cache->get($cacheKey);
            if ($hit !== null) {
                $this->captureRateLimit($hit['headers']);
                $this->lastRequestId = $hit['headers']['x-request-id'] ?? null;
                return is_array($hit['body']) ? $hit['body'] : null;
            }
        }

        $idempotencyKey = $opts['idempotencyKey'] ?? null;
        if ($isMutation && $idempotencyKey === null) {
            $idempotencyKey = self::generateUuidV4();
        }

        $headers = $this->defaultHeaders($opts['headers'] ?? []);
        if ($idempotencyKey !== null) {
            $headers['Idempotency-Key'] = $idempotencyKey;
        }

        $encodedBody = null;
        if ($body !== null) {
            try {
                $encodedBody = json_encode($body, JSON_THROW_ON_ERROR | JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE);
            } catch (\JsonException $e) {
                throw new ConfigException('Failed to JSON-encode request body: ' . $e->getMessage());
            }
        }

        $summary = $this->summarise($method, $path, $query, $body);
        $lastException = null;
        $lastResult = null;

        for ($attempt = 1; $attempt <= $this->maxAttempts; $attempt++) {
            $options = $this->buildCurlOptions($method, $url, $headers, $encodedBody);
            $result = $this->curl->execute($options);
            $lastResult = $result;
            $this->log($method, $url, ['status' => $result->status, 'attempt' => $attempt]);

            if ($result->errno !== 0) {
                $exception = CurlErrorClassifier::classify($result, $summary, $attempt);
                $lastException = $exception;
                if ($this->shouldRetryTransport($result, $method, $idempotencyKey) && $attempt < $this->maxAttempts) {
                    $this->sleep($this->backoffDelay($attempt, null));
                    continue;
                }
                throw $exception;
            }

            $this->captureRateLimit($result->headers);
            $this->lastRequestId = $result->headers['x-request-id'] ?? null;

            if ($result->status >= 200 && $result->status < 300) {
                $decoded = $this->decodeBody($result->body);
                if (!$isMutation && $this->cache !== null && $cacheTtl > 0) {
                    $this->cache->set($cacheKey, [
                        'body' => $decoded,
                        'headers' => $result->headers,
                        'status' => $result->status,
                    ], $cacheTtl);
                } elseif ($isMutation && $this->cache !== null) {
                    // why: any mutation on the resource path should bust the matching GET cache entry.
                    $this->cache->delete($this->cacheKey('GET', $url));
                }
                return $decoded;
            }

            if ($this->isRetriableStatus($result->status) && $attempt < $this->maxAttempts) {
                $delay = $this->backoffDelay($attempt, $this->parseRetryAfter($result->headers));
                $this->sleep($delay);
                continue;
            }

            throw ApiException::fromResponse(
                $result->status,
                $result->body,
                $result->headers,
                $summary,
                $attempt,
            );
        }

        // why: every loop iteration either returns, continues, or throws; this is unreachable.
        if ($lastException !== null) {
            throw $lastException;
        }
        throw new \LogicException('Transport loop exited without a result; should be unreachable');
    }

    /**
     * @return array{limit: ?string, remaining: ?string, reset: ?string}
     */
    public function getLastRateLimit(): array
    {
        return $this->lastRateLimit;
    }

    public function getLastRequestId(): ?string
    {
        return $this->lastRequestId;
    }

    public static function generateUuidV4(): string
    {
        $bytes = random_bytes(16);
        $bytes[6] = chr((ord($bytes[6]) & 0x0f) | 0x40);
        $bytes[8] = chr((ord($bytes[8]) & 0x3f) | 0x80);
        $hex = bin2hex($bytes);
        return sprintf(
            '%s-%s-%s-%s-%s',
            substr($hex, 0, 8),
            substr($hex, 8, 4),
            substr($hex, 12, 4),
            substr($hex, 16, 4),
            substr($hex, 20, 12),
        );
    }

    /**
     * @param array<string, scalar|array<int|string, scalar>>|null $query
     */
    private function buildUrl(string $path, ?array $query): string
    {
        $url = $this->baseUrl . '/' . ltrim($path, '/');
        if ($query !== null && $query !== []) {
            $url .= '?' . http_build_query($query, '', '&', PHP_QUERY_RFC3986);
        }
        return $url;
    }

    /**
     * @param  array<string, string> $extra
     * @return array<string, string>
     */
    private function defaultHeaders(array $extra): array
    {
        $headers = [
            'Authorization' => 'Bearer ' . $this->apiKey,
            'Accept' => 'application/json',
            'Content-Type' => 'application/json',
            'User-Agent' => $this->userAgent,
        ];
        foreach ($extra as $name => $value) {
            $headers[$name] = $value;
        }
        return $headers;
    }

    /**
     * @param  array<string, string> $headers
     * @return array<int, mixed>
     */
    private function buildCurlOptions(string $method, string $url, array $headers, ?string $encodedBody): array
    {
        $opts = [
            CURLOPT_URL => $url,
            CURLOPT_CUSTOMREQUEST => $method,
            CURLOPT_FOLLOWLOCATION => false,
            CURLOPT_CONNECTTIMEOUT_MS => $this->connectTimeoutMs,
            CURLOPT_TIMEOUT_MS => $this->timeoutMs,
            CURLOPT_HTTPHEADER => $this->flattenHeaders($headers),
        ];
        if ($encodedBody !== null) {
            $opts[CURLOPT_POSTFIELDS] = $encodedBody;
        }
        if ($method === 'HEAD') {
            $opts[CURLOPT_NOBODY] = true;
        }
        return $opts;
    }

    /**
     * @param  array<string, string> $headers
     * @return list<string>
     */
    private function flattenHeaders(array $headers): array
    {
        $out = [];
        foreach ($headers as $name => $value) {
            $out[] = $name . ': ' . $value;
        }
        return $out;
    }

    /**
     * @return array<mixed>|null
     */
    private function decodeBody(string $body): ?array
    {
        if ($body === '') {
            return null;
        }
        try {
            $decoded = json_decode($body, true, 512, JSON_THROW_ON_ERROR);
        } catch (\JsonException) {
            return null;
        }
        return is_array($decoded) ? $decoded : null;
    }

    /**
     * @param array<string, string> $headers
     */
    private function captureRateLimit(array $headers): void
    {
        $limit = $headers['x-ratelimit-limit'] ?? null;
        $remaining = $headers['x-ratelimit-remaining'] ?? null;
        $reset = $headers['x-ratelimit-reset'] ?? null;
        $this->lastRateLimit = [
            'limit' => $limit,
            'remaining' => $remaining,
            'reset' => $reset,
        ];
    }

    /**
     * @param array<string, string> $headers
     */
    private function parseRetryAfter(array $headers): ?int
    {
        if (isset($headers['retry-after']) && is_numeric($headers['retry-after'])) {
            return (int) $headers['retry-after'];
        }
        return null;
    }

    private function isRetriableStatus(int $status): bool
    {
        return $status === 429 || $status === 502 || $status === 503 || $status === 504;
    }

    private function shouldRetryTransport(CurlResult $result, string $method, ?string $idempotencyKey): bool
    {
        // why: docs/architecture/transport.md — never retry POST without idempotency key on read timeout.
        $isReadTimeout = $result->errno === CURLE_OPERATION_TIMEOUTED;
        $isMutating = in_array($method, ['POST', 'PUT', 'PATCH', 'DELETE'], true);
        if ($isReadTimeout && $isMutating && $idempotencyKey === null) {
            return false;
        }
        return CurlErrorClassifier::isRetriableErrno($result->errno);
    }

    private function backoffDelay(int $attempt, ?int $retryAfterSeconds): int
    {
        if ($retryAfterSeconds !== null) {
            return min($this->maxDelayMs, $retryAfterSeconds * 1000);
        }
        $expo = $this->baseDelayMs * (2 ** ($attempt - 1));
        $capped = min($this->maxDelayMs, $expo);
        // why: full jitter — distributes retries to avoid thundering herd.
        return random_int(0, max(1, (int) $capped));
    }

    private function sleep(int $milliseconds): void
    {
        if ($milliseconds <= 0) {
            return;
        }
        if ($this->sleeper !== null) {
            ($this->sleeper)($milliseconds);
            return;
        }
        usleep($milliseconds * 1000);
    }

    private function cacheKey(string $method, string $url): string
    {
        return hash('sha256', $method . ' ' . $url . ' ' . substr($this->apiKey, -4));
    }

    /**
     * @param array<string, scalar|array<int|string, scalar>>|null $query
     * @param array<mixed>|null                                    $body
     *
     * @return array<string, mixed>
     */
    private function summarise(string $method, string $path, ?array $query, ?array $body): array
    {
        return [
            'method' => $method,
            'path' => $path,
            'query' => $query,
            'bodyShape' => $body !== null ? $this->describeBody($body) : null,
            'auth' => 'Bearer ' . $this->lastFour($this->apiKey),
        ];
    }

    /**
     * @param  array<mixed> $body
     * @return array<string, int|string>
     */
    private function describeBody(array $body): array
    {
        return [
            'keys' => count($body),
            'type' => array_is_list($body) ? 'list' : 'object',
        ];
    }

    private function lastFour(string $apiKey): string
    {
        return strlen($apiKey) <= 4 ? '****' : substr($apiKey, -4);
    }

    /**
     * @param array<string, mixed> $context
     */
    private function log(string $method, string $url, array $context): void
    {
        if ($this->logger === null) {
            return;
        }
        ($this->logger)($method, $url, $context);
    }
}
