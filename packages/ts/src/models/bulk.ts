/**
 * BulkResult — POST /v2/transactions/bulk response shape.
 * Returns transactions for multiple accounts in a single round-trip.
 */

import type { CursorPagination, Transaction } from './transaction.js';

export interface BulkResult {
  account_id: string;
  transactions: Transaction[];
  pagination: CursorPagination;
}

export interface BulkResponse {
  bulk_results: BulkResult[];
}
