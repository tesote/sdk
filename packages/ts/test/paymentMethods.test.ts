import { describe, expect, it } from 'vitest';
import { PaymentMethodNotFoundError, ValidationError } from '../src/errors.js';
import { V2Client } from '../src/index.js';
import {
  callAt,
  getBody,
  getHeader,
  getMethod,
  jsonResponse,
  makeFetchMock,
  noContent,
} from './helpers.js';

const pm = (id: string) => ({
  id,
  method_type: 'bank_account',
  currency: 'VES',
  label: null,
  details: { bank_code: '0102', account_number: '****', holder_name: 'X' },
  verified: false,
  verified_at: null,
  last_used_at: null,
  counterparty: { id: 'c1', name: 'C' },
  tesote_account: null,
  created_at: 't',
  updated_at: 't',
});

describe('V2 paymentMethods', () => {
  it('list serializes verified=true to "true"', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { items: [pm('p1')], has_more: false, limit: 50, offset: 0 }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.paymentMethods.list({ method_type: 'bank_account', verified: true });
    expect(calls[0]?.url).toContain('/v2/payment_methods');
    expect(calls[0]?.url).toContain('method_type=bank_account');
    expect(calls[0]?.url).toContain('verified=true');
  });

  it('list verified=false serializes "false"', async () => {
    const { fetch, calls } = makeFetchMock([
      jsonResponse(200, { items: [], has_more: false, limit: 50, offset: 0 }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.paymentMethods.list({ verified: false });
    expect(calls[0]?.url).toContain('verified=false');
  });

  it('get → 404 PAYMENT_METHOD_NOT_FOUND', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(404, { error_code: 'PAYMENT_METHOD_NOT_FOUND' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.paymentMethods.get('p1')).rejects.toBeInstanceOf(PaymentMethodNotFoundError);
  });

  it('create wraps body under payment_method', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(201, pm('p1'))]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.paymentMethods.create(
      {
        method_type: 'pago_movil',
        currency: 'VES',
        details: { phone_number: '+58...', identification_number: 'V12345' },
      },
      { idempotencyKey: 'PM1' },
    );
    expect(getMethod(callAt(calls, 0))).toBe('POST');
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toBe('PM1');
    const body = getBody(callAt(calls, 0)) as { payment_method: { method_type: string } };
    expect(body.payment_method.method_type).toBe('pago_movil');
  });

  it('create → 400 VALIDATION_ERROR', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(400, { error_code: 'VALIDATION_ERROR', error: 'bad input' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(
      c.paymentMethods.create({ method_type: 'wire', currency: 'VES', details: {} }),
    ).rejects.toBeInstanceOf(ValidationError);
  });

  it('update PATCHes', async () => {
    const { fetch, calls } = makeFetchMock([jsonResponse(200, pm('p1'))]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await c.paymentMethods.update('p1', { label: 'New label' }, { idempotencyKey: 'U1' });
    expect(getMethod(callAt(calls, 0))).toBe('PATCH');
    expect(calls[0]?.url).toContain('/v2/payment_methods/p1');
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toBe('U1');
    const body = getBody(callAt(calls, 0)) as { payment_method: { label: string } };
    expect(body.payment_method.label).toBe('New label');
  });

  it('delete returns void on 204', async () => {
    const { fetch, calls } = makeFetchMock([noContent(204)]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.paymentMethods.delete('p1')).resolves.toBeUndefined();
    expect(getMethod(callAt(calls, 0))).toBe('DELETE');
    expect(getHeader(callAt(calls, 0), 'idempotency-key')).toMatch(/^[0-9a-f-]{36}$/);
  });

  it('delete → 409 VALIDATION_ERROR (in use)', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(409, { error_code: 'VALIDATION_ERROR', error: 'in use' }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    await expect(c.paymentMethods.delete('p1')).rejects.toBeInstanceOf(ValidationError);
  });

  it('listAll iterates offset pages', async () => {
    const { fetch } = makeFetchMock([
      jsonResponse(200, { items: [pm('p1'), pm('p2')], has_more: true, limit: 2, offset: 0 }),
      jsonResponse(200, { items: [pm('p3')], has_more: false, limit: 2, offset: 2 }),
    ]);
    const c = new V2Client({ apiKey: 'k', fetch });
    const ids: string[] = [];
    for await (const m of c.paymentMethods.listAll({ limit: 2 })) ids.push(m.id);
    expect(ids).toEqual(['p1', 'p2', 'p3']);
  });
});
