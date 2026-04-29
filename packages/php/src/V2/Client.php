<?php

declare(strict_types=1);

namespace Tesote\Sdk\V2;

use Tesote\Sdk\Transport;

/**
 * v2 client. v1 surface plus writes for sync, transaction orders, batches,
 * payment methods, plus bulk + search.
 */
final class Client
{
    public readonly Transport $transport;
    public readonly Accounts $accounts;
    public readonly Transactions $transactions;
    public readonly SyncSessions $syncSessions;
    public readonly TransactionOrders $transactionOrders;
    public readonly Batches $batches;
    public readonly PaymentMethods $paymentMethods;
    public readonly Status $status;

    /**
     * @param array<string, mixed> $config See Transport::__construct.
     */
    public function __construct(array $config)
    {
        $this->transport = $config['transport'] ?? new Transport($config);
        $this->accounts = new Accounts($this->transport);
        $this->transactions = new Transactions($this->transport);
        $this->syncSessions = new SyncSessions($this->transport);
        $this->transactionOrders = new TransactionOrders($this->transport);
        $this->batches = new Batches($this->transport);
        $this->paymentMethods = new PaymentMethods($this->transport);
        $this->status = new Status($this->transport);
    }
}
