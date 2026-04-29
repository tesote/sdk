<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 422 MISSING_DATE_RANGE — endpoint required start_date and/or end_date, none supplied. */
final class MissingDateRangeException extends ApiException
{
}
