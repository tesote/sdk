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

export {
  AccountDisabledError,
  AccountNotFoundError,
  ApiError,
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
  mapApiError,
  NetworkError,
  NotFoundError,
  PaymentMethodNotFoundError,
  RateLimitExceededError,
  type RequestSummary,
  redactBearer,
  ServiceUnavailableError,
  SyncInProgressError,
  SyncRateLimitExceededError,
  SyncSessionNotFoundError,
  TesoteError,
  type TesoteErrorFields,
  TimeoutError,
  TlsError,
  TransactionNotFoundError,
  TransactionOrderNotFoundError,
  TransportError,
  UnauthorizedError,
  UnprocessableContentError,
  ValidationError,
  WorkspaceSuspendedError,
} from './errors.js';

export type {
  Account,
  AccountBank,
  AccountData,
  AccountLegalEntity,
  AccountListResponse,
  PageBasedPagination,
} from './models/account.js';
export type {
  BatchApproveResponse,
  BatchCancelResponse,
  BatchCreateResponse,
  BatchStatus,
  BatchStatusCounts,
  BatchSubmitResponse,
  BatchSummary,
} from './models/batch.js';
export type { BulkResponse, BulkResult } from './models/bulk.js';
export type { OffsetPaginationResponse } from './models/pagination.js';
export type {
  PaymentMethod,
  PaymentMethodCounterparty,
  PaymentMethodDetails,
  PaymentMethodTesoteAccount,
  PaymentMethodType,
} from './models/payment_method.js';
export type { SearchResult } from './models/search.js';
export type { StatusResponse, WhoamiClient, WhoamiResponse } from './models/status.js';
export type {
  SyncSession,
  SyncSessionError,
  SyncSessionPerformance,
} from './models/sync_session.js';
export type {
  SyncRemoved,
  SyncResult,
  SyncTransaction,
} from './models/sync_transaction.js';
export type {
  CursorPagination,
  Transaction,
  TransactionCategory,
  TransactionCounterparty,
  TransactionData,
  TransactionListResponse,
} from './models/transaction.js';
export type {
  Beneficiary,
  TransactionOrder,
  TransactionOrderDestination,
  TransactionOrderFee,
  TransactionOrderLatestAttempt,
  TransactionOrderSourceAccount,
  TransactionOrderStatus,
  TransactionOrderTesoteTransaction,
} from './models/transaction_order.js';
export {
  type CacheBackend,
  type CacheEntry,
  type CacheOptions,
  DEFAULT_BASE_URL,
  DEFAULT_RETRY_POLICY,
  InMemoryLRUCache,
  type LogEvent,
  type LogHook,
  type RateLimitSnapshot,
  type RequestOptions,
  type ResponseEnvelope,
  type RetryPolicy,
  SDK_VERSION,
  Transport,
  type TransportOptions,
} from './transport.js';
export { V1Client, type V1ClientOptions } from './v1/index.js';
export { V2Client, type V2ClientOptions } from './v2/index.js';
