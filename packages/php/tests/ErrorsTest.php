<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests;

use PHPUnit\Framework\Attributes\DataProvider;
use PHPUnit\Framework\TestCase;
use Tesote\Sdk\Errors\AccountDisabledException;
use Tesote\Sdk\Errors\ApiException;
use Tesote\Sdk\Errors\ApiKeyRevokedException;
use Tesote\Sdk\Errors\HistorySyncForbiddenException;
use Tesote\Sdk\Errors\InvalidDateRangeException;
use Tesote\Sdk\Errors\MutationDuringPaginationException;
use Tesote\Sdk\Errors\RateLimitExceededException;
use Tesote\Sdk\Errors\TesoteException;
use Tesote\Sdk\Errors\UnauthorizedException;
use Tesote\Sdk\Errors\UnprocessableContentException;
use Tesote\Sdk\Errors\WorkspaceSuspendedException;

final class ErrorsTest extends TestCase
{
    /**
     * @return array<string, array{int, string, class-string<TesoteException>}>
     */
    public static function errorCodeProvider(): array
    {
        return [
            'unauthorized'              => [401, 'UNAUTHORIZED',           UnauthorizedException::class],
            'api-key-revoked'           => [401, 'API_KEY_REVOKED',        ApiKeyRevokedException::class],
            'workspace-suspended'       => [403, 'WORKSPACE_SUSPENDED',    WorkspaceSuspendedException::class],
            'account-disabled'          => [403, 'ACCOUNT_DISABLED',       AccountDisabledException::class],
            'history-sync-forbidden'    => [403, 'HISTORY_SYNC_FORBIDDEN', HistorySyncForbiddenException::class],
            'mutation-conflict'         => [409, 'MUTATION_CONFLICT',      MutationDuringPaginationException::class],
            'unprocessable-content'     => [422, 'UNPROCESSABLE_CONTENT',  UnprocessableContentException::class],
            'invalid-date-range'        => [422, 'INVALID_DATE_RANGE',     InvalidDateRangeException::class],
            'rate-limit-exceeded'       => [429, 'RATE_LIMIT_EXCEEDED',    RateLimitExceededException::class],
        ];
    }

    /**
     * @param class-string<TesoteException> $expected
     */
    #[DataProvider('errorCodeProvider')]
    public function testErrorCodeMapsToTypedClass(int $status, string $errorCode, string $expected): void
    {
        $body = json_encode(['error' => 'msg', 'error_code' => $errorCode, 'error_id' => 'eid-1', 'retry_after' => 7], JSON_THROW_ON_ERROR);
        $exception = ApiException::fromResponse(
            $status,
            $body,
            ['x-request-id' => 'req-x', 'retry-after' => '11'],
            ['method' => 'GET', 'path' => '/x', 'auth' => 'Bearer 1234'],
            2,
        );

        self::assertInstanceOf($expected, $exception);
        self::assertSame($errorCode, $exception->errorCode);
        self::assertSame($status, $exception->httpStatus);
        self::assertSame('msg', $exception->getMessage());
        self::assertSame('eid-1', $exception->errorId);
        self::assertSame('req-x', $exception->requestId);
        self::assertSame(11, $exception->retryAfter, 'header takes priority over envelope');
        self::assertSame($body, $exception->responseBody);
        self::assertSame(2, $exception->attempts);
        self::assertNotNull($exception->requestSummary);
    }

    public function testUnknownErrorCodeFallsBackToBaseApiException(): void
    {
        $exception = ApiException::fromResponse(
            500,
            '{"error":"boom","error_code":"SOMETHING_NEW"}',
            [],
            null,
            1,
        );
        self::assertSame(ApiException::class, $exception::class);
        self::assertInstanceOf(TesoteException::class, $exception);
    }

    public function testEmptyBodyStillProducesUsableException(): void
    {
        $exception = ApiException::fromResponse(500, '', [], null, 1);
        self::assertSame('UNKNOWN', $exception->errorCode);
        self::assertSame('HTTP 500 (no error envelope)', $exception->getMessage());
    }

    public function testRetryAfterFromEnvelopeWhenHeaderAbsent(): void
    {
        $exception = ApiException::fromResponse(
            429,
            '{"error":"x","error_code":"RATE_LIMIT_EXCEEDED","retry_after":33}',
            [],
            null,
            1,
        );
        self::assertSame(33, $exception->retryAfter);
    }

    public function testSummaryFormat(): void
    {
        $exception = ApiException::fromResponse(
            429,
            '{"error":"slow","error_code":"RATE_LIMIT_EXCEEDED"}',
            ['x-request-id' => 'req-99', 'retry-after' => '5'],
            null,
            4,
        );
        $summary = $exception->summary();
        self::assertStringContainsString('RateLimitExceededException: 429 slow', $summary);
        self::assertStringContainsString('error_code=RATE_LIMIT_EXCEEDED', $summary);
        self::assertStringContainsString('request_id=req-99', $summary);
        self::assertStringContainsString('retry_after=5s', $summary);
        self::assertStringContainsString('attempts=4', $summary);
    }

    public function testRequestSummaryRedactsBearerToken(): void
    {
        $exception = ApiException::fromResponse(
            401,
            '{"error":"x","error_code":"UNAUTHORIZED"}',
            [],
            ['method' => 'POST', 'path' => '/v3/x', 'auth' => 'Bearer abcd'],
            1,
        );
        $serialized = json_encode($exception->requestSummary, JSON_THROW_ON_ERROR);
        self::assertIsString($serialized);
        self::assertStringContainsString('Bearer abcd', $serialized);
        self::assertStringNotContainsString('secret-real-key', $serialized);
    }
}
