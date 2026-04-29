<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V2;

use Tesote\Sdk\Errors\InvalidOrderStateException;
use Tesote\Sdk\Errors\TransactionOrderNotFoundException;
use Tesote\Sdk\Errors\ValidationException;
use Tesote\Sdk\Models\TransactionOrder;
use Tesote\Sdk\Models\TransactionOrderList;
use Tesote\Sdk\Tests\Support\TestCaseBase;
use Tesote\Sdk\V2\Client;

final class TransactionOrdersTest extends TestCaseBase
{
    /**
     * @return array<string, mixed>
     */
    private static function orderFixture(string $status = 'draft'): array
    {
        return [
            'id' => 'ord-1',
            'status' => $status,
            'amount' => '1000.00',
            'currency' => 'VES',
            'description' => 'Rent',
            'reference' => null,
            'external_reference' => null,
            'idempotency_key' => null,
            'batch_id' => null,
            'scheduled_for' => null,
            'approved_at' => null,
            'submitted_at' => null,
            'completed_at' => null,
            'failed_at' => null,
            'cancelled_at' => null,
            'source_account' => ['id' => 'acct-1', 'name' => 'Main', 'payment_method_id' => 'pm-1'],
            'destination' => ['payment_method_id' => 'pm-2', 'counterparty_id' => 'cp-1', 'counterparty_name' => 'Landlord'],
            'fee' => null,
            'execution_strategy' => null,
            'tesote_transaction' => null,
            'latest_attempt' => null,
            'created_at' => '2026-04-28T10:00:00Z',
            'updated_at' => '2026-04-28T10:00:00Z',
        ];
    }

    public function testListForAccount(): void
    {
        $this->enqueueOk([
            'items' => [self::orderFixture()],
            'has_more' => false,
            'limit' => 50,
            'offset' => 0,
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $list = $client->transactionOrders->listForAccount('acct-1', ['status' => 'draft', 'limit' => 50]);

        self::assertInstanceOf(TransactionOrderList::class, $list);
        self::assertCount(1, $list->items);
        self::assertSame('ord-1', $list->items[0]->id);
        self::assertStringContainsString('status=draft', $this->lastUrl());
    }

    public function testGet(): void
    {
        $this->enqueueOk(self::orderFixture('processing'));
        $client = new Client(['transport' => $this->makeTransport()]);
        $order = $client->transactionOrders->get('acct-1', 'ord-1');

        self::assertInstanceOf(TransactionOrder::class, $order);
        self::assertSame('processing', $order->status);
    }

    public function testGetMapsNotFound(): void
    {
        $this->enqueueError(404, ['error' => 'gone', 'error_code' => 'TRANSACTION_ORDER_NOT_FOUND']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(TransactionOrderNotFoundException::class);
        $client->transactionOrders->get('acct-1', 'missing');
    }

    public function testCreateWrapsBodyAndUsesIdempotencyKey(): void
    {
        $this->enqueueOk(self::orderFixture(), 201);
        $client = new Client(['transport' => $this->makeTransport()]);
        $order = $client->transactionOrders->create(
            'acct-1',
            [
                'destination_payment_method_id' => 'pm-2',
                'amount' => '1000.00',
                'currency' => 'VES',
                'description' => 'Rent',
            ],
            'idem-create-1',
        );

        self::assertInstanceOf(TransactionOrder::class, $order);
        self::assertSame('idem-create-1', $this->lastHeaders()['Idempotency-Key']);
        $body = json_decode($this->lastBody() ?? '', true);
        self::assertSame('1000.00', $body['transaction_order']['amount']);
        self::assertSame('Rent', $body['transaction_order']['description']);
    }

    public function testCreateMapsValidationError(): void
    {
        $this->enqueueError(400, ['error' => 'bad amount', 'error_code' => 'VALIDATION_ERROR']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(ValidationException::class);
        $client->transactionOrders->create('acct-1', [
            'amount' => 'NaN', 'currency' => 'VES', 'description' => 'oops',
        ]);
    }

    public function testSubmitMapsInvalidOrderState(): void
    {
        $this->enqueueError(409, ['error' => 'wrong state', 'error_code' => 'INVALID_ORDER_STATE']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(InvalidOrderStateException::class);
        $client->transactionOrders->submit('acct-1', 'ord-1', 'OTP-12345');
    }

    public function testCancelMutates(): void
    {
        $this->enqueueOk(self::orderFixture('cancelled'));
        $client = new Client(['transport' => $this->makeTransport()]);
        $order = $client->transactionOrders->cancel('acct-1', 'ord-1');

        self::assertSame('cancelled', $order->status);
        self::assertSame('POST', $this->lastMethod());
        self::assertStringEndsWith('/cancel', $this->lastUrl());
    }
}
