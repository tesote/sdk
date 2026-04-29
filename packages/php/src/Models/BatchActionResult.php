<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/**
 * Generic result envelope returned by POST /batches/{id}/{approve,submit,cancel}.
 *
 * Each endpoint returns a slightly different mix of counters, kept on a single
 * value object so callers can branch on whichever fields are populated.
 */
final class BatchActionResult
{
    /**
     * @param list<array<string, mixed>> $errors
     */
    public function __construct(
        public readonly ?int $approved,
        public readonly ?int $enqueued,
        public readonly ?int $cancelled,
        public readonly ?int $skipped,
        public readonly ?int $failed,
        public readonly array $errors,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $errors = [];
        foreach ((is_array($data['errors'] ?? null) ? $data['errors'] : []) as $entry) {
            if (is_array($entry)) {
                /** @var array<string, mixed> $entry */
                $errors[] = $entry;
            }
        }

        return new self(
            approved: isset($data['approved']) ? (int) $data['approved'] : null,
            enqueued: isset($data['enqueued']) ? (int) $data['enqueued'] : null,
            cancelled: isset($data['cancelled']) ? (int) $data['cancelled'] : null,
            skipped: isset($data['skipped']) ? (int) $data['skipped'] : null,
            failed: isset($data['failed']) ? (int) $data['failed'] : null,
            errors: $errors,
        );
    }
}
