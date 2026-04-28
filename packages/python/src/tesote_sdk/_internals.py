"""Transport helper types: Response, RetryPolicy, CacheBackend, etc.

Split out so ``transport.py`` stays focused on the request loop.
"""

from __future__ import annotations

import random
import time
from collections import OrderedDict
from typing import (
    Any,
    Callable,
    Dict,
    List,
    Mapping,
    Optional,
    Protocol,
    Tuple,
    Union,
    runtime_checkable,
)

JsonType = Union[Dict[str, Any], List[Any], str, int, float, bool, None]
LoggerCallback = Callable[[Mapping[str, Any]], None]


class Response:
    """Parsed transport response."""

    __slots__ = ("status", "headers", "body", "json", "request_id")

    def __init__(
        self,
        status: int,
        headers: Mapping[str, str],
        body: str,
        parsed: JsonType,
        request_id: Optional[str],
    ) -> None:
        self.status = status
        self.headers = headers
        self.body = body
        self.json = parsed
        self.request_id = request_id


CacheEntry = Tuple[float, Response]


@runtime_checkable
class CacheBackend(Protocol):
    """Pluggable response-cache contract.

    Default implementation is :class:`InMemoryLRUCache`. Users can drop in
    Redis, memcached, etc. by implementing ``get`` / ``set``.
    """

    def get(self, key: str) -> Optional[CacheEntry]:
        ...

    def set(self, key: str, value: CacheEntry) -> None:
        ...


class InMemoryLRUCache:
    """Bounded LRU with per-entry TTL stored alongside the value."""

    def __init__(self, max_entries: int = 256) -> None:
        self._max = max_entries
        self._store: OrderedDict[str, CacheEntry] = OrderedDict()

    def get(self, key: str) -> Optional[CacheEntry]:
        value = self._store.get(key)
        if value is None:
            return None
        expires_at, _response = value
        if expires_at <= time.monotonic():
            # why: drop stale entries on access; cheaper than a sweep thread
            self._store.pop(key, None)
            return None
        self._store.move_to_end(key)
        return value

    def set(self, key: str, value: CacheEntry) -> None:
        self._store[key] = value
        self._store.move_to_end(key)
        while len(self._store) > self._max:
            self._store.popitem(last=False)


class RetryPolicy:
    """Configurable retry policy. Defaults match docs/architecture/transport.md."""

    def __init__(
        self,
        max_attempts: int = 3,
        base_delay: float = 0.25,
        max_delay: float = 8.0,
    ) -> None:
        if max_attempts < 1:
            raise ValueError("max_attempts must be >= 1")
        self.max_attempts = max_attempts
        self.base_delay = base_delay
        self.max_delay = max_delay

    def compute_delay(self, attempt: int, retry_after: Optional[float] = None) -> float:
        if retry_after is not None and retry_after > 0:
            return min(retry_after, self.max_delay)
        capped = min(self.max_delay, self.base_delay * (2 ** attempt))
        # why: full jitter [0, capped] -- spreads retries when many clients fail together
        return random.uniform(0.0, capped)


class RateLimitSnapshot:
    """Read view of the most-recent rate-limit headers."""

    __slots__ = ("limit", "remaining", "reset")

    def __init__(
        self,
        limit: Optional[int],
        remaining: Optional[int],
        reset: Optional[int],
    ) -> None:
        self.limit = limit
        self.remaining = remaining
        self.reset = reset


def redact_bearer(api_key: str) -> str:
    if len(api_key) <= 4:
        return "Bearer ****"
    return "Bearer " + ("*" * (len(api_key) - 4)) + api_key[-4:]


__all__ = [
    "JsonType",
    "LoggerCallback",
    "Response",
    "CacheEntry",
    "CacheBackend",
    "InMemoryLRUCache",
    "RetryPolicy",
    "RateLimitSnapshot",
    "redact_bearer",
]
