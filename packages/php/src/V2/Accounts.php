<?php

declare(strict_types=1);

namespace Tesote\Sdk\V2;

use Tesote\Sdk\Transport;

/** v2 accounts: list, get, sync. */
final class Accounts
{
    public function __construct(private readonly Transport $transport)
    {
    }

    /**
     * @param  array<string, scalar|array<int|string, scalar>>|null $query
     * @return array<mixed>|null
     */
    public function list(?array $query = null): ?array
    {
        return $this->transport->request('GET', '/v2/accounts', $query);
    }

    /**
     * @return array<mixed>|null
     */
    public function get(string $id): ?array
    {
        return $this->transport->request('GET', '/v2/accounts/' . rawurlencode($id));
    }

    /**
     * @return array<mixed>|null
     */
    public function sync(string $id, ?string $idempotencyKey = null): ?array
    {
        return $this->transport->request(
            'POST',
            '/v2/accounts/' . rawurlencode($id) . '/sync',
            null,
            [],
            ['idempotencyKey' => $idempotencyKey],
        );
    }
}
