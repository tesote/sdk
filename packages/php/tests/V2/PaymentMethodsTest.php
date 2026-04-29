<?php

declare(strict_types=1);

namespace Tesote\Sdk\Tests\V2;

use Tesote\Sdk\Errors\PaymentMethodNotFoundException;
use Tesote\Sdk\Errors\ValidationException;
use Tesote\Sdk\Http\CurlResult;
use Tesote\Sdk\Models\PaymentMethod;
use Tesote\Sdk\Models\PaymentMethodList;
use Tesote\Sdk\Tests\Support\TestCaseBase;
use Tesote\Sdk\V2\Client;

final class PaymentMethodsTest extends TestCaseBase
{
    /**
     * @return array<string, mixed>
     */
    private static function fixture(): array
    {
        return [
            'id' => 'pm-1',
            'method_type' => 'bank_account',
            'currency' => 'VES',
            'label' => 'Main',
            'details' => ['bank_code' => '0102', 'account_number' => '****1234', 'holder_name' => 'ACME'],
            'verified' => true,
            'verified_at' => '2026-04-28T00:00:00Z',
            'last_used_at' => null,
            'counterparty' => ['id' => 'cp-1', 'name' => 'Vendor'],
            'tesote_account' => null,
            'created_at' => '2026-01-01T00:00:00Z',
            'updated_at' => '2026-04-01T00:00:00Z',
        ];
    }

    public function testList(): void
    {
        $this->enqueueOk([
            'items' => [self::fixture()],
            'has_more' => false,
            'limit' => 50,
            'offset' => 0,
        ]);
        $client = new Client(['transport' => $this->makeTransport()]);
        $list = $client->paymentMethods->list(['method_type' => 'bank_account', 'verified' => 'true']);

        self::assertInstanceOf(PaymentMethodList::class, $list);
        self::assertSame('pm-1', $list->items[0]->id);
        self::assertStringContainsString('method_type=bank_account', $this->lastUrl());
        self::assertStringContainsString('verified=true', $this->lastUrl());
    }

    public function testGet(): void
    {
        $this->enqueueOk(self::fixture());
        $client = new Client(['transport' => $this->makeTransport()]);
        $pm = $client->paymentMethods->get('pm-1');

        self::assertInstanceOf(PaymentMethod::class, $pm);
        self::assertTrue($pm->verified);
        self::assertSame('Vendor', $pm->counterparty?->name);
    }

    public function testGetMapsNotFound(): void
    {
        $this->enqueueError(404, ['error' => 'gone', 'error_code' => 'PAYMENT_METHOD_NOT_FOUND']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(PaymentMethodNotFoundException::class);
        $client->paymentMethods->get('missing');
    }

    public function testCreateWrapsBody(): void
    {
        $this->enqueueOk(self::fixture(), 201);
        $client = new Client(['transport' => $this->makeTransport()]);
        $pm = $client->paymentMethods->create([
            'method_type' => 'bank_account',
            'currency' => 'VES',
            'label' => 'Main',
            'details' => ['bank_code' => '0102', 'account_number' => '...', 'holder_name' => 'ACME'],
        ], 'idem-pm-1');

        self::assertInstanceOf(PaymentMethod::class, $pm);
        $body = json_decode($this->lastBody() ?? '', true);
        self::assertSame('bank_account', $body['payment_method']['method_type']);
        self::assertSame('idem-pm-1', $this->lastHeaders()['Idempotency-Key']);
    }

    public function testUpdatePatchesPartial(): void
    {
        $this->enqueueOk(self::fixture());
        $client = new Client(['transport' => $this->makeTransport()]);
        $client->paymentMethods->update('pm-1', ['label' => 'Renamed']);

        self::assertSame('PATCH', $this->lastMethod());
        $body = json_decode($this->lastBody() ?? '', true);
        self::assertSame('Renamed', $body['payment_method']['label']);
    }

    public function testUpdateMapsValidation(): void
    {
        $this->enqueueError(400, ['error' => 'bad', 'error_code' => 'VALIDATION_ERROR']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(ValidationException::class);
        $client->paymentMethods->update('pm-1', ['label' => '']);
    }

    public function testDeleteReturnsVoidOn204(): void
    {
        $this->curl->enqueue(new CurlResult(204, '', ['x-request-id' => 'req-x'], 0, ''));
        $client = new Client(['transport' => $this->makeTransport()]);
        $client->paymentMethods->delete('pm-1', 'idem-del-1');

        self::assertSame('DELETE', $this->lastMethod());
        self::assertSame('idem-del-1', $this->lastHeaders()['Idempotency-Key']);
    }

    public function testDeleteMapsValidationWhenInUse(): void
    {
        $this->enqueueError(409, ['error' => 'in use', 'error_code' => 'VALIDATION_ERROR']);
        $client = new Client(['transport' => $this->makeTransport()]);
        $this->expectException(ValidationException::class);
        $client->paymentMethods->delete('pm-1');
    }
}
