import { describe, expect, it } from 'vitest';
import {
  AccountNotFoundError,
  BatchNotFoundError,
  BatchValidationError,
  InvalidOrderStateError,
} from '../src/errors.js';
import { V2Client } from '../src/index.js';
import { callAt, getBody, getHeader, getMethod, jsonResponse, makeFetchMock } from './helpers.js';

const orderStub = (id: string) => ({
  id,
  status: 'draft',
  amount: 100,
  currency: 'VES',
  description: 'pay',
  reference: null,
  external_reference: null,
  idempotency_key: null,
  batch_id: 'b1',
  scheduled_for: null,
  approved_at: null,
  submitted_at: null,
  completed_at: null,
  failed_at: null,
  cancelled_at: null,
  source_account: { id: 'a1', name: 'A1', payment_method_id: 'pm1' },
  destination: { payment_method_id: 'pm2', counterparty_id: 'c1', counterparty_name: 'X' },
  fee: null,
  execution_strategy: null,
  tesote_transaction: null,
  latest_attempt: null,
  created_at: 't',
  updated_at: 't',
});

describe('V2 batches', () => {
  it('create POSTs orders array', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(201, { batch_id: 'b1', orders: [orderStub('o1')], errors: [] }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.batches.create(
      'a1',
      {
        orders: [{ amount: '1', currency: 'VES', description: 'a', beneficiary: { name: 'X' } }],
      },
      { idempotencyKey: 'B1' },
    );
    expect(getMethod(callAt(calls, 0))).toBe('POST');
    expect(calls[0]?.url).toContain('/v2/accounts/a1/batches');
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toBe('B1');
    const body = getBody(callAt(calls, 0)) as { orders: unknown[] };
    expect(body.orders).toHaveLength(1);
  });

  it('create → 400 BATCH_VALIDATION_ERROR', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(400, { error_code: 'BATCH_VALIDATION_ERROR', error: 'bad batch' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.batches.create('a1', { orders: [] })).rejects.toBeInstanceOf(
      BatchValidationError,
    );
  });

  it('get returns summary', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, {
        batch_id: 'b1',
        total_orders: 1,
        total_amount_cents: 100,
        amount_currency: 'VES',
        statuses: { draft: 1 },
        batch_status: 'draft',
        created_at: 't',
        orders: [orderStub('o1')],
      }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.batches.get('a1', 'b1');
    expect(r.batch_status).toBe('draft');
    expect(calls[0]?.url).toContain('/v2/accounts/a1/batches/b1');
  });

  it('get → 404 BATCH_NOT_FOUND', async () => {
    const { fetch } = makeFetchMock([jsonResponse(404, { error_code: 'BATCH_NOT_FOUND' })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.batches.get('a1', 'b1')).rejects.toBeInstanceOf(BatchNotFoundError);
  });

  it('approve POSTs and returns counts', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, { approved: 5, failed: 0 })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.batches.approve('a1', 'b1');
    expect(r.approved).toBe(5);
    expect(calls[0]?.url).toContain('/v2/accounts/a1/batches/b1/approve');
    expect(getMethod(callAt(calls, 0))).toBe('POST');
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toMatch(/^[0-9a-f-]{36}$/);
  });

  it('submit posts token', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, { enqueued: 3, failed: 0 })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.batches.submit('a1', 'b1', { token: 'OTP' });
    expect(getBody(callAt(calls, 0))).toEqual({ token: 'OTP' });
    expect(calls[0]?.url).toContain('/v2/accounts/a1/batches/b1/submit');
  });

  it('cancel returns shape', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { cancelled: 2, skipped: 1, errors: [] }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const r = await c.batches.cancel('a1', 'b1');
    expect(r.cancelled).toBe(2);
    expect(r.skipped).toBe(1);
    expect(calls[0]?.url).toContain('/v2/accounts/a1/batches/b1/cancel');
  });

  it('approve → 409 INVALID_ORDER_STATE', async () => {
    const { fetch } = makeFetchMock([jsonResponse(409, { error_code: 'INVALID_ORDER_STATE' })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.batches.approve('a1', 'b1')).rejects.toBeInstanceOf(InvalidOrderStateError);
  });

  it('create → 404 ACCOUNT_NOT_FOUND', async () => {
    const { fetch } = makeFetchMock([jsonResponse(404, { error_code: 'ACCOUNT_NOT_FOUND' })]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.batches.create('zz', { orders: [] })).rejects.toBeInstanceOf(
      AccountNotFoundError,
    );
  });
});
