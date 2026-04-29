<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Result envelope returned by POST /v2/.../transactions/sync. */
final readonly class SyncResult
{
    /**
     * @param list<SyncTransaction> $added
     * @param list<SyncTransaction> $modified
     * @param list<SyncRemoval>     $removed
     */
    public function __construct(
        public array $added,
        public array $modified,
        public array $removed,
        public ?string $nextCursor,
        public bool $hasMore,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $added = [];
        foreach ((is_array($data['added'] ?? null) ? $data['added'] : []) as $entry) {
            if (is_array($entry)) {
                $added[] = SyncTransaction::fromArray($entry);
            }
        }
        $modified = [];
        foreach ((is_array($data['modified'] ?? null) ? $data['modified'] : []) as $entry) {
            if (is_array($entry)) {
                $modified[] = SyncTransaction::fromArray($entry);
            }
        }
        $removed = [];
        foreach ((is_array($data['removed'] ?? null) ? $data['removed'] : []) as $entry) {
            if (is_array($entry)) {
                $removed[] = SyncRemoval::fromArray($entry);
            }
        }

        return new self(
            added: $added,
            modified: $modified,
            removed: $removed,
            nextCursor: isset($data['next_cursor']) ? (string) $data['next_cursor'] : null,
            hasMore: (bool) ($data['has_more'] ?? false),
        );
    }
}
