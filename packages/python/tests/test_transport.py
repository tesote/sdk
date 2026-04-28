"""Transport behavior matches docs/architecture/transport.md.

Mocks `urllib.request.urlopen` via the injected `opener` callback so each test
fully controls request/response cycles without touching the network.
"""

from __future__ import annotations

import time
from typing import Any, Dict, List

import pytest

from tesote_sdk import RateLimitExceededError, UnauthorizedError
from tesote_sdk.errors import (
    ApiError,
    NetworkError,
    UnprocessableContentError,
)
from tesote_sdk.errors import (
    TimeoutError as SdkTimeoutError,
)
from tesote_sdk.transport import (
    DEFAULT_BASE_URL,
    InMemoryLRUCache,
    RetryPolicy,
    Transport,
)

from .conftest import ScriptedOpener, http_error, ok_response


def _make(opener: ScriptedOpener, **kwargs: Any) -> Transport:
    return Transport(
        api_key="sk_test_abcdef1234",
        retry_policy=RetryPolicy(max_attempts=3, base_delay=0.0, max_delay=0.0),
        opener=opener,
        **kwargs,
    )


# ---------- happy path ----------


def test_get_succeeds_and_parses_json() -> None:
    opener = ScriptedOpener([ok_response({"id": "acct_1"})])
    t = _make(opener)
    response = t.request("GET", "/v3/accounts/acct_1")
    assert response.status == 200
    assert response.json == {"id": "acct_1"}
    assert opener.calls[0]["url"] == f"{DEFAULT_BASE_URL}/v3/accounts/acct_1"
    headers = dict(opener.calls[0]["headers"])
    assert headers["Authorization"] == "Bearer sk_test_abcdef1234"
    assert headers["Accept"] == "application/json"
    assert "tesote-sdk-py/" in headers["User-agent"]


def test_query_string_is_url_encoded_and_drops_none() -> None:
    opener = ScriptedOpener([ok_response([])])
    t = _make(opener)
    t.request("GET", "/v3/accounts", query={"cursor": "abc", "limit": 10, "skip": None})
    url = opener.calls[0]["url"]
    assert "cursor=abc" in url
    assert "limit=10" in url
    assert "skip" not in url


# ---------- mutating methods ----------


def test_post_auto_generates_idempotency_key_when_missing() -> None:
    opener = ScriptedOpener([ok_response({"ok": True}, status=201)])
    t = _make(opener)
    t.request("POST", "/v3/accounts/acct_1/sync", body={})
    headers = dict(opener.calls[0]["headers"])
    assert "Idempotency-key" in headers and len(headers["Idempotency-key"]) >= 16
    assert headers["Content-type"] == "application/json"


def test_post_uses_caller_supplied_idempotency_key() -> None:
    opener = ScriptedOpener([ok_response({}, status=201)])
    t = _make(opener)
    t.request("POST", "/v3/accounts/acct_1/sync", body={"foo": "bar"}, idempotency_key="k-123")
    headers = dict(opener.calls[0]["headers"])
    assert headers["Idempotency-key"] == "k-123"


# ---------- rate-limit header capture ----------


def test_rate_limit_headers_are_captured_on_success() -> None:
    opener = ScriptedOpener(
        [
            ok_response(
                {"data": []},
                headers={
                    "X-RateLimit-Limit": "200",
                    "X-RateLimit-Remaining": "157",
                    "X-RateLimit-Reset": "60",
                },
            )
        ]
    )
    t = _make(opener)
    t.request("GET", "/v3/accounts")
    snap = t.last_rate_limit
    assert snap is not None
    assert snap.limit == 200
    assert snap.remaining == 157
    assert snap.reset == 60


# ---------- request id propagation ----------


def test_request_id_propagates_into_thrown_error() -> None:
    opener = ScriptedOpener(
        [
            http_error(
                401,
                {"X-Request-Id": "rid-001"},
                {"error": "bad key", "error_code": "UNAUTHORIZED"},
            )
        ]
    )
    t = _make(opener)
    with pytest.raises(UnauthorizedError) as exc:
        t.request("GET", "/v3/accounts")
    err = exc.value
    assert err.request_id == "rid-001"
    assert err.error_code == "UNAUTHORIZED"
    assert err.http_status == 401
    assert err.attempts == 1


def test_422_maps_to_unprocessable_content() -> None:
    opener = ScriptedOpener(
        [
            http_error(
                422,
                {"X-Request-Id": "rid-422"},
                {"error": "bad", "error_code": "UNPROCESSABLE_CONTENT", "error_id": "eid-1"},
            )
        ]
    )
    t = _make(opener)
    with pytest.raises(UnprocessableContentError) as exc:
        t.request("POST", "/v3/transactions/bulk", body={"items": []})
    assert exc.value.error_id == "eid-1"
    assert exc.value.attempts == 1


# ---------- retries ----------


def test_429_retries_then_raises_rate_limit_exceeded() -> None:
    err = http_error(
        429,
        {"Retry-After": "0", "X-Request-Id": "rid-429"},
        {"error": "slow down", "error_code": "RATE_LIMIT_EXCEEDED", "retry_after": 0},
    )
    opener = ScriptedOpener([err, err, err])
    t = _make(opener)
    with pytest.raises(RateLimitExceededError) as exc:
        t.request("GET", "/v3/accounts")
    assert len(opener.calls) == 3
    assert exc.value.attempts == 3
    assert exc.value.retry_after == 0


def test_503_retries_until_success() -> None:
    opener = ScriptedOpener(
        [
            http_error(503, {}, {"error": "down", "error_code": "SERVICE_UNAVAILABLE"}),
            ok_response({"ok": True}),
        ]
    )
    t = _make(opener)
    response = t.request("GET", "/v3/accounts")
    assert response.status == 200
    assert len(opener.calls) == 2


def test_4xx_other_than_429_does_not_retry() -> None:
    opener = ScriptedOpener(
        [http_error(401, {}, {"error": "nope", "error_code": "UNAUTHORIZED"})]
    )
    t = _make(opener)
    with pytest.raises(UnauthorizedError):
        t.request("GET", "/v3/accounts")
    assert len(opener.calls) == 1


def test_network_error_retries_then_raises() -> None:
    opener = ScriptedOpener(
        [
            ConnectionResetError("reset"),
            ConnectionResetError("reset"),
            ok_response({"ok": True}),
        ]
    )
    t = _make(opener)
    response = t.request("GET", "/v3/accounts")
    assert response.status == 200
    assert len(opener.calls) == 3


def test_post_without_idempotency_does_not_retry_on_timeout() -> None:
    # why: rule from transport.md -- timeout on non-idempotent POST without an
    # idempotency key may have succeeded server-side, so we must surface, not retry.
    from tesote_sdk.transport import Transport as _T

    assert _T._should_retry_error(SdkTimeoutError("t"), "POST", None) is False
    assert _T._should_retry_error(SdkTimeoutError("t"), "POST", "key-1") is True
    assert _T._should_retry_error(SdkTimeoutError("t"), "GET", None) is True


def test_get_timeout_does_retry() -> None:
    import socket

    opener = ScriptedOpener([socket.timeout("read"), ok_response({"ok": True})])
    t = _make(opener)
    response = t.request("GET", "/v3/accounts")
    assert response.status == 200
    assert len(opener.calls) == 2


# ---------- bearer redaction ----------


def test_bearer_token_is_redacted_in_logs_and_request_summary() -> None:
    captured: List[Dict[str, Any]] = []

    def logger(payload: Dict[str, Any]) -> None:
        captured.append(payload)

    opener = ScriptedOpener(
        [http_error(401, {}, {"error": "nope", "error_code": "UNAUTHORIZED"})]
    )
    t = _make(opener, logger=logger)
    with pytest.raises(UnauthorizedError) as exc:
        t.request("POST", "/v3/accounts/acct_1/sync", body={"a": 1})

    # logger never sees the raw bearer in request-phase entries
    request_entries = [e for e in captured if e.get("phase") == "request"]
    assert request_entries, "expected at least one request log entry"
    for entry in request_entries:
        auth = entry["headers"].get("Authorization", "")
        assert "sk_test_abcdef1234" not in auth
        assert auth.endswith("1234")  # last4 visible

    # error.request_summary doesn't leak it either
    summary = exc.value.request_summary
    assert summary is not None
    assert "sk_test_abcdef1234" not in str(summary)
    assert summary["authorization"].endswith("1234")


# ---------- cache ----------


def test_ttl_cache_hits_skip_network() -> None:
    opener = ScriptedOpener([ok_response({"id": "acct_1"})])
    t = _make(opener, cache_backend=InMemoryLRUCache())
    a = t.request("GET", "/v3/accounts/acct_1", cache_ttl=30.0)
    b = t.request("GET", "/v3/accounts/acct_1", cache_ttl=30.0)
    assert a.json == {"id": "acct_1"}
    assert b.json == {"id": "acct_1"}
    assert len(opener.calls) == 1


def test_ttl_cache_expires() -> None:
    opener = ScriptedOpener([ok_response({"v": 1}), ok_response({"v": 2})])
    cache = InMemoryLRUCache()
    t = _make(opener, cache_backend=cache)
    t.request("GET", "/v3/accounts/acct_1", cache_ttl=0.01)
    time.sleep(0.02)
    second = t.request("GET", "/v3/accounts/acct_1", cache_ttl=0.01)
    assert second.json == {"v": 2}
    assert len(opener.calls) == 2


# ---------- attempts counter ----------


def test_attempts_set_on_eventual_success_error_chain() -> None:
    opener = ScriptedOpener(
        [
            http_error(503, {}, {"error": "down", "error_code": "X"}),
            http_error(503, {}, {"error": "down", "error_code": "X"}),
            http_error(503, {}, {"error": "down", "error_code": "X"}),
        ]
    )
    t = _make(opener)
    with pytest.raises(ApiError) as exc:
        t.request("GET", "/v3/accounts")
    assert exc.value.attempts == 3


# ---------- transport classifies low-level errors ----------


def test_classify_connection_reset_as_network_error() -> None:
    opener = ScriptedOpener([ConnectionResetError("boom")] * 3)
    t = _make(opener)
    with pytest.raises(NetworkError):
        t.request("GET", "/v3/accounts")
