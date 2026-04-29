/**
 * Account model — identical between v1 and v2.
 * Wire shape uses snake_case; we preserve it on the model.
 */

export interface AccountBank {
  name: string;
}

export interface AccountLegalEntity {
  id: string | null;
  legal_name: string | null;
}

export interface AccountData {
  masked_account_number: string;
  currency: string;
  transactions_data_current_as_of: string | null;
  balance_data_current_as_of: string | null;
  custom_user_provided_identifier: string | null;
  /** Conditional: only present if `display_balances_in_api` is enabled. */
  balance_cents?: string;
  /** Conditional: only present if `display_balances_in_api` is enabled. */
  available_balance_cents?: string;
}

export interface Account {
  id: string;
  name: string;
  data: AccountData;
  bank: AccountBank;
  legal_entity: AccountLegalEntity;
  tesote_created_at: string;
  tesote_updated_at: string;
}

export interface PageBasedPagination {
  current_page: number;
  per_page: number;
  total_pages: number;
  total_count: number;
}

export interface AccountListResponse {
  total: number;
  accounts: Account[];
  pagination: PageBasedPagination;
}
