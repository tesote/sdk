<?php

declare(strict_types=1);

namespace Tesote\Sdk\V1;

use Tesote\Sdk\NotImplemented;
use Tesote\Sdk\Transport;

/**
 * v1 client. Read-only foundation: accounts + transactions, plus status/whoami.
 *
 * Constructor accepts the same shared config shape as V2 — see Transport.
 */
final class Client
{
    public readonly Transport $transport;
    public readonly Accounts $accounts;
    public readonly NotImplemented $transactions;
    public readonly NotImplemented $status;

    /**
     * @param array<string, mixed> $config See Transport::__construct.
     */
    public function __construct(array $config)
    {
        $this->transport = $config['transport'] ?? new Transport($config);
        $this->accounts = new Accounts($this->transport);
        $this->transactions = new NotImplemented('transactions');
        $this->status = new NotImplemented('status');
    }
}
