"""Typed error hierarchy for the Tesote SDK.

Mirrors ``docs/architecture/errors.md``: every instance carries enough context
to debug without re-running the request. Bearer tokens are never persisted on
errors -- transport callers must redact before constructing ``request_summary``.
"""

from __future__ import annotations

from typing import Any, Dict, Mapping, Optional, Type


class TesoteError(Exception):
    """Base for everything the SDK raises.

    Catch this only as a last resort. Prefer the narrowest typed subclass.
    """

    def __init__(
        self,
        message: str,
        *,
        error_code: Optional[str] = None,
        http_status: Optional[int] = None,
        request_id: Optional[str] = None,
        error_id: Optional[str] = None,
        retry_after: Optional[int] = None,
        response_body: Optional[str] = None,
        request_summary: Optional[Mapping[str, Any]] = None,
        attempts: int = 1,
    ) -> None:
        super().__init__(message)
        self.message: str = message
        self.error_code: Optional[str] = error_code
        self.http_status: Optional[int] = http_status
        self.request_id: Optional[str] = request_id
        self.error_id: Optional[str] = error_id
        self.retry_after: Optional[int] = retry_after
        self.response_body: Optional[str] = response_body
        self.request_summary: Optional[Mapping[str, Any]] = request_summary
        self.attempts: int = attempts

    def __str__(self) -> str:
        parts = [self.message]
        if self.error_code:
            parts.append(f"error_code={self.error_code}")
        if self.http_status is not None:
            parts.append(f"http_status={self.http_status}")
        if self.request_id:
            parts.append(f"request_id={self.request_id}")
        if self.attempts > 1:
            parts.append(f"attempts={self.attempts}")
        return " ".join(parts)


class ConfigError(TesoteError):
    """Bad SDK configuration. Raised synchronously at construction."""


class EndpointRemovedError(TesoteError):
    """Method exists in this SDK version but its upstream endpoint is gone."""


class TransportError(TesoteError):
    """No usable HTTP response was received."""


class NetworkError(TransportError):
    """DNS failure, connection refused, connection reset, etc."""


class TimeoutError(TransportError):  # noqa: A001 -- intentional shadow of builtins
    """Connect or read timeout."""


class TlsError(TransportError):
    """Certificate or TLS handshake failure."""


class ApiError(TesoteError):
    """Server returned a structured error envelope."""


# 401 family -----------------------------------------------------------------


class UnauthorizedError(ApiError):
    """401 ``UNAUTHORIZED``."""


class ApiKeyRevokedError(ApiError):
    """401 ``API_KEY_REVOKED``."""


# 403 family -----------------------------------------------------------------


class WorkspaceSuspendedError(ApiError):
    """403 ``WORKSPACE_SUSPENDED``."""


class AccountDisabledError(ApiError):
    """403 ``ACCOUNT_DISABLED``."""


class HistorySyncForbiddenError(ApiError):
    """403 ``HISTORY_SYNC_FORBIDDEN``."""


# 404 family -----------------------------------------------------------------


class NotFoundError(ApiError):
    """Generic 404 fallback when a more specific subclass does not match."""


class AccountNotFoundError(NotFoundError):
    """404 ``ACCOUNT_NOT_FOUND``."""


class TransactionNotFoundError(NotFoundError):
    """404 ``TRANSACTION_NOT_FOUND``."""


class SyncSessionNotFoundError(NotFoundError):
    """404 ``SYNC_SESSION_NOT_FOUND``."""


class PaymentMethodNotFoundError(NotFoundError):
    """404 ``PAYMENT_METHOD_NOT_FOUND``."""


class TransactionOrderNotFoundError(NotFoundError):
    """404 ``TRANSACTION_ORDER_NOT_FOUND``."""


class BatchNotFoundError(NotFoundError):
    """404 ``BATCH_NOT_FOUND``."""


class BankConnectionNotFoundError(NotFoundError):
    """404 ``BANK_CONNECTION_NOT_FOUND``."""


# 409 family -----------------------------------------------------------------


class MutationDuringPaginationError(ApiError):
    """409 ``MUTATION_CONFLICT`` while iterating a cursor."""


class SyncInProgressError(ApiError):
    """409 ``SYNC_IN_PROGRESS``."""


class InvalidOrderStateError(ApiError):
    """409 ``INVALID_ORDER_STATE``."""


# 400/422 validation family --------------------------------------------------


class ValidationError(ApiError):
    """400 ``VALIDATION_ERROR``."""


class BatchValidationError(ApiError):
    """400 ``BATCH_VALIDATION_ERROR``."""


class UnprocessableContentError(ApiError):
    """422 ``UNPROCESSABLE_CONTENT`` -- generic validation failure."""


class InvalidDateRangeError(ApiError):
    """422 ``INVALID_DATE_RANGE``."""


class MissingDateRangeError(ApiError):
    """422 ``MISSING_DATE_RANGE``."""


class InvalidCursorError(ApiError):
    """422 ``INVALID_CURSOR``."""


class InvalidCountError(ApiError):
    """422 ``INVALID_COUNT``."""


class InvalidLimitError(ApiError):
    """422 ``INVALID_LIMIT``."""


class InvalidQueryError(ApiError):
    """422 ``INVALID_QUERY``."""


class BankSubmissionError(ApiError):
    """422 ``BANK_SUBMISSION_ERROR``."""


# 429 family -----------------------------------------------------------------


class RateLimitExceededError(ApiError):
    """429 ``RATE_LIMIT_EXCEEDED``."""


class SyncRateLimitExceededError(ApiError):
    """429 ``SYNC_RATE_LIMIT_EXCEEDED``."""


# 5xx family -----------------------------------------------------------------


class ServiceUnavailableError(ApiError):
    """503 -- platform pause mode."""


class BankUnderMaintenanceError(ApiError):
    """503 ``BANK_UNDER_MAINTENANCE``."""


class InternalError(ApiError):
    """500 ``INTERNAL_ERROR``."""


# why: error_code dispatch table; keep in sync with docs/architecture/errors.md
ERROR_CODE_TO_CLASS: Dict[str, Type[ApiError]] = {
    "UNAUTHORIZED": UnauthorizedError,
    "API_KEY_REVOKED": ApiKeyRevokedError,
    "WORKSPACE_SUSPENDED": WorkspaceSuspendedError,
    "ACCOUNT_DISABLED": AccountDisabledError,
    "HISTORY_SYNC_FORBIDDEN": HistorySyncForbiddenError,
    "ACCOUNT_NOT_FOUND": AccountNotFoundError,
    "TRANSACTION_NOT_FOUND": TransactionNotFoundError,
    "SYNC_SESSION_NOT_FOUND": SyncSessionNotFoundError,
    "PAYMENT_METHOD_NOT_FOUND": PaymentMethodNotFoundError,
    "TRANSACTION_ORDER_NOT_FOUND": TransactionOrderNotFoundError,
    "BATCH_NOT_FOUND": BatchNotFoundError,
    "BANK_CONNECTION_NOT_FOUND": BankConnectionNotFoundError,
    "MUTATION_CONFLICT": MutationDuringPaginationError,
    "SYNC_IN_PROGRESS": SyncInProgressError,
    "INVALID_ORDER_STATE": InvalidOrderStateError,
    "VALIDATION_ERROR": ValidationError,
    "BATCH_VALIDATION_ERROR": BatchValidationError,
    "UNPROCESSABLE_CONTENT": UnprocessableContentError,
    "INVALID_DATE_RANGE": InvalidDateRangeError,
    "MISSING_DATE_RANGE": MissingDateRangeError,
    "INVALID_CURSOR": InvalidCursorError,
    "INVALID_COUNT": InvalidCountError,
    "INVALID_LIMIT": InvalidLimitError,
    "INVALID_QUERY": InvalidQueryError,
    "BANK_SUBMISSION_ERROR": BankSubmissionError,
    "RATE_LIMIT_EXCEEDED": RateLimitExceededError,
    "SYNC_RATE_LIMIT_EXCEEDED": SyncRateLimitExceededError,
    "BANK_UNDER_MAINTENANCE": BankUnderMaintenanceError,
    "INTERNAL_ERROR": InternalError,
}


# why: HTTP-status fallback when the server omits or sends an unknown error_code
HTTP_STATUS_TO_CLASS: Dict[int, Type[ApiError]] = {
    401: UnauthorizedError,
    403: WorkspaceSuspendedError,
    404: NotFoundError,
    409: MutationDuringPaginationError,
    422: UnprocessableContentError,
    429: RateLimitExceededError,
    500: InternalError,
    503: ServiceUnavailableError,
}


def classify_api_error(
    error_code: Optional[str],
    http_status: Optional[int],
) -> Type[ApiError]:
    """Pick the narrowest ``ApiError`` subclass for a server response."""
    if error_code and error_code in ERROR_CODE_TO_CLASS:
        return ERROR_CODE_TO_CLASS[error_code]
    if http_status is not None and http_status in HTTP_STATUS_TO_CLASS:
        return HTTP_STATUS_TO_CLASS[http_status]
    return ApiError


__all__ = [
    "TesoteError",
    "ConfigError",
    "EndpointRemovedError",
    "TransportError",
    "NetworkError",
    "TimeoutError",
    "TlsError",
    "ApiError",
    # 401
    "UnauthorizedError",
    "ApiKeyRevokedError",
    # 403
    "WorkspaceSuspendedError",
    "AccountDisabledError",
    "HistorySyncForbiddenError",
    # 404
    "NotFoundError",
    "AccountNotFoundError",
    "TransactionNotFoundError",
    "SyncSessionNotFoundError",
    "PaymentMethodNotFoundError",
    "TransactionOrderNotFoundError",
    "BatchNotFoundError",
    "BankConnectionNotFoundError",
    # 409
    "MutationDuringPaginationError",
    "SyncInProgressError",
    "InvalidOrderStateError",
    # 400/422 validation
    "ValidationError",
    "BatchValidationError",
    "UnprocessableContentError",
    "InvalidDateRangeError",
    "MissingDateRangeError",
    "InvalidCursorError",
    "InvalidCountError",
    "InvalidLimitError",
    "InvalidQueryError",
    "BankSubmissionError",
    # 429
    "RateLimitExceededError",
    "SyncRateLimitExceededError",
    # 5xx
    "ServiceUnavailableError",
    "BankUnderMaintenanceError",
    "InternalError",
    # tables
    "ERROR_CODE_TO_CLASS",
    "HTTP_STATUS_TO_CLASS",
    "classify_api_error",
]
