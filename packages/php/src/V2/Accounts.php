<?php

declare(strict_types=1);

namespace Tesote\Sdk\V2;

use Tesote\Sdk\Models\Account;
use Tesote\Sdk\Models\AccountList;
use Tesote\Sdk\Models\ExportFile;
use Tesote\Sdk\Models\SyncResult;
use Tesote\Sdk\Models\SyncStarted;
use Tesote\Sdk\Models\TransactionList;
use Tesote\Sdk\Transport;

/** v2 accounts: list, get, sync, transactions index, transaction sync, export. */
final class Accounts
{
    private const LIST_TTL = 60;
    private const SHOW_TTL = 300;
    private const TXN_TTL = 60;

    public function __construct(private readonly Transport $transport)
    {
    }

    /**
     * @param array{
     *     page?: int,
     *     per_page?: int,
     *     include?: string,
     *     sort?: string,
     * } $query
     */
    public function list(array $query = []): AccountList
    {
        $body = $this->transport->request('GET', '/v2/accounts', $query, null, ['cacheTtl' => self::LIST_TTL]) ?? [];
        return AccountList::fromArray($body);
    }

    public function get(string $id): Account
    {
        $body = $this->transport->request(
            'GET',
            '/v2/accounts/' . rawurlencode($id),
            null,
            null,
            ['cacheTtl' => self::SHOW_TTL],
        ) ?? [];
        return Account::fromArray($body);
    }

    public function sync(string $id, ?string $idempotencyKey = null): SyncStarted
    {
        $body = $this->transport->request(
            'POST',
            '/v2/accounts/' . rawurlencode($id) . '/sync',
            null,
            [],
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return SyncStarted::fromArray($body);
    }

    /**
     * @param array{
     *     start_date?: string,
     *     end_date?: string,
     *     scope?: string,
     *     page?: int,
     *     per_page?: int,
     *     transactions_after_id?: string,
     *     transactions_before_id?: string,
     *     transaction_date_after?: string,
     *     transaction_date_before?: string,
     *     created_after?: string,
     *     updated_after?: string,
     *     amount_min?: float|int|string,
     *     amount_max?: float|int|string,
     *     amount?: float|int|string,
     *     status?: string,
     *     category_id?: string,
     *     counterparty_id?: string,
     *     q?: string,
     *     type?: string,
     *     reference_code?: string,
     * } $query
     */
    public function listTransactions(string $accountId, array $query = []): TransactionList
    {
        $body = $this->transport->request(
            'GET',
            '/v2/accounts/' . rawurlencode($accountId) . '/transactions',
            $query,
            null,
            ['cacheTtl' => self::TXN_TTL],
        ) ?? [];
        return TransactionList::fromArray($body);
    }

    /**
     * Export the transactions for an account as CSV or JSON.
     *
     * @param array{
     *     format?: 'csv'|'json',
     *     start_date?: string,
     *     end_date?: string,
     *     scope?: string,
     *     status?: string,
     *     category_id?: string,
     *     counterparty_id?: string,
     *     q?: string,
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
    public function exportTransactions(string $accountId, array $query = []): ExportFile
    {
        $format = isset($query['format']) ? (string) $query['format'] : 'csv';
        $raw = $this->transport->requestRaw(
            'GET',
            '/v2/accounts/' . rawurlencode($accountId) . '/transactions/export',
            $query,
        );
        return new ExportFile(
            body: $raw['body'],
            format: $format,
            filename: self::parseFilename($raw['headers']['content-disposition'] ?? null),
        );
    }

    /**
     * @param array{
     *     count?: int,
     *     cursor?: string|null,
     *     options?: array{include_running_balance?: bool},
     * } $body
     */
    public function syncTransactions(string $accountId, array $body = [], ?string $idempotencyKey = null): SyncResult
    {
        $payload = $body !== [] ? $body : (object) [];
        $decoded = $this->transport->request(
            'POST',
            '/v2/accounts/' . rawurlencode($accountId) . '/transactions/sync',
            null,
            (array) $payload,
            ['idempotencyKey' => $idempotencyKey],
        ) ?? [];
        return SyncResult::fromArray($decoded);
    }

    private static function parseFilename(?string $disposition): ?string
    {
        if ($disposition === null) {
            return null;
        }
        // why: matches both filename="name.csv" and filename=name.csv (RFC 6266 minimum).
        if (preg_match('/filename\*?=(?:"([^"]+)"|([^;\s]+))/i', $disposition, $m) === 1) {
            return $m[1] !== '' ? $m[1] : ($m[2] ?? null);
        }
        return null;
    }
}
