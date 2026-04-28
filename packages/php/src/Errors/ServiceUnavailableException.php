<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 503 — service unavailable (pause mode, planned downtime). Backoff and retry. */
final class ServiceUnavailableException extends ApiException
{
}
