<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /v2/accounts/{id}/sync_sessions. */
final readonly class SyncSessionList
{
    /**
     * @param list<SyncSession> $syncSessions
     */
    public function __construct(
        public array $syncSessions,
        public int $limit,
        public int $offset,
        public bool $hasMore,
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
