<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** TLS handshake / certificate-validation failure. Almost always a config issue. */
final class TlsException extends TransportException
{
}
