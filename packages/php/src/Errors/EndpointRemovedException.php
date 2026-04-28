<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

use RuntimeException;

/**
 * Thrown when a method's upstream endpoint has been removed in the API
 * version the caller targeted. Carries a hint at the replacement.
 */
final class EndpointRemovedException extends RuntimeException
{
    public function __construct(string $method, public readonly ?string $replacement = null)
    {
        $msg = sprintf('Endpoint %s has been removed.', $method);
        if ($replacement !== null) {
            $msg .= ' Use ' . $replacement . ' instead.';
        }
        parent::__construct($msg);
    }
}
