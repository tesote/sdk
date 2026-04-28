<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/**
 * Anything that prevented us from getting a usable HTTP response.
 *
 * Subclasses split DNS/connection (NetworkException), timeouts
 * (TimeoutException) and TLS handshake failures (TlsException).
 */
class TransportException extends TesoteException
{
}
