<?php

declare(strict_types=1);

namespace Tesote\Sdk\Http;

/**
 * Production CurlInterface backed by ext-curl.
 *
 * Captures response headers via CURLOPT_HEADERFUNCTION because PHP's stock
 * include-header-in-body mode forces brittle string parsing, especially when
 * intermediaries inject extra status lines (100-continue, redirects).
 */
final class ExtCurl implements CurlInterface
{
    public function execute(array $options): CurlResult
    {
        $handle = curl_init();
        if ($handle === false) {
            // why: curl_init() can fail under extreme resource pressure; caller treats this as a NetworkError upstream.
            return new CurlResult(0, '', [], 1, 'curl_init failed');
        }

        $headers = [];
        $options[CURLOPT_HEADERFUNCTION] = static function ($_, string $headerLine) use (&$headers): int {
            $trim = trim($headerLine);
            if ($trim === '' || stripos($trim, 'HTTP/') === 0) {
                return strlen($headerLine);
            }
            $colon = strpos($trim, ':');
            if ($colon !== false) {
                $name = strtolower(trim(substr($trim, 0, $colon)));
                $value = trim(substr($trim, $colon + 1));
                $headers[$name] = $value;
            }
            return strlen($headerLine);
        };
        $options[CURLOPT_RETURNTRANSFER] = true;

        curl_setopt_array($handle, $options);

        /** @var string|false $body */
        $body = curl_exec($handle);
        $status = (int) curl_getinfo($handle, CURLINFO_RESPONSE_CODE);
        $errno = curl_errno($handle);
        $errorMessage = curl_error($handle);
        curl_close($handle);

        return new CurlResult(
            status: $status,
            body: $body === false ? '' : $body,
            headers: $headers,
            errno: $errno,
            errorMessage: $errorMessage,
        );
    }
}
