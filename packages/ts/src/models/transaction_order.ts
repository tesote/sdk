/**
 * TransactionOrder — outbound payment lifecycle object.
 * State machine: draft -> pending_approval -> approved -> processing -> completed|failed|cancelled.
 */

export type TransactionOrderStatus =
  | 'draft'
  | 'pending_approval'
  | 'approved'
  | 'processing'
  | 'completed'
  | 'failed'
  | 'cancelled';

export interface Beneficiary {
  name: string;
  bank_code?: string | null;
  account_number?: string | null;
  identification_type?: string | null;
  identification_number?: string | null;
}

export interface TransactionOrderSourceAccount {
  id: string;
  name: string;
  payment_method_id: string;
}

export interface TransactionOrderDestination {
  payment_method_id: string;
  counterparty_id: string;
  counterparty_name: string;
}

export interface TransactionOrderFee {
  amount: number;
  currency: string;
}

export interface TransactionOrderTesoteTransaction {
  id: string;
  status: string;
}

export interface TransactionOrderLatestAttempt {
  id: string;
  status: string;
  attempt_number: number;
  external_reference: string | null;
  submitted_at: string | null;
  completed_at: string | null;
  error_code: string | null;
  error_message: string | null;
}

export interface TransactionOrder {
  id: string;
  status: TransactionOrderStatus;
  amount: number;
  currency: string;
  description: string;
  reference: string | null;
  external_reference: string | null;
  idempotency_key: string | null;
  batch_id: string | null;
  scheduled_for: string | null;
  approved_at: string | null;
  submitted_at: string | null;
  completed_at: string | null;
  failed_at: string | null;
  cancelled_at: string | null;
  source_account: TransactionOrderSourceAccount;
  destination: TransactionOrderDestination;
  /** Null when fee_cents is zero. */
  fee: TransactionOrderFee | null;
  execution_strategy: string | null;
  tesote_transaction: TransactionOrderTesoteTransaction | null;
  latest_attempt: TransactionOrderLatestAttempt | null;
  created_at: string;
  updated_at: string;
}

export interface TransactionOrderListResponse {
  items: TransactionOrder[];
  has_more: boolean;
  limit: number;
  offset: number;
}
