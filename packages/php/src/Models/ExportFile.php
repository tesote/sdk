<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/**
 * Response from GET /v2/accounts/{id}/transactions/export.
 *
 * `body` is the raw file payload (CSV bytes or pretty-printed JSON string).
 * `format` echoes back the requested format. `filename` reflects the server's
 * Content-Disposition suggestion (or null if the SDK couldn't parse it).
 */
final class ExportFile
{
    public function __construct(
        public readonly string $body,
        public readonly string $format,
        public readonly ?string $filename,
    ) {
    }
}
