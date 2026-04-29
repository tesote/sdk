/**
 * SyncTransaction — v2 sync response shape (Plaid-compatible flattened model).
 * Returned in `added` / `modified` arrays by POST /v2/accounts/{id}/transactions/sync.
 */

export interface SyncTransaction {
  transaction_id: string;
  account_id: string;
  amount: number;
  iso_currency_code: string;
  unofficial_currency_code: string;
  date: string;
  datetime: string | null;
  name: string;
  merchant_name: string | null;
  pending: boolean;
  category: string[];
  /** Conditional: only when running balances are enabled. */
  running_balance_cents?: number;
}

export interface SyncRemoved {
  transaction_id: string;
  account_id: string;
}

export interface SyncResult {
  added: SyncTransaction[];
  modified: SyncTransaction[];
  removed: SyncRemoved[];
  next_cursor: string | null;
  has_more: boolean;
}
