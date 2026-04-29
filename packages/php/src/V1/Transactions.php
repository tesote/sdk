<?php

declare(strict_types=1);

namespace Tesote\Sdk\V1;

use Tesote\Sdk\Models\Transaction;
use Tesote\Sdk\Transport;

/** v1 transactions: lookup by id. */
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
            '/v1/transactions/' . rawurlencode($id),
            null,
            null,
            ['cacheTtl' => self::SHOW_TTL],
        ) ?? [];
        return Transaction::fromArray($body);
    }
}
