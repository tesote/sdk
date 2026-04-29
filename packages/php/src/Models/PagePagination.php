<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Page-based pagination block (used by accounts list endpoints). */
final readonly class PagePagination
{
    public function __construct(
        public int $currentPage,
        public int $perPage,
        public int $totalPages,
        public int $totalCount,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            currentPage: (int) ($data['current_page'] ?? 0),
            perPage: (int) ($data['per_page'] ?? 0),
            totalPages: (int) ($data['total_pages'] ?? 0),
            totalCount: (int) ($data['total_count'] ?? 0),
        );
    }
}
