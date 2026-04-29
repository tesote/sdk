<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Cursor pagination block (used by transactions index endpoints). */
final readonly class CursorPagination
{
    public function __construct(
        public bool $hasMore,
        public int $perPage,
        public ?string $afterId,
        public ?string $beforeId,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            hasMore: (bool) ($data['has_more'] ?? false),
            perPage: (int) ($data['per_page'] ?? 0),
            afterId: isset($data['after_id']) ? (string) $data['after_id'] : null,
            beforeId: isset($data['before_id']) ? (string) $data['before_id'] : null,
        );
    }
}
