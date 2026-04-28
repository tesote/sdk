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


class UnauthorizedError(ApiError):
    """401 ``UNAUTHORIZED``."""


class ApiKeyRevokedError(ApiError):
    """401 ``API_KEY_REVOKED``."""


class WorkspaceSuspendedError(ApiError):
    """403 ``WORKSPACE_SUSPENDED``."""


class AccountDisabledError(ApiError):
    """403 ``ACCOUNT_DISABLED``."""


class HistorySyncForbiddenError(ApiError):
    """403 ``HISTORY_SYNC_FORBIDDEN``."""


class MutationDuringPaginationError(ApiError):
    """409 ``MUTATION_CONFLICT`` while iterating a cursor."""


class UnprocessableContentError(ApiError):
    """422 ``UNPROCESSABLE_CONTENT`` -- generic validation failure."""


class InvalidDateRangeError(ApiError):
    """422 ``INVALID_DATE_RANGE``."""


class RateLimitExceededError(ApiError):
    """429 ``RATE_LIMIT_EXCEEDED``."""


class ServiceUnavailableError(ApiError):
    """503 -- platform pause mode."""


# why: error_code dispatch table; keep in sync with docs/architecture/errors.md
ERROR_CODE_TO_CLASS: Dict[str, Type[ApiError]] = {
    "UNAUTHORIZED": UnauthorizedError,
    "API_KEY_REVOKED": ApiKeyRevokedError,
    "WORKSPACE_SUSPENDED": WorkspaceSuspendedError,
    "ACCOUNT_DISABLED": AccountDisabledError,
    "HISTORY_SYNC_FORBIDDEN": HistorySyncForbiddenError,
    "MUTATION_CONFLICT": MutationDuringPaginationError,
    "UNPROCESSABLE_CONTENT": UnprocessableContentError,
    "INVALID_DATE_RANGE": InvalidDateRangeError,
    "RATE_LIMIT_EXCEEDED": RateLimitExceededError,
}


# why: HTTP-status fallback when the server omits or sends an unknown error_code
HTTP_STATUS_TO_CLASS: Dict[int, Type[ApiError]] = {
    401: UnauthorizedError,
    403: WorkspaceSuspendedError,
    409: MutationDuringPaginationError,
    422: UnprocessableContentError,
    429: RateLimitExceededError,
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
    "UnauthorizedError",
    "ApiKeyRevokedError",
    "WorkspaceSuspendedError",
    "AccountDisabledError",
    "HistorySyncForbiddenError",
    "MutationDuringPaginationError",
    "UnprocessableContentError",
    "InvalidDateRangeError",
    "RateLimitExceededError",
    "ServiceUnavailableError",
    "ERROR_CODE_TO_CLASS",
    "HTTP_STATUS_TO_CLASS",
    "classify_api_error",
]
