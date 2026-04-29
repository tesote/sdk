/**
 * SearchResult — GET /v2/transactions/search response shape.
 */

import type { Transaction } from './transaction.js';

export interface SearchResult {
  transactions: Transaction[];
  total: number;
}
