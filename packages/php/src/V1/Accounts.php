<?php

declare(strict_types=1);

namespace Tesote\Sdk\V1;

use Tesote\Sdk\Transport;

/**
 * v1 accounts: read-only listing and lookup.
 *
 * Other v2/v3 mutations (sync) are not part of v1.
 */
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
        return $this->transport->request('GET', '/v1/accounts', $query);
    }

    /**
     * @return array<mixed>|null
     */
    public function get(string $id): ?array
    {
        return $this->transport->request('GET', '/v1/accounts/' . rawurlencode($id));
    }
}
