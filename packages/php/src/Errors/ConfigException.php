<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

use RuntimeException;

/**
 * Bad SDK config detected at construction. Raised before any HTTP call.
 *
 * Intentionally not a TesoteException: it predates the request lifecycle,
 * so the request-summary / attempts fields aren't meaningful.
 */
final class ConfigException extends RuntimeException
{
}
