<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 422 UNPROCESSABLE_CONTENT — server validation failure; inspect responseBody for details. */
final class UnprocessableContentException extends ApiException
{
}
