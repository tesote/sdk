<?php

declare(strict_types=1);

namespace Tesote\Sdk\V2;

use Tesote\Sdk\NotImplemented;
use Tesote\Sdk\Transport;

/**
 * v2 client. v1 surface plus writes for sync, transaction orders, batches,
 * payment methods, plus bulk + search.
 */
final class Client
{
    public readonly Transport $transport;
    public readonly Accounts $accounts;
    public readonly NotImplemented $transactions;
    public readonly NotImplemented $syncSessions;
    public readonly NotImplemented $transactionOrders;
    public readonly NotImplemented $batches;
    public readonly NotImplemented $paymentMethods;
    public readonly NotImplemented $status;

    /**
     * @param array<string, mixed> $config See Transport::__construct.
     */
    public function __construct(array $config)
    {
        $this->transport = $config['transport'] ?? new Transport($config);
        $this->accounts = new Accounts($this->transport);
        $this->transactions = new NotImplemented('transactions');
        $this->syncSessions = new NotImplemented('sync_sessions');
        $this->transactionOrders = new NotImplemented('transaction_orders');
        $this->batches = new NotImplemented('batches');
        $this->paymentMethods = new NotImplemented('payment_methods');
        $this->status = new NotImplemented('status');
    }
}
