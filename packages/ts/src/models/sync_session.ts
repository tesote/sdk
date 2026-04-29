/**
 * SyncSession — record of a single bank-sync attempt.
 * From GET /v2/accounts/{id}/sync_sessions and the POST /sync response.
 */

export type SyncSessionStatus = 'pending' | 'started' | 'completed' | 'failed' | 'skipped';

export interface SyncSessionError {
  type: string;
  message: string;
}

export interface SyncSessionPerformance {
  total_duration: number;
  complexity_score: number;
  sync_speed_score: number;
}

export interface SyncSession {
  id: string;
  status: SyncSessionStatus;
  started_at: string;
  completed_at: string | null;
  transactions_synced: number;
  accounts_count: number;
  /** Only present when status === 'failed'. */
  error: SyncSessionError | null;
  /** Optional metrics from the SpeedMetric association. */
  performance: SyncSessionPerformance | null;
}

export interface SyncStartResponse {
  message: string;
  sync_session_id: string;
  status: SyncSessionStatus;
  started_at: string;
}

export interface SyncSessionListResponse {
  sync_sessions: SyncSession[];
  limit: number;
  offset: number;
  has_more: boolean;
}
