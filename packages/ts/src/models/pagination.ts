/**
 * Generic offset-pagination response envelope.
 * Used by sync_sessions, transaction_orders, payment_methods.
 */

export interface OffsetPaginationResponse<T> {
  items: T[];
  has_more: boolean;
  limit: number;
  offset: number;
}
