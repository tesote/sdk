/**
 * PaymentMethod — destination/source for transactions.
 * `method_type` determines which fields the `details` object carries.
 */

export type PaymentMethodType =
  | 'bank_account'
  | 'pago_movil'
  | 'wire'
  | 'crypto_wallet'
  | 'fx_rail'
  | 'ach'
  | 'eft';

export interface PaymentMethodDetails {
  bank_code?: string;
  account_number?: string;
  holder_name?: string;
  identification_type?: string | null;
  identification_number?: string | null;
  /** method_type-specific extras (wallet address, IBAN, routing number, ...). */
  [field: string]: string | number | boolean | null | undefined;
}

export interface PaymentMethodCounterparty {
  id: string;
  name: string;
}

export interface PaymentMethodTesoteAccount {
  id: string;
  name: string;
}

export interface PaymentMethod {
  id: string;
  method_type: PaymentMethodType;
  currency: string;
  label: string | null;
  details: PaymentMethodDetails;
  verified: boolean;
  verified_at: string | null;
  last_used_at: string | null;
  /** Mutually exclusive with `tesote_account` — set when the method is a destination. */
  counterparty: PaymentMethodCounterparty | null;
  /** Mutually exclusive with `counterparty` — set when the method is a source. */
  tesote_account: PaymentMethodTesoteAccount | null;
  created_at: string;
  updated_at: string;
}

export interface PaymentMethodListResponse {
  items: PaymentMethod[];
  has_more: boolean;
  limit: number;
  offset: number;
}
