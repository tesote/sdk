<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 500 INTERNAL_ERROR — unhandled server failure; errorId points at server-side trace. */
final class InternalErrorException extends ApiException
{
}
