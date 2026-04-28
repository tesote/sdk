/**
 * @tesote/sdk — official TypeScript SDK for the Tesote API.
 *
 * Versioned clients ship side-by-side. Pick the version explicitly:
 *
 * ```ts
 * import { V1Client, V2Client } from '@tesote/sdk';
 * const v2 = new V2Client({ apiKey: process.env.TESOTE_API_KEY! });
 * await v2.accounts.list();
 * ```
 */

export { V1Client, type V1ClientOptions } from './v1/index.js';
export { V2Client, type V2ClientOptions } from './v2/index.js';

export {
  ApiError,
  AccountDisabledError,
  ApiKeyRevokedError,
  ConfigError,
  EndpointRemovedError,
  HistorySyncForbiddenError,
  InvalidDateRangeError,
  MutationDuringPaginationError,
  NetworkError,
  RateLimitExceededError,
  ServiceUnavailableError,
  TesoteError,
  TimeoutError,
  TlsError,
  TransportError,
  UnauthorizedError,
  UnprocessableContentError,
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
