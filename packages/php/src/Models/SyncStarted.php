<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response envelope returned by POST /v2/accounts/{id}/sync (202 Accepted). */
final readonly class SyncStarted
{
    public function __construct(
        public string $message,
        public string $syncSessionId,
        public string $status,
        public string $startedAt,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            message: (string) ($data['message'] ?? ''),
            syncSessionId: (string) ($data['sync_session_id'] ?? ''),
            status: (string) ($data['status'] ?? ''),
            startedAt: (string) ($data['started_at'] ?? ''),
        );
    }
}
