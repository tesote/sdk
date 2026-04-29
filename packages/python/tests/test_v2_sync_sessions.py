"""v2 sync_sessions resource: list, iter, get."""

from __future__ import annotations

import pytest

from tesote_sdk.errors import (
    BankConnectionNotFoundError,
    SyncSessionNotFoundError,
)
from tesote_sdk.v2.sync_sessions import SyncSessionsResource

from ._helpers import make_transport
from .conftest import http_error, ok_response


def _session(sid: str, status: str = "completed") -> dict:
    return {
        "id": sid,
        "status": status,
        "started_at": "2026-04-01T00:00:00Z",
        "completed_at": "2026-04-01T00:00:30Z",
        "transactions_synced": 10,
        "accounts_count": 1,
    }


def test_list_returns_paginated_sessions() -> None:
    transport, opener = make_transport(
        [
            ok_response(
                {
                    "sync_sessions": [_session("ss_1"), _session("ss_2", "failed")],
                    "limit": 50,
                    "offset": 0,
                    "has_more": False,
                }
            )
        ]
    )
    result = SyncSessionsResource(transport).list("acct_1", status="completed")
    assert len(result.sync_sessions) == 2
    assert result.has_more is False
    assert "status=completed" in opener.calls[0]["url"]


def test_iter_follows_offset_pagination() -> None:
    p1 = ok_response(
        {
            "sync_sessions": [_session("ss_1"), _session("ss_2")],
            "limit": 2,
            "offset": 0,
            "has_more": True,
        }
    )
    p2 = ok_response(
        {
            "sync_sessions": [_session("ss_3")],
            "limit": 2,
            "offset": 2,
            "has_more": False,
        }
    )
    transport, opener = make_transport([p1, p2])
    sessions = list(
        SyncSessionsResource(transport).iter("acct_1", page_size=2)
    )
    assert [s.id for s in sessions] == ["ss_1", "ss_2", "ss_3"]
    assert "offset=2" in opener.calls[1]["url"]


def test_get_returns_typed_session_with_error_block_when_failed() -> None:
    transport, _ = make_transport(
        [
            ok_response(
                {
                    **_session("ss_1", "failed"),
                    "error": {"type": "BankError", "message": "down"},
                    "performance": {
                        "total_duration": 1.5,
                        "complexity_score": 0.8,
                        "sync_speed_score": 0.9,
                    },
                }
            )
        ]
    )
    session = SyncSessionsResource(transport).get("acct_1", "ss_1")
    assert session.error is not None
    assert session.error.message == "down"
    assert session.performance is not None
    assert session.performance.total_duration == 1.5


def test_get_404_session_not_found() -> None:
    transport, _ = make_transport(
        [
            http_error(
                404,
                {},
                {"error": "missing", "error_code": "SYNC_SESSION_NOT_FOUND"},
            )
        ]
    )
    with pytest.raises(SyncSessionNotFoundError):
        SyncSessionsResource(transport).get("acct_1", "ss_missing")


def test_list_404_bank_connection_not_found() -> None:
    transport, _ = make_transport(
        [
            http_error(
                404,
                {},
                {"error": "no bank link", "error_code": "BANK_CONNECTION_NOT_FOUND"},
            )
        ]
    )
    with pytest.raises(BankConnectionNotFoundError):
        SyncSessionsResource(transport).list("acct_1")
