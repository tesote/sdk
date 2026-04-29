<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 429 SYNC_RATE_LIMIT_EXCEEDED — bank-connection-level sync throttle hit; honour retryAfter. */
final class SyncRateLimitExceededException extends ApiException
{
}
