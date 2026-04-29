<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 422 INVALID_CURSOR — supplied pagination cursor is malformed or expired. */
final class InvalidCursorException extends ApiException
{
}
