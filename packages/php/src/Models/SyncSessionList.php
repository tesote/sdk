<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /v2/accounts/{id}/sync_sessions. */
final class SyncSessionList
{
    /**
     * @param list<SyncSession> $syncSessions
     */
    public function __construct(
        public readonly array $syncSessions,
        public readonly int $limit,
        public readonly int $offset,
        public readonly bool $hasMore,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $sessions = [];
        foreach ((is_array($data['sync_sessions'] ?? null) ? $data['sync_sessions'] : []) as $entry) {
            if (is_array($entry)) {
                $sessions[] = SyncSession::fromArray($entry);
            }
        }

        return new self(
            syncSessions: $sessions,
            limit: (int) ($data['limit'] ?? 0),
            offset: (int) ($data['offset'] ?? 0),
            hasMore: (bool) ($data['has_more'] ?? false),
        );
    }
}
