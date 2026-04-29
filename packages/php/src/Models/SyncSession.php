<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** SyncSession (one row from /v2/accounts/{id}/sync_sessions). */
final readonly class SyncSession
{
    /**
     * @param array{type: string, message: string}|null                                          $error
     * @param array{total_duration: float, complexity_score: float, sync_speed_score: float}|null $performance
     */
    public function __construct(
        public string $id,
        public string $status,
        public string $startedAt,
        public ?string $completedAt,
        public int $transactionsSynced,
        public int $accountsCount,
        public ?array $error,
        public ?array $performance,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $error = null;
        if (is_array($data['error'] ?? null)) {
            $error = [
                'type' => (string) ($data['error']['type'] ?? ''),
                'message' => (string) ($data['error']['message'] ?? ''),
            ];
        }
        $performance = null;
        if (is_array($data['performance'] ?? null)) {
            $performance = [
                'total_duration' => (float) ($data['performance']['total_duration'] ?? 0.0),
                'complexity_score' => (float) ($data['performance']['complexity_score'] ?? 0.0),
                'sync_speed_score' => (float) ($data['performance']['sync_speed_score'] ?? 0.0),
            ];
        }

        return new self(
            id: (string) ($data['id'] ?? ''),
            status: (string) ($data['status'] ?? ''),
            startedAt: (string) ($data['started_at'] ?? ''),
            completedAt: isset($data['completed_at']) ? (string) $data['completed_at'] : null,
            transactionsSynced: (int) ($data['transactions_synced'] ?? 0),
            accountsCount: (int) ($data['accounts_count'] ?? 0),
            error: $error,
            performance: $performance,
        );
    }
}
