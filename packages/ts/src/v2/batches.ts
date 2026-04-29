import type {
  BatchApproveResponse,
  BatchCancelResponse,
  BatchCreateResponse,
  BatchSubmitResponse,
  BatchSummary,
} from '../models/batch.js';
import type { Beneficiary } from '../models/transaction_order.js';
import type { Transport } from '../transport.js';

export interface BatchOrderInput {
  destination_payment_method_id?: string | null;
  beneficiary?: Beneficiary;
  amount: string;
  currency: string;
  description: string;
  scheduled_for?: string | null;
  metadata?: Record<string, unknown>;
}

export interface BatchCreateInput {
  orders: BatchOrderInput[];
}

export class V2BatchesClient {
  constructor(private readonly transport: Transport) {}

  /** POST /v2/accounts/{id}/batches */
  async create(
    accountId: string,
    body: BatchCreateInput,
    opts: { idempotencyKey?: string } = {},
  ): Promise<BatchCreateResponse> {
    const res = await this.transport.request<BatchCreateResponse>({
      method: 'POST',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/batches`,
      body,
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }

  /** GET /v2/accounts/{id}/batches/{batch_id} */
  async get(accountId: string, batchId: string): Promise<BatchSummary> {
    const res = await this.transport.request<BatchSummary>({
      method: 'GET',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/batches/${encodeURIComponent(batchId)}`,
    });
    return res.data;
  }

  /** POST /v2/accounts/{id}/batches/{batch_id}/approve */
  async approve(
    accountId: string,
    batchId: string,
    opts: { idempotencyKey?: string } = {},
  ): Promise<BatchApproveResponse> {
    const res = await this.transport.request<BatchApproveResponse>({
      method: 'POST',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/batches/${encodeURIComponent(batchId)}/approve`,
      body: {},
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }

  /** POST /v2/accounts/{id}/batches/{batch_id}/submit */
  async submit(
    accountId: string,
    batchId: string,
    opts: { token?: string | null; idempotencyKey?: string } = {},
  ): Promise<BatchSubmitResponse> {
    const body: Record<string, unknown> = {};
    if (opts.token !== undefined) body.token = opts.token;
    const res = await this.transport.request<BatchSubmitResponse>({
      method: 'POST',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/batches/${encodeURIComponent(batchId)}/submit`,
      body,
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }

  /** POST /v2/accounts/{id}/batches/{batch_id}/cancel */
  async cancel(
    accountId: string,
    batchId: string,
    opts: { idempotencyKey?: string } = {},
  ): Promise<BatchCancelResponse> {
    const res = await this.transport.request<BatchCancelResponse>({
      method: 'POST',
      path: `/v2/accounts/${encodeURIComponent(accountId)}/batches/${encodeURIComponent(batchId)}/cancel`,
      body: {},
      ...(opts.idempotencyKey !== undefined ? { idempotencyKey: opts.idempotencyKey } : {}),
    });
    return res.data;
  }
}
