<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from POST /v2/transactions/bulk. */
final readonly class BulkResult
{
    /**
     * @param list<BulkAccountResult> $bulkResults
     */
    public function __construct(
        public array $bulkResults,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $rows = [];
        foreach ((is_array($data['bulk_results'] ?? null) ? $data['bulk_results'] : []) as $entry) {
            if (is_array($entry)) {
                $rows[] = BulkAccountResult::fromArray($entry);
            }
        }

        return new self(bulkResults: $rows);
    }
}
