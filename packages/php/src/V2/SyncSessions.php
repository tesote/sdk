<?php

declare(strict_types=1);

namespace Tesote\Sdk\V2;

use Tesote\Sdk\Models\SyncSession;
use Tesote\Sdk\Models\SyncSessionList;
use Tesote\Sdk\Transport;

/** GET /v2/accounts/{id}/sync_sessions[/{session_id}]. */
final class SyncSessions
{
    public function __construct(private readonly Transport $transport)
    {
    }

    /**
     * @param array{
     *     limit?: int,
     *     offset?: int,
     *     status?: string,
     * } $query
     */
    public function listForAccount(string $accountId, array $query = []): SyncSessionList
    {
        $body = $this->transport->request(
            'GET',
            '/v2/accounts/' . rawurlencode($accountId) . '/sync_sessions',
            $query,
        ) ?? [];
        return SyncSessionList::fromArray($body);
    }

    public function get(string $accountId, string $sessionId): SyncSession
    {
        $body = $this->transport->request(
            'GET',
            '/v2/accounts/' . rawurlencode($accountId) . '/sync_sessions/' . rawurlencode($sessionId),
        ) ?? [];
        return SyncSession::fromArray($body);
    }
}
