/**
 * Batch — wraps multiple TransactionOrder rows for atomic creation/lifecycle.
 */

import type { TransactionOrder } from './transaction_order.js';

export type BatchStatus = 'draft' | 'mixed' | 'approved' | 'processing' | 'completed' | string;

export interface BatchStatusCounts {
  draft?: number;
  pending_approval?: number;
  approved?: number;
  processing?: number;
  completed?: number;
  failed?: number;
  cancelled?: number;
  [k: string]: number | undefined;
}

export interface BatchSummary {
  batch_id: string;
  total_orders: number;
  total_amount_cents: number;
  amount_currency: string;
  statuses: BatchStatusCounts;
  batch_status: BatchStatus;
  created_at: string;
  orders: TransactionOrder[];
}

export interface BatchCreateError {
  /** Index of the order in the request array that failed validation. */
  index: number;
  message: string;
  details?: unknown;
}

export interface BatchCreateResponse {
  batch_id: string;
  orders: TransactionOrder[];
  errors: BatchCreateError[];
}

export interface BatchApproveResponse {
  approved: number;
  failed: number;
}

export interface BatchSubmitResponse {
  enqueued: number;
  failed: number;
}

export interface BatchCancelResponse {
  cancelled: number;
  skipped: number;
  errors: BatchCreateError[];
}
