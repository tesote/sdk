/**
 * Typed error hierarchy for the Tesote SDK.
 * Mirrors docs/architecture/errors.md exactly.
 */

export interface RequestSummary {
  method: string;
  path: string;
  query?: Record<string, string | number | boolean | null | undefined>;
  bodyShape?: string;
  /** Bearer token redacted to `Bearer <last4>`; never the full secret. */
  authorization?: string;
}

export interface TesoteErrorFields {
  errorCode: string;
  message: string;
  httpStatus: number | null;
  requestId: string | null;
  errorId: string | null;
  retryAfter: number | null;
  responseBody: string | null;
  requestSummary: RequestSummary | null;
  attempts: number;
  cause?: unknown;
}

const DEFAULT_FIELDS: Omit<TesoteErrorFields, 'errorCode' | 'message'> = {
  httpStatus: null,
  requestId: null,
  errorId: null,
  retryAfter: null,
  responseBody: null,
  requestSummary: null,
  attempts: 1,
};

/**
 * Base for everything the SDK throws. Catch this only as a last resort —
 * prefer the narrow subclasses.
 */
export class TesoteError extends Error {
  public readonly errorCode: string;
  public readonly httpStatus: number | null;
  public readonly requestId: string | null;
  public readonly errorId: string | null;
  public readonly retryAfter: number | null;
  public readonly responseBody: string | null;
  public readonly requestSummary: RequestSummary | null;
  public readonly attempts: number;

  constructor(
    fields: Partial<TesoteErrorFields> & Pick<TesoteErrorFields, 'errorCode' | 'message'>,
  ) {
    super(fields.message);
    this.name = new.target.name;
    this.errorCode = fields.errorCode;
    this.httpStatus = fields.httpStatus ?? DEFAULT_FIELDS.httpStatus;
    this.requestId = fields.requestId ?? DEFAULT_FIELDS.requestId;
    this.errorId = fields.errorId ?? DEFAULT_FIELDS.errorId;
    this.retryAfter = fields.retryAfter ?? DEFAULT_FIELDS.retryAfter;
    this.responseBody = fields.responseBody ?? DEFAULT_FIELDS.responseBody;
    this.requestSummary = fields.requestSummary ?? DEFAULT_FIELDS.requestSummary;
    this.attempts = fields.attempts ?? DEFAULT_FIELDS.attempts;
    if (fields.cause !== undefined) {
      // why: Node 18+ supports Error.cause natively; assign to keep the chain
      (this as { cause?: unknown }).cause = fields.cause;
    }
    // why: V8 stack-trace tweak; harmless on other engines (no-op on undefined).
    if (
      typeof (Error as unknown as { captureStackTrace?: unknown }).captureStackTrace === 'function'
    ) {
      (
        Error as unknown as { captureStackTrace: (t: object, c: unknown) => void }
      ).captureStackTrace(this, new.target);
    }
  }
}

/** Server returned a structured API error. */
export class ApiError extends TesoteError {}

// 401
export class UnauthorizedError extends ApiError {}
export class ApiKeyRevokedError extends ApiError {}

// 403
export class WorkspaceSuspendedError extends ApiError {}
export class AccountDisabledError extends ApiError {}
export class HistorySyncForbiddenError extends ApiError {}

// 404
export class NotFoundError extends ApiError {}
export class AccountNotFoundError extends NotFoundError {}
export class TransactionNotFoundError extends NotFoundError {}
export class SyncSessionNotFoundError extends NotFoundError {}
export class PaymentMethodNotFoundError extends NotFoundError {}
export class TransactionOrderNotFoundError extends NotFoundError {}
export class BatchNotFoundError extends NotFoundError {}
export class BankConnectionNotFoundError extends NotFoundError {}

// 409
export class MutationDuringPaginationError extends ApiError {}
export class InvalidOrderStateError extends ApiError {}
export class SyncInProgressError extends ApiError {}

// 422
export class UnprocessableContentError extends ApiError {}
export class InvalidDateRangeError extends ApiError {}
export class InvalidCursorError extends ApiError {}
export class InvalidCountError extends ApiError {}
export class InvalidLimitError extends ApiError {}
export class InvalidQueryError extends ApiError {}
export class MissingDateRangeError extends ApiError {}
export class BankSubmissionError extends ApiError {}

// 400
export class ValidationError extends ApiError {}
export class BatchValidationError extends ApiError {}

// 429
export class RateLimitExceededError extends ApiError {}
export class SyncRateLimitExceededError extends ApiError {}

// 503
export class ServiceUnavailableError extends ApiError {}
export class BankUnderMaintenanceError extends ApiError {}

// 500
export class InternalServerError extends ApiError {}

/** Transport-level failure: no usable HTTP response. */
export class TransportError extends TesoteError {}
export class NetworkError extends TransportError {}
export class TimeoutError extends TransportError {}
export class TlsError extends TransportError {}

/** Bad SDK config; raised at client construction. */
export class ConfigError extends TesoteError {}

/** Method exists in the SDK but the upstream endpoint is gone in this API version. */
export class EndpointRemovedError extends TesoteError {}

interface ApiErrorEnvelope {
  error?: string;
  error_code?: string;
  error_id?: string;
  retry_after?: number;
}

type ApiErrorCtor = new (
  f: Partial<TesoteErrorFields> & Pick<TesoteErrorFields, 'errorCode' | 'message'>,
) => ApiError;

const ERROR_CODE_MAP: Record<string, ApiErrorCtor> = {
  // 401
  UNAUTHORIZED: UnauthorizedError,
  API_KEY_REVOKED: ApiKeyRevokedError,
  // 403
  WORKSPACE_SUSPENDED: WorkspaceSuspendedError,
  ACCOUNT_DISABLED: AccountDisabledError,
  HISTORY_SYNC_FORBIDDEN: HistorySyncForbiddenError,
  // 404
  ACCOUNT_NOT_FOUND: AccountNotFoundError,
  TRANSACTION_NOT_FOUND: TransactionNotFoundError,
  SYNC_SESSION_NOT_FOUND: SyncSessionNotFoundError,
  PAYMENT_METHOD_NOT_FOUND: PaymentMethodNotFoundError,
  TRANSACTION_ORDER_NOT_FOUND: TransactionOrderNotFoundError,
  BATCH_NOT_FOUND: BatchNotFoundError,
  BANK_CONNECTION_NOT_FOUND: BankConnectionNotFoundError,
  // 409
  MUTATION_CONFLICT: MutationDuringPaginationError,
  INVALID_ORDER_STATE: InvalidOrderStateError,
  SYNC_IN_PROGRESS: SyncInProgressError,
  // 422
  UNPROCESSABLE_CONTENT: UnprocessableContentError,
  INVALID_DATE_RANGE: InvalidDateRangeError,
  INVALID_CURSOR: InvalidCursorError,
  INVALID_COUNT: InvalidCountError,
  INVALID_LIMIT: InvalidLimitError,
  INVALID_QUERY: InvalidQueryError,
  MISSING_DATE_RANGE: MissingDateRangeError,
  BANK_SUBMISSION_ERROR: BankSubmissionError,
  // 400
  VALIDATION_ERROR: ValidationError,
  BATCH_VALIDATION_ERROR: BatchValidationError,
  // 429
  RATE_LIMIT_EXCEEDED: RateLimitExceededError,
  SYNC_RATE_LIMIT_EXCEEDED: SyncRateLimitExceededError,
  // 503
  BANK_UNDER_MAINTENANCE: BankUnderMaintenanceError,
  // 500
  INTERNAL_ERROR: InternalServerError,
};

function pickStatusFallback(httpStatus: number): {
  cls: ApiErrorCtor;
  errorCode: string;
} {
  if (httpStatus === 401) return { cls: UnauthorizedError, errorCode: 'UNAUTHORIZED' };
  if (httpStatus === 403) return { cls: ApiError, errorCode: 'FORBIDDEN' };
  if (httpStatus === 404) return { cls: NotFoundError, errorCode: 'NOT_FOUND' };
  if (httpStatus === 409)
    return { cls: MutationDuringPaginationError, errorCode: 'MUTATION_CONFLICT' };
  if (httpStatus === 422)
    return { cls: UnprocessableContentError, errorCode: 'UNPROCESSABLE_CONTENT' };
  if (httpStatus === 429) return { cls: RateLimitExceededError, errorCode: 'RATE_LIMIT_EXCEEDED' };
  if (httpStatus === 500) return { cls: InternalServerError, errorCode: 'INTERNAL_ERROR' };
  if (httpStatus === 503) return { cls: ServiceUnavailableError, errorCode: 'SERVICE_UNAVAILABLE' };
  return { cls: ApiError, errorCode: `HTTP_${httpStatus}` };
}

export interface MapApiErrorInput {
  httpStatus: number;
  requestId: string | null;
  retryAfterHeader: string | null;
  responseBody: string | null;
  parsedBody: unknown;
  requestSummary: RequestSummary;
  attempts: number;
}

function parseRetryAfter(header: string | null, envelope: number | null): number | null {
  if (envelope !== null && Number.isFinite(envelope)) return envelope;
  if (header === null) return null;
  const n = Number(header);
  return Number.isFinite(n) ? n : null;
}

function envelopeFrom(parsed: unknown): ApiErrorEnvelope {
  if (parsed === null || typeof parsed !== 'object') return {};
  const obj = parsed as Record<string, unknown>;
  const env: ApiErrorEnvelope = {};
  if (typeof obj.error === 'string') env.error = obj.error;
  if (typeof obj.error_code === 'string') env.error_code = obj.error_code;
  if (typeof obj.error_id === 'string') env.error_id = obj.error_id;
  if (typeof obj.retry_after === 'number') env.retry_after = obj.retry_after;
  return env;
}

/**
 * Dispatch an API response (4xx/5xx) into a typed error.
 */
export function mapApiError(input: MapApiErrorInput): ApiError {
  const env = envelopeFrom(input.parsedBody);
  const code = env.error_code;
  let cls: ApiErrorCtor;
  let errorCode: string;
  if (code !== undefined && code in ERROR_CODE_MAP) {
    const mapped = ERROR_CODE_MAP[code];
    if (mapped === undefined) {
      const fb = pickStatusFallback(input.httpStatus);
      cls = fb.cls;
      errorCode = code;
    } else {
      cls = mapped;
      errorCode = code;
    }
  } else if (input.httpStatus === 503) {
    cls = ServiceUnavailableError;
    errorCode = code ?? 'SERVICE_UNAVAILABLE';
  } else {
    const fb = pickStatusFallback(input.httpStatus);
    cls = fb.cls;
    errorCode = code ?? fb.errorCode;
  }

  const message = env.error ?? `${input.httpStatus} ${cls.name.replace(/Error$/, '')}`;

  return new cls({
    errorCode,
    message,
    httpStatus: input.httpStatus,
    requestId: input.requestId,
    errorId: env.error_id ?? null,
    retryAfter: parseRetryAfter(input.retryAfterHeader, env.retry_after ?? null),
    responseBody: input.responseBody,
    requestSummary: input.requestSummary,
    attempts: input.attempts,
  });
}

/**
 * Redact a bearer token to `Bearer <last4>` for safe logging.
 * Empty / short tokens collapse to `Bearer <redacted>`.
 */
export function redactBearer(apiKey: string): string {
  if (apiKey.length < 4) return 'Bearer <redacted>';
  return `Bearer ${apiKey.slice(-4)}`;
}
