/**
 * @tesote.com/sdk — official TypeScript SDK for the Tesote API.
 *
 * Versioned clients ship side-by-side. Pick the version explicitly:
 *
 * ```ts
 * import { V1Client, V2Client } from '@tesote.com/sdk';
 * const v2 = new V2Client({ apiKey: process.env.TESOTE_API_KEY! });
 * await v2.accounts.list();
 * ```
 */

export { V1Client, type V1ClientOptions } from './v1/index.js';
export { V2Client, type V2ClientOptions } from './v2/index.js';

export {
  ApiError,
  AccountDisabledError,
  AccountNotFoundError,
  ApiKeyRevokedError,
  BankConnectionNotFoundError,
  BankSubmissionError,
  BankUnderMaintenanceError,
  BatchNotFoundError,
  BatchValidationError,
  ConfigError,
  EndpointRemovedError,
  HistorySyncForbiddenError,
  InternalServerError,
  InvalidCountError,
  InvalidCursorError,
  InvalidDateRangeError,
  InvalidLimitError,
  InvalidOrderStateError,
  InvalidQueryError,
  MissingDateRangeError,
  MutationDuringPaginationError,
  NetworkError,
  NotFoundError,
  PaymentMethodNotFoundError,
  RateLimitExceededError,
  ServiceUnavailableError,
  SyncInProgressError,
  SyncRateLimitExceededError,
  SyncSessionNotFoundError,
  TesoteError,
  TimeoutError,
  TlsError,
  TransactionNotFoundError,
  TransactionOrderNotFoundError,
  TransportError,
  UnauthorizedError,
  UnprocessableContentError,
  ValidationError,
  WorkspaceSuspendedError,
  mapApiError,
  redactBearer,
  type RequestSummary,
  type TesoteErrorFields,
} from './errors.js';

export {
  Transport,
  InMemoryLRUCache,
  DEFAULT_BASE_URL,
  DEFAULT_RETRY_POLICY,
  SDK_VERSION,
  type CacheBackend,
  type CacheEntry,
  type CacheOptions,
  type LogEvent,
  type LogHook,
  type RateLimitSnapshot,
  type RequestOptions,
  type ResponseEnvelope,
  type RetryPolicy,
  type TransportOptions,
} from './transport.js';

export type {
  Account,
  AccountData,
  AccountBank,
  AccountLegalEntity,
  AccountListResponse,
  PageBasedPagination,
} from './models/account.js';
export type {
  Transaction,
  TransactionData,
  TransactionCategory,
  TransactionCounterparty,
  TransactionListResponse,
  CursorPagination,
} from './models/transaction.js';
export type {
  SyncTransaction,
  SyncRemoved,
  SyncResult,
} from './models/sync_transaction.js';
export type {
  SyncSession,
  SyncSessionError,
  SyncSessionPerformance,
} from './models/sync_session.js';
export type {
  TransactionOrder,
  TransactionOrderStatus,
  TransactionOrderSourceAccount,
  TransactionOrderDestination,
  TransactionOrderFee,
  TransactionOrderTesoteTransaction,
  TransactionOrderLatestAttempt,
  Beneficiary,
} from './models/transaction_order.js';
export type {
  PaymentMethod,
  PaymentMethodType,
  PaymentMethodDetails,
  PaymentMethodCounterparty,
  PaymentMethodTesoteAccount,
} from './models/payment_method.js';
export type {
  BatchSummary,
  BatchStatus,
  BatchStatusCounts,
  BatchCreateResponse,
  BatchApproveResponse,
  BatchSubmitResponse,
  BatchCancelResponse,
} from './models/batch.js';
export type { BulkResult, BulkResponse } from './models/bulk.js';
export type { SearchResult } from './models/search.js';
export type { OffsetPaginationResponse } from './models/pagination.js';
export type { StatusResponse, WhoamiResponse, WhoamiClient } from './models/status.js';
