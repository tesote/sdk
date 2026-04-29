<?php

declare(strict_types=1);

namespace Tesote\Sdk;

use Tesote\Sdk\Cache\CacheBackend;
use Tesote\Sdk\Errors\ApiException;
use Tesote\Sdk\Http\CurlErrorClassifier;
use Tesote\Sdk\Http\CurlInterface;
use Tesote\Sdk\Internal\BackoffStrategy;
use Tesote\Sdk\Internal\CacheKey;
use Tesote\Sdk\Internal\JsonCodec;
use Tesote\Sdk\Internal\RateLimitTracker;
use Tesote\Sdk\Internal\RequestBuilder;
use Tesote\Sdk\Internal\RequestSummarizer;
use Tesote\Sdk\Internal\RetryPolicy;
use Tesote\Sdk\Internal\TransportConfig;
use Tesote\Sdk\Internal\Uuid;

/**
 * Single HTTP client for the SDK. Resource clients call request() and get
 * back a parsed JSON array (or null for 204), or a typed exception.
 *
 * Owns: bearer injection, retries with exp-backoff + jitter, rate-limit
 * header capture, idempotency-key auto-gen for mutations, request-id
 * propagation into exceptions, optional TTL cache. Cross-cutting concerns
 * live in single-purpose collaborators in Tesote\Sdk\Internal\*.
 */
final class Transport
{
    public const VERSION = '0.2.0';
    public const DEFAULT_BASE_URL = 'https://equipo.tesote.com/api';

    private readonly string $apiKey;
    private readonly int $maxAttempts;
    private readonly ?CacheBackend $cache;
    private readonly CurlInterface $curl;
    /** @var (callable(string, string, array<string, mixed>): void)|null */
    private $logger;

    private readonly RequestBuilder $builder;
    private readonly BackoffStrategy $backoff;
    private readonly RequestSummarizer $summarizer;
    private readonly RateLimitTracker $rateLimit;

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
        $cfg = new TransportConfig($config);
        $this->apiKey = $cfg->apiKey;
        $this->maxAttempts = $cfg->maxAttempts;
        $this->cache = $cfg->cache;
        $this->curl = $cfg->curl;
        $this->logger = $cfg->logger;

        $this->builder = new RequestBuilder(
            baseUrl: $cfg->baseUrl,
            apiKey: $cfg->apiKey,
            userAgent: $cfg->userAgent,
            connectTimeoutMs: $cfg->connectTimeoutMs,
            timeoutMs: $cfg->timeoutMs,
        );
        $this->backoff = new BackoffStrategy($cfg->baseDelayMs, $cfg->maxDelayMs, $cfg->sleeper);
        $this->summarizer = new RequestSummarizer($cfg->apiKey);
        $this->rateLimit = new RateLimitTracker();
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
        $url = $this->builder->buildUrl($path, $query);
        $isMutation = RetryPolicy::isMutating($method);
        $cacheKey = CacheKey::for($method, $url, $this->apiKey);
        $cacheTtl = isset($opts['cacheTtl']) ? (int) $opts['cacheTtl'] : 0;

        if (!$isMutation && $this->cache !== null && $cacheTtl > 0) {
            $hit = $this->cache->get($cacheKey);
            if ($hit !== null) {
                $this->rateLimit->record($hit['headers']);
                return is_array($hit['body']) ? $hit['body'] : null;
            }
        }

        $idempotencyKey = $opts['idempotencyKey'] ?? null;
        if ($isMutation && $idempotencyKey === null) {
            $idempotencyKey = Uuid::v4();
        }

        $headers = $this->builder->defaultHeaders($opts['headers'] ?? []);
        if ($idempotencyKey !== null) {
            $headers['Idempotency-Key'] = $idempotencyKey;
        }

        $encodedBody = $body !== null ? JsonCodec::encode($body) : null;
        $summary = $this->summarizer->summarise($method, $path, $query, $body);
        $lastException = null;

        for ($attempt = 1; $attempt <= $this->maxAttempts; $attempt++) {
            $options = $this->builder->buildCurlOptions($method, $url, $headers, $encodedBody);
            $result = $this->curl->execute($options);
            $this->log($method, $url, ['status' => $result->status, 'attempt' => $attempt]);

            if ($result->errno !== 0) {
                $exception = CurlErrorClassifier::classify($result, $summary, $attempt);
                $lastException = $exception;
                if (RetryPolicy::shouldRetryTransport($result, $method, $idempotencyKey) && $attempt < $this->maxAttempts) {
                    $this->backoff->sleep($this->backoff->delayMs($attempt, null));
                    continue;
                }
                throw $exception;
            }

            $this->rateLimit->record($result->headers);

            if ($result->status >= 200 && $result->status < 300) {
                $decoded = JsonCodec::decode($result->body);
                if (!$isMutation && $this->cache !== null && $cacheTtl > 0) {
                    $this->cache->set($cacheKey, [
                        'body' => $decoded,
                        'headers' => $result->headers,
                        'status' => $result->status,
                    ], $cacheTtl);
                } elseif ($isMutation && $this->cache !== null) {
                    // why: any mutation on the resource path should bust the matching GET cache entry.
                    $this->cache->delete(CacheKey::for('GET', $url, $this->apiKey));
                }
                return $decoded;
            }

            if (RetryPolicy::isRetriableStatus($result->status) && $attempt < $this->maxAttempts) {
                $delay = $this->backoff->delayMs($attempt, RetryPolicy::parseRetryAfter($result->headers));
                $this->backoff->sleep($delay);
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
     * Like request() but returns the raw body string + response headers.
     *
     * Used by export endpoints that return CSV / pretty-printed JSON files
     * rather than the standard JSON envelope. Skips the response cache.
     *
     * @param array<string, scalar|array<int|string, scalar>>|null $query
     * @param array<mixed>|null                                    $body
     * @param array{
     *     idempotencyKey?: string|null,
     *     headers?: array<string, string>,
     * } $opts
     *
     * @return array{body: string, headers: array<string, string>, status: int}
     */
    public function requestRaw(
        string $method,
        string $path,
        ?array $query = null,
        ?array $body = null,
        array $opts = [],
    ): array {
        $method = strtoupper($method);
        $url = $this->builder->buildUrl($path, $query);
        $isMutation = RetryPolicy::isMutating($method);

        $idempotencyKey = $opts['idempotencyKey'] ?? null;
        if ($isMutation && $idempotencyKey === null) {
            $idempotencyKey = Uuid::v4();
        }

        $headers = $this->builder->defaultHeaders($opts['headers'] ?? []);
        if ($idempotencyKey !== null) {
            $headers['Idempotency-Key'] = $idempotencyKey;
        }

        $encodedBody = $body !== null ? JsonCodec::encode($body) : null;
        $summary = $this->summarizer->summarise($method, $path, $query, $body);
        $lastException = null;

        for ($attempt = 1; $attempt <= $this->maxAttempts; $attempt++) {
            $options = $this->builder->buildCurlOptions($method, $url, $headers, $encodedBody);
            $result = $this->curl->execute($options);
            $this->log($method, $url, ['status' => $result->status, 'attempt' => $attempt]);

            if ($result->errno !== 0) {
                $exception = CurlErrorClassifier::classify($result, $summary, $attempt);
                $lastException = $exception;
                if (RetryPolicy::shouldRetryTransport($result, $method, $idempotencyKey) && $attempt < $this->maxAttempts) {
                    $this->backoff->sleep($this->backoff->delayMs($attempt, null));
                    continue;
                }
                throw $exception;
            }

            $this->rateLimit->record($result->headers);

            if ($result->status >= 200 && $result->status < 300) {
                return [
                    'body' => $result->body,
                    'headers' => $result->headers,
                    'status' => $result->status,
                ];
            }

            if (RetryPolicy::isRetriableStatus($result->status) && $attempt < $this->maxAttempts) {
                $delay = $this->backoff->delayMs($attempt, RetryPolicy::parseRetryAfter($result->headers));
                $this->backoff->sleep($delay);
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
        return $this->rateLimit->lastRateLimit();
    }

    public function getLastRequestId(): ?string
    {
        return $this->rateLimit->lastRequestId();
    }

    public static function generateUuidV4(): string
    {
        return Uuid::v4();
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
