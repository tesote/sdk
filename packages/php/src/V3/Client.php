<?php

declare(strict_types=1);

namespace Tesote\Sdk\V3;

use Tesote\Sdk\NotImplemented;
use Tesote\Sdk\Transport;

/**
 * v3 client. Full surface: v2 plus categories, counterparties, legal entities,
 * connections, webhooks, reports, balance history, workspace, MCP.
 *
 * Webhook signature verification is a stateless static helper; it does not
 * live on this object.
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
    public readonly NotImplemented $categories;
    public readonly NotImplemented $counterparties;
    public readonly NotImplemented $legalEntities;
    public readonly NotImplemented $connections;
    public readonly NotImplemented $webhooks;
    public readonly NotImplemented $reports;
    public readonly NotImplemented $balanceHistory;
    public readonly NotImplemented $workspace;
    public readonly NotImplemented $mcp;
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
        $this->categories = new NotImplemented('categories');
        $this->counterparties = new NotImplemented('counterparties');
        $this->legalEntities = new NotImplemented('legal_entities');
        $this->connections = new NotImplemented('connections');
        $this->webhooks = new NotImplemented('webhooks');
        $this->reports = new NotImplemented('reports');
        $this->balanceHistory = new NotImplemented('balance_history');
        $this->workspace = new NotImplemented('workspace');
        $this->mcp = new NotImplemented('mcp');
        $this->status = new NotImplemented('status');
    }
}
