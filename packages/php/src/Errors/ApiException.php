<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

use Throwable;

/**
 * Server-returned error: the API responded with a JSON error envelope.
 *
 * This is the parent of every error_code-specific class. Use the
 * fromResponse() factory to dispatch on error_code into the right subclass.
 */
class ApiException extends TesoteException
{
    /**
     * @param array<string, string>      $headers Lowercased response headers.
     * @param array<string, mixed>|null  $requestSummary
     */
    public static function fromResponse(
        int $status,
        string $rawBody,
        array $headers,
        ?array $requestSummary,
        int $attempts,
        ?Throwable $previous = null,
    ): self {
        $envelope = self::decodeEnvelope($rawBody);
        $errorCode = is_string($envelope['error_code'] ?? null)
            ? (string) $envelope['error_code']
            : 'UNKNOWN';
        $message = is_string($envelope['error'] ?? null)
            ? (string) $envelope['error']
            : sprintf('HTTP %d (no error envelope)', $status);
        $errorId = is_string($envelope['error_id'] ?? null)
            ? (string) $envelope['error_id']
            : null;
        $retryAfter = self::pickRetryAfter($envelope, $headers);
        $requestId = $headers['x-request-id'] ?? null;

        $class = self::dispatch($errorCode);

        return new $class(
            $message,
            $errorCode,
            $status,
            $requestId,
            $errorId,
            $retryAfter,
            $rawBody,
            $requestSummary,
            $attempts,
            $previous,
        );
    }

    /**
     * @return array<string, mixed>
     */
    private static function decodeEnvelope(string $rawBody): array
    {
        if ($rawBody === '') {
            return [];
        }
        try {
            $decoded = json_decode($rawBody, true, 16, JSON_THROW_ON_ERROR);
        } catch (\JsonException) {
            return [];
        }
        return is_array($decoded) ? $decoded : [];
    }

    /**
     * @param array<string, mixed>  $envelope
     * @param array<string, string> $headers
     */
    private static function pickRetryAfter(array $envelope, array $headers): ?int
    {
        if (isset($headers['retry-after']) && is_numeric($headers['retry-after'])) {
            return (int) $headers['retry-after'];
        }
        if (isset($envelope['retry_after']) && is_numeric($envelope['retry_after'])) {
            return (int) $envelope['retry_after'];
        }
        return null;
    }

    /**
     * @return class-string<self>
     */
    private static function dispatch(string $errorCode): string
    {
        return match ($errorCode) {
            'UNAUTHORIZED' => UnauthorizedException::class,
            'API_KEY_REVOKED' => ApiKeyRevokedException::class,
            'WORKSPACE_SUSPENDED' => WorkspaceSuspendedException::class,
            'ACCOUNT_DISABLED' => AccountDisabledException::class,
            'HISTORY_SYNC_FORBIDDEN' => HistorySyncForbiddenException::class,
            'MUTATION_CONFLICT' => MutationDuringPaginationException::class,
            'UNPROCESSABLE_CONTENT' => UnprocessableContentException::class,
            'INVALID_DATE_RANGE' => InvalidDateRangeException::class,
            'MISSING_DATE_RANGE' => MissingDateRangeException::class,
            'INVALID_CURSOR' => InvalidCursorException::class,
            'INVALID_COUNT' => InvalidCountException::class,
            'INVALID_LIMIT' => InvalidLimitException::class,
            'INVALID_QUERY' => InvalidQueryException::class,
            'RATE_LIMIT_EXCEEDED' => RateLimitExceededException::class,
            'SYNC_RATE_LIMIT_EXCEEDED' => SyncRateLimitExceededException::class,
            'SYNC_IN_PROGRESS' => SyncInProgressException::class,
            'BANK_UNDER_MAINTENANCE' => BankUnderMaintenanceException::class,
            'BANK_CONNECTION_NOT_FOUND' => BankConnectionNotFoundException::class,
            'ACCOUNT_NOT_FOUND' => AccountNotFoundException::class,
            'TRANSACTION_NOT_FOUND' => TransactionNotFoundException::class,
            'SYNC_SESSION_NOT_FOUND' => SyncSessionNotFoundException::class,
            'PAYMENT_METHOD_NOT_FOUND' => PaymentMethodNotFoundException::class,
            'TRANSACTION_ORDER_NOT_FOUND' => TransactionOrderNotFoundException::class,
            'BATCH_NOT_FOUND' => BatchNotFoundException::class,
            'VALIDATION_ERROR' => ValidationException::class,
            'INVALID_ORDER_STATE' => InvalidOrderStateException::class,
            'BANK_SUBMISSION_ERROR' => BankSubmissionException::class,
            'BATCH_VALIDATION_ERROR' => BatchValidationException::class,
            'INTERNAL_ERROR' => InternalErrorException::class,
            default => self::class,
        };
    }
}
