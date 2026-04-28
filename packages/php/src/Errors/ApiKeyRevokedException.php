<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 401 API_KEY_REVOKED — the key existed but was revoked; rotate it. */
final class ApiKeyRevokedException extends ApiException
{
}
