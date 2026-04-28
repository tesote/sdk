<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 401 UNAUTHORIZED — bearer missing, malformed, or unrecognised. */
final class UnauthorizedException extends ApiException
{
}
