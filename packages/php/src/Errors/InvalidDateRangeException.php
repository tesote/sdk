<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 422 INVALID_DATE_RANGE — from > to, or range exceeds the API's window. */
final class InvalidDateRangeException extends ApiException
{
}
