"""Public API for the Tesote Python SDK.

Re-exports the versioned clients and every error class.
"""

from ._version import __version__
from .errors import (
    AccountDisabledError,
    ApiError,
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
)
from .v1 import V1Client
from .v2 import V2Client

__all__ = [
    "__version__",
    "V1Client",
    "V2Client",
    # error hierarchy
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
]
