"""v2 accounts resource: list, get, sync."""

from __future__ import annotations

import pytest

from tesote_sdk.errors import (
    AccountNotFoundError,
    BankUnderMaintenanceError,
    SyncInProgressError,
    SyncRateLimitExceededError,
    UnprocessableContentError,
)
from tesote_sdk.v2.accounts import AccountsResource

from ._helpers import make_transport
from .conftest import http_error, ok_response


def test_list_returns_account_list() -> None:
    transport, _ = make_transport(
        [
            ok_response(
                {
                    "total": 0,
                    "accounts": [],
                    "pagination": {
                        "current_page": 1,
                        "per_page": 50,
                        "total_pages": 1,
                        "total_count": 0,
                    },
                }
            )
        ]
    )
    result = AccountsResource(transport).list(page=1)
    assert result.total == 0
    assert result.pagination is not None
    assert result.pagination.current_page == 1


def test_get_uses_v2_prefix() -> None:
    transport, opener = make_transport(
        [
            ok_response(
                {
                    "id": "acct_1",
                    "name": "Op",
                    "data": {},
                    "bank": {},
                    "legal_entity": {},
                }
            )
        ]
    )
    AccountsResource(transport).get("acct_1")
    assert opener.calls[0]["url"].endswith("/v2/accounts/acct_1")


def test_sync_returns_started_response_and_sets_idempotency_header() -> None:
    transport, opener = make_transport(
        [
            ok_response(
                {
                    "message": "Sync started",
                    "sync_session_id": "ss_1",
                    "status": "pending",
                    "started_at": "2026-04-28T19:21:00Z",
                },
                status=202,
            )
        ]
    )
    res = AccountsResource(transport).sync("acct_1", idempotency_key="my-key")
    assert res.sync_session_id == "ss_1"
    assert res.status == "pending"
    headers = dict(opener.calls[0]["headers"])
    assert headers["Idempotency-key"] == "my-key"
    # Content-Type required because POST body sent
    assert headers["Content-type"] == "application/json"


def test_sync_auto_generates_idempotency_key_when_omitted() -> None:
    transport, opener = make_transport(
        [
            ok_response(
                {
                    "message": "Sync started",
                    "sync_session_id": "ss_1",
                    "status": "pending",
                    "started_at": "2026-04-28T19:21:00Z",
                },
                status=202,
            )
        ]
    )
    AccountsResource(transport).sync("acct_1")
    headers = dict(opener.calls[0]["headers"])
    assert "Idempotency-key" in headers
    assert len(headers["Idempotency-key"]) >= 16


def test_sync_409_raises_sync_in_progress() -> None:
    transport, _ = make_transport(
        [
            http_error(
                409,
                {},
                {"error": "in flight", "error_code": "SYNC_IN_PROGRESS"},
            )
        ]
    )
    with pytest.raises(SyncInProgressError):
        AccountsResource(transport).sync("acct_1")


def test_sync_429_raises_sync_rate_limit_exceeded() -> None:
    transport, _ = make_transport(
        [
            http_error(
                429,
                {"Retry-After": "0"},
                {"error": "wait", "error_code": "SYNC_RATE_LIMIT_EXCEEDED"},
            ),
            http_error(
                429,
                {"Retry-After": "0"},
                {"error": "wait", "error_code": "SYNC_RATE_LIMIT_EXCEEDED"},
            ),
            http_error(
                429,
                {"Retry-After": "0"},
                {"error": "wait", "error_code": "SYNC_RATE_LIMIT_EXCEEDED"},
            ),
        ]
    )
    with pytest.raises(SyncRateLimitExceededError):
        AccountsResource(transport).sync("acct_1")


def test_sync_503_bank_under_maintenance() -> None:
    transport, _ = make_transport(
        [
            http_error(503, {}, {"error": "down", "error_code": "BANK_UNDER_MAINTENANCE"})
            for _ in range(3)
        ]
    )
    with pytest.raises(BankUnderMaintenanceError):
        AccountsResource(transport).sync("acct_1")


def test_sync_404_raises_account_not_found() -> None:
    transport, _ = make_transport(
        [http_error(404, {}, {"error": "no", "error_code": "ACCOUNT_NOT_FOUND"})]
    )
    with pytest.raises(AccountNotFoundError):
        AccountsResource(transport).sync("acct_missing")


def test_415_when_content_type_missing_maps_to_unprocessable_content() -> None:
    """Server returns 415 when POST/PATCH lack Content-Type. Transport always sets it,
    but we verify the typed-error mapping path here using a forced 415 response."""
    transport, _ = make_transport(
        [
            http_error(
                415,
                {},
                {"error": "need json", "error_code": "UNPROCESSABLE_CONTENT"},
            )
        ]
    )
    with pytest.raises(UnprocessableContentError):
        AccountsResource(transport).sync("acct_1")
