<?php

declare(strict_types=1);

namespace Tesote\Sdk\Http;

use Tesote\Sdk\Errors\NetworkException;
use Tesote\Sdk\Errors\TimeoutException;
use Tesote\Sdk\Errors\TlsException;
use Tesote\Sdk\Errors\TransportException;

/**
 * Maps a libcurl errno + message into the right typed TransportException.
 *
 * Lives outside Transport so the classification table can grow without
 * dragging the transport file past the 500-LOC guideline.
 */
final class CurlErrorClassifier
{
    /**
     * @param array<string, mixed>|null $requestSummary
     */
    public static function classify(CurlResult $result, ?array $requestSummary, int $attempts): TransportException
    {
        if ($result->errno === CURLE_OPERATION_TIMEOUTED) {
            return self::build(TimeoutException::class, 'TIMEOUT', $result, $requestSummary, $attempts, 'Request timed out');
        }
        if (self::isTlsErrno($result->errno)) {
            return self::build(TlsException::class, 'TLS', $result, $requestSummary, $attempts, 'TLS handshake failed');
        }
        return self::build(NetworkException::class, 'NETWORK', $result, $requestSummary, $attempts, 'Network error');
    }

    public static function isRetriableErrno(int $errno): bool
    {
        return in_array($errno, [
            CURLE_COULDNT_RESOLVE_HOST,
            CURLE_COULDNT_CONNECT,
            CURLE_OPERATION_TIMEOUTED,
            CURLE_GOT_NOTHING,
            CURLE_RECV_ERROR,
            CURLE_SEND_ERROR,
        ], true);
    }

    private static function isTlsErrno(int $errno): bool
    {
        // why: some libcurl builds omit certain TLS error constants; guard each one with defined().
        $tlsConstants = ['CURLE_SSL_CONNECT_ERROR', 'CURLE_PEER_FAILED_VERIFICATION', 'CURLE_SSL_CERTPROBLEM', 'CURLE_SSL_CACERT'];
        foreach ($tlsConstants as $name) {
            if (defined($name) && constant($name) === $errno) {
                return true;
            }
        }
        return false;
    }

    /**
     * @param class-string<TransportException> $class
     * @param array<string, mixed>|null        $requestSummary
     */
    private static function build(
        string $class,
        string $code,
        CurlResult $result,
        ?array $requestSummary,
        int $attempts,
        string $fallbackMessage,
    ): TransportException {
        return new $class(
            $result->errorMessage !== '' ? $result->errorMessage : $fallbackMessage,
            $code,
            0,
            null,
            null,
            null,
            null,
            $requestSummary,
            $attempts,
        );
    }
}
