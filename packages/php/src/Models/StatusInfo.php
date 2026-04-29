<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /status and GET /v2/status. */
final class StatusInfo
{
    public function __construct(
        public readonly string $status,
        public readonly bool $authenticated,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            status: (string) ($data['status'] ?? ''),
            authenticated: (bool) ($data['authenticated'] ?? false),
        );
    }
}
