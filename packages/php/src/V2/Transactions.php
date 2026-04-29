<?php

declare(strict_types=1);

namespace Tesote\Sdk\V2;

use Tesote\Sdk\Models\BulkResult;
use Tesote\Sdk\Models\SyncResult;
use Tesote\Sdk\Models\Transaction;
use Tesote\Sdk\Models\TransactionSearchResult;
use Tesote\Sdk\Transport;

/** v2 transactions: lookup, legacy sync, bulk, search. */
final class Transactions
{
    private const SHOW_TTL = 300;

    public function __construct(private readonly Transport $transport)
    {
    }

    public function get(string $id): Transaction
    {
        $body = $this->transport->request(
            'GET',
            '/v2/transactions/' . rawurlencode($id),
            null,
            null,
            ['cacheTtl' => self::SHOW_TTL],
        ) ?? [];
        return Transaction::fromArray($body);
    }

    /**
     * Legacy non-nested sync — POST /v2/transactions/sync. Prefer V2\Accounts::syncTransactions.
     *
     * @param array<string, mixed> $body
     */
    public function sync(array $body = [], ?string $idempotencyKey = null): SyncResult
    {
        $payload = $body !== [] ? $body : (object) [];
        $decoded = $this->transport->request(
            'POST',
            '/v2/transactions/sync',
            null,
            (array) $payload,
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return SyncResult::fromArray($decoded);
    }

    /**
     * @param array{
     *     account_ids: list<string>,
     *     page?: int,
     *     per_page?: int,
     *     limit?: int,
     *     offset?: int,
     * } $body
     */
    public function bulk(array $body, ?string $idempotencyKey = null): BulkResult
    {
        $decoded = $this->transport->request(
            'POST',
            '/v2/transactions/bulk',
            null,
            $body,
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return BulkResult::fromArray($decoded);
    }

    /**
     * @param array{
     *     q: string,
     *     account_id?: string,
     *     limit?: int,
     *     offset?: int,
     *     start_date?: string,
     *     end_date?: string,
     *     status?: string,
     *     category_id?: string,
     *     counterparty_id?: string,
     *     type?: string,
     *     reference_code?: string,
     *     amount_min?: float|int|string,
     *     amount_max?: float|int|string,
     *     amount?: float|int|string,
     *     transaction_date_after?: string,
     *     transaction_date_before?: string,
     *     created_after?: string,
     *     updated_after?: string,
     * } $query
     */
    public function search(array $query): TransactionSearchResult
    {
        $body = $this->transport->request('GET', '/v2/transactions/search', $query) ?? [];
        return TransactionSearchResult::fromArray($body);
    }
}
