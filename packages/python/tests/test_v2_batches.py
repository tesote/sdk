"""v2 batches resource: create, get, approve, submit, cancel."""

from __future__ import annotations

import json

import pytest

from tesote_sdk.errors import (
    BatchNotFoundError,
    BatchValidationError,
    InvalidOrderStateError,
)
from tesote_sdk.v2.batches import BatchesResource

from ._helpers import make_transport
from .conftest import http_error, ok_response


def _order(oid: str = "o1", status: str = "draft") -> dict:
    return {
        "id": oid,
        "status": status,
        "amount": 50.0,
        "currency": "VES",
        "description": "x",
        "created_at": "2026-04-01T00:00:00Z",
        "updated_at": "2026-04-01T00:00:00Z",
    }


def test_create_wraps_orders_in_body_and_returns_batch_create_result() -> None:
    transport, opener = make_transport(
        [
            ok_response(
                {
                    "batch_id": "b_1",
                    "orders": [_order("o1"), _order("o2")],
                    "errors": [],
                },
                status=201,
            )
        ]
    )
    res = BatchesResource(transport).create(
        "acct_1",
        [
            {"amount": "10.00", "currency": "VES", "description": "a"},
            {"amount": "20.00", "currency": "VES", "description": "b"},
        ],
        idempotency_key="bk-1",
    )
    assert res.batch_id == "b_1"
    assert len(res.orders) == 2
    body = json.loads(opener.calls[0]["body"].decode("utf-8"))
    assert "orders" in body and len(body["orders"]) == 2
    headers = dict(opener.calls[0]["headers"])
    assert headers["Idempotency-key"] == "bk-1"


def test_create_batch_validation_error() -> None:
    transport, _ = make_transport(
        [
            http_error(
                400,
                {},
                {"error": "bad batch", "error_code": "BATCH_VALIDATION_ERROR"},
            )
        ]
    )
    with pytest.raises(BatchValidationError):
        BatchesResource(transport).create("acct_1", [])


def test_get_returns_summary_with_statuses() -> None:
    transport, _ = make_transport(
        [
            ok_response(
                {
                    "batch_id": "b_1",
                    "total_orders": 5,
                    "total_amount_cents": 50000,
                    "amount_currency": "VES",
                    "statuses": {"draft": 3, "approved": 2},
                    "batch_status": "mixed",
                    "created_at": "2026-04-01T00:00:00Z",
                    "orders": [_order("o1"), _order("o2", "approved")],
                }
            )
        ]
    )
    summary = BatchesResource(transport).get("acct_1", "b_1")
    assert summary.batch_id == "b_1"
    assert summary.statuses == {"draft": 3, "approved": 2}
    assert summary.batch_status == "mixed"
    assert len(summary.orders) == 2


def test_get_404_batch_not_found() -> None:
    transport, _ = make_transport(
        [http_error(404, {}, {"error": "missing", "error_code": "BATCH_NOT_FOUND"})]
    )
    with pytest.raises(BatchNotFoundError):
        BatchesResource(transport).get("acct_1", "missing")


def test_approve_returns_counts() -> None:
    transport, opener = make_transport([ok_response({"approved": 5, "failed": 0})])
    res = BatchesResource(transport).approve("acct_1", "b_1")
    assert res.approved == 5
    assert res.failed == 0
    assert opener.calls[0]["url"].endswith("/v2/accounts/acct_1/batches/b_1/approve")


def test_approve_invalid_order_state() -> None:
    transport, _ = make_transport(
        [http_error(409, {}, {"error": "bad state", "error_code": "INVALID_ORDER_STATE"})]
    )
    with pytest.raises(InvalidOrderStateError):
        BatchesResource(transport).approve("acct_1", "b_1")


def test_submit_with_token() -> None:
    transport, opener = make_transport([ok_response({"enqueued": 5, "failed": 0})])
    res = BatchesResource(transport).submit("acct_1", "b_1", token="otp")
    assert res.enqueued == 5
    body = json.loads(opener.calls[0]["body"].decode("utf-8"))
    assert body == {"token": "otp"}


def test_cancel_returns_counts_and_skipped() -> None:
    transport, opener = make_transport(
        [ok_response({"cancelled": 4, "skipped": 1, "errors": []})]
    )
    res = BatchesResource(transport).cancel("acct_1", "b_1", idempotency_key="bk-cancel")
    assert res.cancelled == 4
    assert res.skipped == 1
    headers = dict(opener.calls[0]["headers"])
    assert headers["Idempotency-key"] == "bk-cancel"
