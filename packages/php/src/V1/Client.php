<?php

declare(strict_types=1);

namespace Tesote\Sdk\V1;

use Tesote\Sdk\Transport;

/**
 * v1 client. Read-only foundation: status, whoami, accounts, transactions.
 *
 * Constructor accepts the same shared config shape as V2 — see Transport.
 */
final class Client
{
    public readonly Transport $transport;
    public readonly Accounts $accounts;
    public readonly Transactions $transactions;
    public readonly Status $status;

    /**
     * @param array<string, mixed> $config See Transport::__construct.
     */
    public function __construct(array $config)
    {
        $this->transport = $config['transport'] ?? new Transport($config);
        $this->accounts = new Accounts($this->transport);
        $this->transactions = new Transactions($this->transport);
        $this->status = new Status($this->transport);
    }
}
