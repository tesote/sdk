"""Single HTTP client for every resource module across v1/v2.

Owns: bearer injection, retries with exponential backoff + jitter, rate-limit
header capture, idempotency-key auto-generation, request-id propagation into
errors, opt-in TTL LRU cache, bearer redaction in logs and error summaries.

Stdlib only -- ``urllib.request`` + ``json``. Zero runtime dependencies.
"""

from __future__ import annotations

import json
import socket
import ssl
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
import uuid
from typing import Any, Callable, Dict, Mapping, Optional

from ._internals import (
    CacheBackend,
    InMemoryLRUCache,
    JsonType,
    LoggerCallback,
    RateLimitSnapshot,
    Response,
    RetryPolicy,
    redact_bearer,
)
from ._version import __version__
from .errors import (
    ApiError,
    NetworkError,
    TesoteError,
    TimeoutError,
    TlsError,
    TransportError,
    classify_api_error,
)

DEFAULT_BASE_URL = "https://equipo.tesote.com/api"
DEFAULT_CONNECT_TIMEOUT = 5.0
DEFAULT_READ_TIMEOUT = 30.0

_MUTATING_METHODS = frozenset({"POST", "PUT", "PATCH", "DELETE"})
_RETRY_STATUS = frozenset({429, 502, 503, 504})


def _sleep(seconds: float) -> None:  # why: indirection makes tests deterministic
    if seconds > 0:
        time.sleep(seconds)


# why: kept for back-compat with tests that import this private helper directly
_redact_bearer = redact_bearer


class Transport:
    """Single HTTP client used by every versioned resource module."""

    def __init__(
        self,
        api_key: str,
        *,
        base_url: str = DEFAULT_BASE_URL,
        connect_timeout: float = DEFAULT_CONNECT_TIMEOUT,
        read_timeout: float = DEFAULT_READ_TIMEOUT,
        retry_policy: Optional[RetryPolicy] = None,
        cache_backend: Optional[CacheBackend] = None,
        user_agent: Optional[str] = None,
        logger: Optional[LoggerCallback] = None,
        opener: Optional[Callable[..., Any]] = None,
    ) -> None:
        self._api_key = api_key
        self._base_url = base_url.rstrip("/")
        self._connect_timeout = connect_timeout
        self._read_timeout = read_timeout
        self._retry_policy = retry_policy or RetryPolicy()
        self._cache_backend = cache_backend
        self._user_agent = user_agent or self._default_user_agent()
        self._logger = logger
        # why: caller-injected opener allows tests to swap urlopen without monkeypatching
        self._opener = opener or urllib.request.urlopen
        self.last_rate_limit: Optional[RateLimitSnapshot] = None

    @staticmethod
    def _default_user_agent() -> str:
        v = sys.version_info
        return f"tesote-sdk-py/{__version__} (python/{v.major}.{v.minor}.{v.micro})"

    @property
    def base_url(self) -> str:
        return self._base_url

    def request(
        self,
        method: str,
        path: str,
        *,
        query: Optional[Mapping[str, Any]] = None,
        body: Optional[Mapping[str, Any]] = None,
        idempotency_key: Optional[str] = None,
        cache_ttl: Optional[float] = None,
        extra_headers: Optional[Mapping[str, str]] = None,
    ) -> Response:
        """Execute one logical request, including retries.

        Returns a parsed :class:`Response`. On a non-2xx that we do not retry
        (or retries exhausted), raises the typed :class:`ApiError` subclass.
        """
        method_upper = method.upper()
        full_url = self._build_url(path, query)
        headers = self._build_headers(method_upper, idempotency_key, extra_headers)

        cache_key = self._cache_key(method_upper, full_url) if cache_ttl else None
        if cache_key and self._cache_backend is not None:
            cached = self._cache_backend.get(cache_key)
            if cached is not None:
                _, response = cached
                return response

        encoded_body = json.dumps(body).encode("utf-8") if body is not None else None
        request_summary = self._summarize_request(method_upper, path, query, body)

        attempt = 0
        last_error: Optional[TesoteError] = None
        while attempt < self._retry_policy.max_attempts:
            attempt += 1
            try:
                response = self._perform(method_upper, full_url, headers, encoded_body)
            except TesoteError as exc:
                exc.attempts = attempt
                last_error = exc
                if not self._should_retry_error(exc, method_upper, idempotency_key):
                    raise
                _sleep(self._retry_policy.compute_delay(attempt))
                continue

            self._capture_rate_limit(response.headers)
            if 200 <= response.status < 300:
                if cache_key and self._cache_backend is not None and cache_ttl:
                    expires_at = time.monotonic() + cache_ttl
                    self._cache_backend.set(cache_key, (expires_at, response))
                return response

            api_error = self._build_api_error(
                response, request_summary=request_summary, attempts=attempt
            )
            last_error = api_error
            if response.status in _RETRY_STATUS and attempt < self._retry_policy.max_attempts:
                retry_after = self._parse_retry_after(response.headers)
                _sleep(self._retry_policy.compute_delay(attempt, retry_after))
                continue
            raise api_error

        # why: loop exited via `continue` only; surface the last typed error
        assert last_error is not None
        raise last_error

    # ------------------------------------------------------------------
    # internals
    # ------------------------------------------------------------------

    def _build_url(self, path: str, query: Optional[Mapping[str, Any]]) -> str:
        if not path.startswith("/"):
            path = "/" + path
        url = self._base_url + path
        if query:
            filtered = {k: v for k, v in query.items() if v is not None}
            if filtered:
                url = url + "?" + urllib.parse.urlencode(filtered, doseq=True)
        return url

    def _build_headers(
        self,
        method: str,
        idempotency_key: Optional[str],
        extra: Optional[Mapping[str, str]],
    ) -> Dict[str, str]:
        headers: Dict[str, str] = {
            "Authorization": f"Bearer {self._api_key}",
            "Accept": "application/json",
            "User-Agent": self._user_agent,
        }
        if method in _MUTATING_METHODS:
            headers["Content-Type"] = "application/json"
            key = idempotency_key or str(uuid.uuid4())
            headers["Idempotency-Key"] = key
        if extra:
            headers.update(extra)
        return headers

    @staticmethod
    def _cache_key(method: str, full_url: str) -> str:
        return f"{method} {full_url}"

    def _summarize_request(
        self,
        method: str,
        path: str,
        query: Optional[Mapping[str, Any]],
        body: Optional[Mapping[str, Any]],
    ) -> Dict[str, Any]:
        body_shape: Optional[Dict[str, Any]] = None
        if body is not None:
            body_shape = {"keys": sorted(body.keys()), "size": len(json.dumps(body))}
        return {
            "method": method,
            "path": path,
            "query": dict(query) if query else None,
            "body_shape": body_shape,
            # why: explicit redaction so a logger that captures the summary can't leak the key
            "authorization": redact_bearer(self._api_key),
        }

    def _perform(
        self,
        method: str,
        url: str,
        headers: Mapping[str, str],
        body: Optional[bytes],
    ) -> Response:
        req = urllib.request.Request(url=url, data=body, method=method)
        for k, v in headers.items():
            req.add_header(k, v)

        log_payload: Dict[str, Any] = {
            "phase": "request",
            "method": method,
            "url": url,
            "headers": {**headers, "Authorization": redact_bearer(self._api_key)},
        }
        self._log(log_payload)

        # why: urlopen `timeout` covers connect+read; the SDK exposes both for
        # symmetry with other languages -- we apply the larger of the two.
        timeout = max(self._connect_timeout, self._read_timeout)
        try:
            raw = self._opener(req, timeout=timeout)
        except urllib.error.HTTPError as exc:
            return self._read_http_error(exc)
        except urllib.error.URLError as exc:
            raise self._classify_url_error(exc) from exc
        except socket.timeout as exc:
            raise TimeoutError(
                f"Request timed out after {timeout}s",
            ) from exc
        except (ConnectionError, OSError) as exc:
            raise NetworkError(f"Network error: {exc}") from exc

        return self._read_http_response(raw)

    def _read_http_response(self, raw: Any) -> Response:
        status = int(getattr(raw, "status", None) or raw.getcode())
        headers = self._headers_to_dict(raw.headers)
        body_bytes: bytes = raw.read()
        try:
            raw.close()
        except Exception:  # noqa: BLE001 -- best-effort cleanup
            pass
        body_text = body_bytes.decode("utf-8", errors="replace")
        parsed = self._parse_json(body_text)
        request_id = headers.get("X-Request-Id") or headers.get("x-request-id")
        self._log(
            {
                "phase": "response",
                "status": status,
                "request_id": request_id,
                "headers": headers,
            }
        )
        return Response(status, headers, body_text, parsed, request_id)

    def _read_http_error(self, exc: urllib.error.HTTPError) -> Response:
        status = exc.code
        headers = self._headers_to_dict(exc.headers) if exc.headers else {}
        try:
            body_bytes = exc.read()
        except Exception:  # noqa: BLE001
            body_bytes = b""
        body_text = body_bytes.decode("utf-8", errors="replace")
        parsed = self._parse_json(body_text)
        request_id = headers.get("X-Request-Id") or headers.get("x-request-id")
        self._log(
            {
                "phase": "response",
                "status": status,
                "request_id": request_id,
                "headers": headers,
            }
        )
        return Response(status, headers, body_text, parsed, request_id)

    @staticmethod
    def _headers_to_dict(headers: Any) -> Dict[str, str]:
        if hasattr(headers, "items"):
            return {str(k): str(v) for k, v in headers.items()}
        return {}

    @staticmethod
    def _parse_json(text: str) -> JsonType:
        if not text:
            return None
        try:
            return json.loads(text)  # type: ignore[no-any-return]
        except json.JSONDecodeError:
            return None

    def _classify_url_error(self, exc: urllib.error.URLError) -> TransportError:
        reason = exc.reason
        if isinstance(reason, ssl.SSLError):
            return TlsError(f"TLS error: {reason}")
        if isinstance(reason, socket.timeout):
            return TimeoutError(f"Connection timed out: {reason}")
        return NetworkError(f"Network error: {reason}")

    def _capture_rate_limit(self, headers: Mapping[str, str]) -> None:
        def _int(name: str) -> Optional[int]:
            value = headers.get(name) or headers.get(name.lower())
            if value is None:
                return None
            try:
                return int(value)
            except ValueError:
                return None

        self.last_rate_limit = RateLimitSnapshot(
            limit=_int("X-RateLimit-Limit"),
            remaining=_int("X-RateLimit-Remaining"),
            reset=_int("X-RateLimit-Reset"),
        )

    @staticmethod
    def _parse_retry_after(headers: Mapping[str, str]) -> Optional[float]:
        value = headers.get("Retry-After") or headers.get("retry-after")
        if value is None:
            return None
        try:
            return float(value)
        except ValueError:
            return None

    @staticmethod
    def _should_retry_error(
        exc: TesoteError, method: str, idempotency_key: Optional[str]
    ) -> bool:
        if isinstance(exc, TimeoutError):
            # why: read-timeouts on non-idempotent calls may have succeeded server-side
            if method in _MUTATING_METHODS and idempotency_key is None:
                return False
            return True
        if isinstance(exc, NetworkError):
            return True
        return False

    def _build_api_error(
        self,
        response: Response,
        *,
        request_summary: Mapping[str, Any],
        attempts: int,
    ) -> ApiError:
        envelope: Dict[str, Any] = {}
        if isinstance(response.json, dict):
            envelope = response.json
        error_code = envelope.get("error_code")
        error_id = envelope.get("error_id")
        message = envelope.get("error") or f"HTTP {response.status}"
        retry_after = self._parse_retry_after(response.headers)
        if retry_after is None and isinstance(envelope.get("retry_after"), (int, float)):
            retry_after = float(envelope["retry_after"])

        cls = classify_api_error(error_code, response.status)
        return cls(
            message,
            error_code=error_code,
            http_status=response.status,
            request_id=response.request_id,
            error_id=error_id,
            retry_after=int(retry_after) if retry_after is not None else None,
            response_body=response.body,
            request_summary=request_summary,
            attempts=attempts,
        )

    def _log(self, payload: Mapping[str, Any]) -> None:
        if self._logger is None:
            return
        try:
            self._logger(payload)
        except Exception:  # noqa: BLE001 -- a misbehaving logger must not break a request
            pass


__all__ = [
    "Transport",
    "Response",
    "RetryPolicy",
    "RateLimitSnapshot",
    "CacheBackend",
    "InMemoryLRUCache",
    "LoggerCallback",
    "DEFAULT_BASE_URL",
    "DEFAULT_CONNECT_TIMEOUT",
    "DEFAULT_READ_TIMEOUT",
]
