<?php

declare(strict_types=1);

namespace Tesote\Sdk\Internal;

/**
 * Builds the redacted request summary attached to thrown exceptions.
 *
 * The bearer token must never appear in the summary; only the last four
 * characters are exposed so support tickets can correlate without leaking.
 */
final class RequestSummarizer
{
    public function __construct(private readonly string $apiKey)
    {
    }

    /**
     * @param array<string, scalar|array<int|string, scalar>>|null $query
     * @param array<mixed>|null                                    $body
     *
     * @return array<string, mixed>
     */
    public function summarise(string $method, string $path, ?array $query, ?array $body): array
    {
        return [
            'method' => $method,
            'path' => $path,
            'query' => $query,
            'bodyShape' => $body !== null ? $this->describeBody($body) : null,
            'auth' => 'Bearer ' . $this->lastFour(),
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

    private function lastFour(): string
    {
        return strlen($this->apiKey) <= 4 ? '****' : substr($this->apiKey, -4);
    }
}
