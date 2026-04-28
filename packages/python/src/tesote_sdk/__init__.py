"""Public API for the Tesote Python SDK.

Re-exports the three versioned clients, every error class, and the v3
``verify_webhook_signature`` helper.
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
from .v3 import V3Client, verify_webhook_signature

__all__ = [
    "__version__",
    "V1Client",
    "V2Client",
    "V3Client",
    "verify_webhook_signature",
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
