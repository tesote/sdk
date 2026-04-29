<?php

declare(strict_types=1);

namespace Tesote\Sdk\Internal;

/**
 * Builds the URL, header set, and curl option array for a single request.
 *
 * Stateless aside from constructor-bound configuration. The Transport owns
 * the orchestration loop; this class owns wire-format assembly.
 */
final class RequestBuilder
{
    public function __construct(
        private readonly string $baseUrl,
        private readonly string $apiKey,
        private readonly string $userAgent,
        private readonly int $connectTimeoutMs,
        private readonly int $timeoutMs,
    ) {
    }

    /**
     * @param array<string, scalar|array<int|string, scalar>>|null $query
     */
    public function buildUrl(string $path, ?array $query): string
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
    public function defaultHeaders(array $extra): array
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
    public function buildCurlOptions(string $method, string $url, array $headers, ?string $encodedBody): array
    {
        $opts = [
            CURLOPT_URL => $url,
            CURLOPT_CUSTOMREQUEST => $method,
            CURLOPT_FOLLOWLOCATION => false,
            CURLOPT_CONNECTTIMEOUT_MS => $this->connectTimeoutMs,
            CURLOPT_TIMEOUT_MS => $this->timeoutMs,
            CURLOPT_HTTPHEADER => self::flattenHeaders($headers),
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
    private static function flattenHeaders(array $headers): array
    {
        $out = [];
        foreach ($headers as $name => $value) {
            $out[] = $name . ': ' . $value;
        }
        return $out;
    }
}
