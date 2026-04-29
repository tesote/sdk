/**
 * Transaction model — v1 schema, also used by GET /v2/transactions/{id}.
 * Wire shape uses snake_case; we preserve it on the model.
 */

export interface TransactionCategory {
  name: string;
  external_category_code: string | null;
  created_at: string;
  updated_at: string;
}

export interface TransactionCounterparty {
  name: string;
}

export interface TransactionData {
  amount_cents: number;
  currency: string;
  description: string;
  transaction_date: string;
  created_at: string | null;
  created_at_date: string | null;
  note: string | null;
  external_service_id: string | null;
  /** Conditional: only present when running balances are enabled for the workspace. */
  running_balance_cents?: number;
}

export type TransactionStatus = 'posted' | 'pending' | 'failed' | string;

export interface Transaction {
  id: string;
  status: TransactionStatus;
  data: TransactionData;
  tesote_imported_at: string;
  tesote_updated_at: string;
  transaction_categories: TransactionCategory[];
  counterparty: TransactionCounterparty | null;
}

export interface CursorPagination {
  has_more: boolean;
  per_page: number;
  after_id: string | null;
  before_id: string | null;
}

export interface TransactionListResponse {
  total: number;
  transactions: Transaction[];
  pagination: CursorPagination;
}
