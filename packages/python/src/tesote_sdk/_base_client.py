"""Shared client construction logic for v1/v2/v3 clients."""

from __future__ import annotations

from typing import Optional

from .errors import ConfigError
from .transport import (
    DEFAULT_BASE_URL,
    DEFAULT_CONNECT_TIMEOUT,
    DEFAULT_READ_TIMEOUT,
    CacheBackend,
    LoggerCallback,
    RetryPolicy,
    Transport,
)


def build_transport(
    api_key: str,
    *,
    base_url: str = DEFAULT_BASE_URL,
    connect_timeout: float = DEFAULT_CONNECT_TIMEOUT,
    read_timeout: float = DEFAULT_READ_TIMEOUT,
    retry_policy: Optional[RetryPolicy] = None,
    cache_backend: Optional[CacheBackend] = None,
    user_agent: Optional[str] = None,
    logger: Optional[LoggerCallback] = None,
) -> Transport:
    if not api_key:
        raise ConfigError("api_key is required")
    return Transport(
        api_key=api_key,
        base_url=base_url,
        connect_timeout=connect_timeout,
        read_timeout=read_timeout,
        retry_policy=retry_policy,
        cache_backend=cache_backend,
        user_agent=user_agent,
        logger=logger,
    )


__all__ = ["build_transport"]
