<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** Connect or read timeout. Distinct from NetworkException so callers can react differently. */
final class TimeoutException extends TransportException
{
}
