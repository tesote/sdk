import { describe, expect, it } from 'vitest';
import {
  AccountNotFoundError,
  BankSubmissionError,
  InvalidOrderStateError,
  TransactionOrderNotFoundError,
  ValidationError,
} from '../src/errors.js';
import { V2Client } from '../src/index.js';
import { callAt, getBody, getHeader, getMethod, jsonResponse, makeFetchMock } from './helpers.js';

const order = (id: string, status: 'draft' | 'processing' | 'cancelled' = 'draft') => ({
  id,
  status,
  amount: 100,
  currency: 'VES',
  description: 'pay',
  reference: null,
  external_reference: null,
  idempotency_key: null,
  batch_id: null,
  scheduled_for: null,
  approved_at: null,
  submitted_at: null,
  completed_at: null,
  failed_at: null,
  cancelled_at: null,
  source_account: { id: 'a1', name: 'A1', payment_method_id: 'pm1' },
  destination: { payment_method_id: 'pm2', counterparty_id: 'c1', counterparty_name: 'Bob' },
  fee: null,
  execution_strategy: null,
  tesote_transaction: null,
  latest_attempt: null,
  created_at: 't',
  updated_at: 't',
});

describe('V2 transactionOrders', () => {
  it('list serializes filters', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { items: [order('o1')], has_more: false, limit: 50, offset: 0 }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.transactionOrders.list('a1', { status: 'draft', limit: 25 });
    expect(calls[0]?.url).toContain('/v2/accounts/a1/transaction_orders');
    expect(calls[0]?.url).toContain('status=draft');
    expect(calls[0]?.url).toContain('limit=25');
  });

  it('get returns order', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, order('o1'))]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.transactionOrders.get('a1', 'o1');
    expect(r.id).toBe('o1');
    expect(calls[0]?.url).toContain('/v2/accounts/a1/transaction_orders/o1');
  });

  it('get → 404 TRANSACTION_ORDER_NOT_FOUND', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(404, { error_code: 'TRANSACTION_ORDER_NOT_FOUND' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.transactionOrders.get('a1', 'o1')).rejects.toBeInstanceOf(
      TransactionOrderNotFoundError,
    );
  });

  it('create wraps body under transaction_order key', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(201, order('o1'))]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.transactionOrders.create(
      'a1',
      {
        amount: '10.00',
        currency: 'VES',
        description: 'pay',
        beneficiary: { name: 'Bob' },
      },
      { idempotencyKey: 'I1' },
    );
    expect(getMethod(callAt(calls, 0))).toBe('POST');
    expect(getBody(callAt(calls, 0))).toEqual({
      transaction_order: {
        amount: '10.00',
        currency: 'VES',
        description: 'pay',
        beneficiary: { name: 'Bob' },
      },
    });
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toBe('I1');
    expect(getHeader(callAt(calls, 0), 'content-type')).toBe('application/json');
  });

  it('create → 400 VALIDATION_ERROR', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(400, { error_code: 'VALIDATION_ERROR', error: 'amount required' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(
      c.transactionOrders.create('a1', { amount: '', currency: 'VES', description: '' }),
    ).rejects.toBeInstanceOf(ValidationError);
  });

  it('create → 404 ACCOUNT_NOT_FOUND', async () => {
    const { fetch } = makeFetchMock([jsonResponse(404, { error_code: 'ACCOUNT_NOT_FOUND' })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(
      c.transactionOrders.create('zz', { amount: '1', currency: 'VES', description: 'x' }),
    ).rejects.toBeInstanceOf(AccountNotFoundError);
  });

  it('submit posts token body', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(202, order('o1', 'processing'))]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.transactionOrders.submit('a1', 'o1', { token: 'OTP123' });
    expect(r.status).toBe('processing');
    expect(calls[0]?.url).toContain('/v2/accounts/a1/transaction_orders/o1/submit');
    expect(getBody(callAt(calls, 0))).toEqual({ token: 'OTP123' });
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toMatch(/^[0-9a-f-]{36}$/);
  });

  it('submit no token body is empty object', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(202, order('o1', 'processing'))]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.transactionOrders.submit('a1', 'o1');
    expect(getBody(callAt(calls, 0))).toEqual({});
  });

  it('submit → 409 INVALID_ORDER_STATE', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(409, { error_code: 'INVALID_ORDER_STATE', error: 'wrong state' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.transactionOrders.submit('a1', 'o1')).rejects.toBeInstanceOf(
      InvalidOrderStateError,
    );
  });

  it('submit → 422 BANK_SUBMISSION_ERROR', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(422, { error_code: 'BANK_SUBMISSION_ERROR', error: 'rejected' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.transactionOrders.submit('a1', 'o1')).rejects.toBeInstanceOf(
      BankSubmissionError,
    );
  });

  it('cancel POSTs to /cancel', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, order('o1', 'cancelled'))]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.transactionOrders.cancel('a1', 'o1', { idempotencyKey: 'C1' });
    expect(r.status).toBe('cancelled');
    expect(calls[0]?.url).toContain('/v2/accounts/a1/transaction_orders/o1/cancel');
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toBe('C1');
  });

  it('listAll iterates offset pagination', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(200, {
        items: [order('o1'), order('o2')],
        has_more: true,
        limit: 2,
        offset: 0,
      }),
      jsonResponse(200, { items: [order('o3')], has_more: false, limit: 2, offset: 2 }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const ids: string[] = [];
    for await (const o of c.transactionOrders.listAll('a1', { limit: 2 })) ids.push(o.id);
    expect(ids).toEqual(['o1', 'o2', 'o3']);
  });
});
