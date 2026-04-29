import type {
  Beneficiary,
  TransactionOrder,
  TransactionOrderListResponse,
  TransactionOrderStatus,
} from '../models/transaction_order.js';
import type { Transport } from '../transport.js';

export interface TransactionOrderListParams {
  limit?: number;
  offset?: number;
  status?: TransactionOrderStatus;
  created_after?: string;
  created_before?: string;
  batch_id?: string;
}

export interface TransactionOrderCreateInput {
  destination_payment_method_id?: string | null;
  beneficiary?: Beneficiary;
  amount: string;
  currency: string;
  description: string;
  scheduled_for?: string | null;
  idempotency_key?: string | null;
  metadata?: Record<string, unknown>;
}

export interface TransactionOrderSubmitOptions {
  /** Token (e.g. MFA/OTP) some banks require. */
  token?: string | null;
  /** Idempotency key for the submit request itself. */
  idempotencyKey?: string;
}

export class V2TransactionOrdersClient {
  constructor(private readonly transport: Transport) {}

  /** GET /v2/accounts/{id}/transaction_orders */
  async list(
    accountId: string,
    params: TransactionOrderListParams = {},
  ): Promise<TransactionOrderListResponse> {
    const res = await this.transport.request<TransactionOrderListResponse>({
      method: 'GET',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/transaction_orders`,
      query: { ...params },
    });
    return res.data;
  }

  async *listAll(
    accountId: string,
    params: TransactionOrderListParams = {},
  ): AsyncGenerator<TransactionOrder, void, void> {
    let offset = params.offset ?? 0;
    const limit = params.limit ?? 50;
    while (true) {
      const page = await this.list(accountId, { ...params, limit, offset });
      for (const order of page.items) yield order;
      if (!page.has_more) return;
      offset += page.items.length;
    }
  }

  /** GET /v2/accounts/{id}/transaction_orders/{order_id} */
  async get(accountId: string, orderId: string): Promise<TransactionOrder> {
    const res = await this.transport.request<TransactionOrder>({
      method: 'GET',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/transaction_orders/${encodeURIComponent(orderId)}`,
    });
    return res.data;
  }

  /** POST /v2/accounts/{id}/transaction_orders */
  async create(
    accountId: string,
    body: TransactionOrderCreateInput,
    opts: { idempotencyKey?: string } = {},
  ): Promise<TransactionOrder> {
    const res = await this.transport.request<TransactionOrder>({
      method: 'POST',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/transaction_orders`,
      body: { transaction_order: body },
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }

  /** POST /v2/accounts/{id}/transaction_orders/{order_id}/submit */
  async submit(
    accountId: string,
    orderId: string,
    opts: TransactionOrderSubmitOptions = {},
  ): Promise<TransactionOrder> {
    const body: Record<string, unknown> = {};
    if (opts.token !== undefined) body.token = opts.token;
    const res = await this.transport.request<TransactionOrder>({
      method: 'POST',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/transaction_orders/${encodeURIComponent(orderId)}/submit`,
      body,
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }

  /** POST /v2/accounts/{id}/transaction_orders/{order_id}/cancel */
  async cancel(
    accountId: string,
    orderId: string,
    opts: { idempotencyKey?: string } = {},
  ): Promise<TransactionOrder> {
    const res = await this.transport.request<TransactionOrder>({
      method: 'POST',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/transaction_orders/${encodeURIComponent(orderId)}/cancel`,
      body: {},
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }
}
