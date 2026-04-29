import type {
  PaymentMethod,
  PaymentMethodDetails,
  PaymentMethodListResponse,
  PaymentMethodType,
} from '../models/payment_method.js';
import type { Transport } from '../transport.js';

export interface PaymentMethodListParams {
  limit?: number;
  offset?: number;
  method_type?: PaymentMethodType;
  currency?: string;
  counterparty_id?: string;
  /** Stringified boolean — the API takes "true" / "false" strings. */
  verified?: boolean;
}

export interface PaymentMethodCreateInput {
  method_type: PaymentMethodType;
  currency: string;
  label?: string | null;
  counterparty_id?: string | null;
  counterparty?: { name: string };
  details: PaymentMethodDetails;
}

export interface PaymentMethodUpdateInput {
  method_type?: PaymentMethodType;
  currency?: string;
  label?: string | null;
  counterparty_id?: string | null;
  counterparty?: { name: string };
  details?: Partial<PaymentMethodDetails>;
}

export class V2PaymentMethodsClient {
  constructor(private readonly transport: Transport) {}

  /** GET /v2/payment_methods */
  async list(params: PaymentMethodListParams = {}): Promise<PaymentMethodListResponse> {
    const query: Record<string, string | number | boolean | null | undefined> = {};
    if (params.limit !== undefined) query.limit = params.limit;
    if (params.offset !== undefined) query.offset = params.offset;
    if (params.method_type !== undefined) query.method_type = params.method_type;
    if (params.currency !== undefined) query.currency = params.currency;
    if (params.counterparty_id !== undefined) query.counterparty_id = params.counterparty_id;
    if (params.verified !== undefined) query.verified = params.verified ? 'true' : 'false';
    const res = await this.transport.request<PaymentMethodListResponse>({
      method: 'GET',
      path: '/v2/payment_methods',
      query,
    });
    return res.data;
  }

  async *listAll(params: PaymentMethodListParams = {}): AsyncGenerator<PaymentMethod, void, void> {
    let offset = params.offset ?? 0;
    const limit = params.limit ?? 50;
    while (true) {
      const page = await this.list({ ...params, limit, offset });
      for (const pm of page.items) yield pm;
      if (!page.has_more) return;
      offset += page.items.length;
    }
  }

  /** GET /v2/payment_methods/{id} */
  async get(id: string): Promise<PaymentMethod> {
    const res = await this.transport.request<PaymentMethod>({
      method: 'GET',
      path: `/v2/payment_methods/${encodeURIComponent(id)}`,
    });
    return res.data;
  }

  /** POST /v2/payment_methods */
  async create(
    body: PaymentMethodCreateInput,
    opts: { idempotencyKey?: string } = {},
  ): Promise<PaymentMethod> {
    const res = await this.transport.request<PaymentMethod>({
      method: 'POST',
      path: '/v2/payment_methods',
      body: { payment_method: body },
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }

  /** PATCH /v2/payment_methods/{id} */
  async update(
    id: string,
    body: PaymentMethodUpdateInput,
    opts: { idempotencyKey?: string } = {},
  ): Promise<PaymentMethod> {
    const res = await this.transport.request<PaymentMethod>({
      method: 'PATCH',
      path: `/v2/payment_methods/${encodeURIComponent(id)}`,
      body: { payment_method: body },
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }

  /** DELETE /v2/payment_methods/{id} — returns 204; resolves to void. */
  async delete(id: string, opts: { idempotencyKey?: string } = {}): Promise<void> {
    await this.transport.request<unknown>({
      method: 'DELETE',
      path: `/v2/payment_methods/${encodeURIComponent(id)}`,
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
  }
}
