import { describe, expect, it } from 'vitest';
import {
  AccountDisabledError,
  ApiError,
  ApiKeyRevokedError,
  HistorySyncForbiddenError,
  InvalidDateRangeError,
  MutationDuringPaginationError,
  RateLimitExceededError,
  type RequestSummary,
  ServiceUnavailableError,
  type TesoteError,
  UnauthorizedError,
  UnprocessableContentError,
  WorkspaceSuspendedError,
  mapApiError,
  redactBearer,
} from '../src/errors.js';

const SUMMARY: RequestSummary = {
  method: 'GET',
  path: '/v3/accounts',
  authorization: 'Bearer abcd',
};

const cases: Array<{
  code: string;
  status: number;
  cls: typeof TesoteError;
}> = [
  { code: 'UNAUTHORIZED', status: 401, cls: UnauthorizedError },
  { code: 'API_KEY_REVOKED', status: 401, cls: ApiKeyRevokedError },
  { code: 'WORKSPACE_SUSPENDED', status: 403, cls: WorkspaceSuspendedError },
  { code: 'ACCOUNT_DISABLED', status: 403, cls: AccountDisabledError },
  { code: 'HISTORY_SYNC_FORBIDDEN', status: 403, cls: HistorySyncForbiddenError },
  { code: 'MUTATION_CONFLICT', status: 409, cls: MutationDuringPaginationError },
  { code: 'UNPROCESSABLE_CONTENT', status: 422, cls: UnprocessableContentError },
  { code: 'INVALID_DATE_RANGE', status: 422, cls: InvalidDateRangeError },
  { code: 'RATE_LIMIT_EXCEEDED', status: 429, cls: RateLimitExceededError },
];

describe('mapApiError — error_code → typed class', () => {
  for (const { code, status, cls } of cases) {
    it(`${code} → ${cls.name}`, () => {
      const err = mapApiError({
        httpStatus: status,
        requestId: 'req-x',
        retryAfterHeader: null,
        responseBody: '{}',
        parsedBody: { error: 'msg', error_code: code, error_id: 'eid' },
        requestSummary: SUMMARY,
        attempts: 1,
      });
      expect(err).toBeInstanceOf(cls);
      expect(err).toBeInstanceOf(ApiError);
      expect(err.errorCode).toBe(code);
      expect(err.errorId).toBe('eid');
      expect(err.httpStatus).toBe(status);
      expect(err.requestId).toBe('req-x');
      expect(err.requestSummary).toBe(SUMMARY);
    });
  }

  it('falls back to ServiceUnavailableError on bare 503', () => {
    const err = mapApiError({
      httpStatus: 503,
      requestId: null,
      retryAfterHeader: '5',
      responseBody: '',
      parsedBody: null,
      requestSummary: SUMMARY,
      attempts: 2,
    });
    expect(err).toBeInstanceOf(ServiceUnavailableError);
    expect(err.retryAfter).toBe(5);
    expect(err.attempts).toBe(2);
  });

  it('prefers envelope retry_after over header', () => {
    const err = mapApiError({
      httpStatus: 429,
      requestId: null,
      retryAfterHeader: '10',
      responseBody: '',
      parsedBody: { error_code: 'RATE_LIMIT_EXCEEDED', retry_after: 42 },
      requestSummary: SUMMARY,
      attempts: 1,
    });
    expect(err.retryAfter).toBe(42);
  });

  it('preserves all required fields', () => {
    const err = mapApiError({
      httpStatus: 422,
      requestId: 'r1',
      retryAfterHeader: null,
      responseBody: '{"error":"x","error_code":"UNPROCESSABLE_CONTENT","error_id":"eid"}',
      parsedBody: { error: 'x', error_code: 'UNPROCESSABLE_CONTENT', error_id: 'eid' },
      requestSummary: SUMMARY,
      attempts: 3,
    });
    expect(err.errorCode).toBe('UNPROCESSABLE_CONTENT');
    expect(err.message).toBe('x');
    expect(err.httpStatus).toBe(422);
    expect(err.requestId).toBe('r1');
    expect(err.errorId).toBe('eid');
    expect(err.responseBody).toBe(
      '{"error":"x","error_code":"UNPROCESSABLE_CONTENT","error_id":"eid"}',
    );
    expect(err.requestSummary).toBe(SUMMARY);
    expect(err.attempts).toBe(3);
  });
});

describe('redactBearer', () => {
  it('keeps only the last 4 chars', () => {
    expect(redactBearer('sk_super_secret_abcd')).toBe('Bearer abcd');
  });

  it('redacts short tokens entirely', () => {
    expect(redactBearer('abc')).toBe('Bearer <redacted>');
  });
});
