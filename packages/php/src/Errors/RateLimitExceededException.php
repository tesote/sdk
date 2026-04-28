<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 429 RATE_LIMIT_EXCEEDED — retries exhausted; honour retryAfter before retrying. */
final class RateLimitExceededException extends ApiException
{
}
