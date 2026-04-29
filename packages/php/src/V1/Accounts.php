<?php

declare(strict_types=1);

namespace Tesote\Sdk\V1;

use Tesote\Sdk\Models\Account;
use Tesote\Sdk\Models\AccountList;
use Tesote\Sdk\Models\TransactionList;
use Tesote\Sdk\Transport;

/** v1 accounts: list, get, list-transactions. Read-only. */
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
        $body = $this->transport->request('GET', '/v1/accounts', $query, null, ['cacheTtl' => self::LIST_TTL]) ?? [];
        return AccountList::fromArray($body);
    }

    public function get(string $id): Account
    {
        $body = $this->transport->request(
            'GET',
            '/v1/accounts/' . rawurlencode($id),
            null,
            null,
            ['cacheTtl' => self::SHOW_TTL],
        ) ?? [];
        return Account::fromArray($body);
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
     * } $query
     */
    public function listTransactions(string $accountId, array $query = []): TransactionList
    {
        $body = $this->transport->request(
            'GET',
            '/v1/accounts/' . rawurlencode($accountId) . '/transactions',
            $query,
            null,
            ['cacheTtl' => self::TXN_TTL],
        ) ?? [];
        return TransactionList::fromArray($body);
    }
}
