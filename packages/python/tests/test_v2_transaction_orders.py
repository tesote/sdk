"""v2 transaction_orders resource: list, get, create, submit, cancel."""

from __future__ import annotations

import json

import pytest

from tesote_sdk.errors import (
    InvalidOrderStateError,
    TransactionOrderNotFoundError,
    ValidationError,
)
from tesote_sdk.v2.transaction_orders import TransactionOrdersResource

from ._helpers import make_transport
from .conftest import http_error, ok_response


def _order(oid: str, status: str = "draft") -> dict:
    return {
        "id": oid,
        "status": status,
        "amount": 100.0,
        "currency": "VES",
        "description": "test",
        "source_account": {"id": "acct_1", "name": "Op", "payment_method_id": "pm_src"},
        "destination": {
            "payment_method_id": "pm_dst",
            "counterparty_id": "cp_1",
            "counterparty_name": "Vendor",
        },
        "fee": {"amount": 1.5, "currency": "VES"},
        "created_at": "2026-04-01T00:00:00Z",
        "updated_at": "2026-04-01T00:00:00Z",
    }


def test_list_returns_paginated_orders() -> None:
    transport, opener = make_transport(
        [
            ok_response(
                {
                    "items": [_order("o1"), _order("o2", "approved")],
                    "has_more": False,
                    "limit": 50,
                    "offset": 0,
                }
            )
        ]
    )
    result = TransactionOrdersResource(transport).list(
        "acct_1", limit=50, status="draft", batch_id="b_1"
    )
    assert len(result.items) == 2
    assert result.items[1].status == "approved"
    url = opener.calls[0]["url"]
    assert "status=draft" in url
    assert "batch_id=b_1" in url


def test_iter_follows_offset() -> None:
    p1 = ok_response(
        {"items": [_order("o1"), _order("o2")], "has_more": True, "limit": 2, "offset": 0}
    )
    p2 = ok_response(
        {"items": [_order("o3")], "has_more": False, "limit": 2, "offset": 2}
    )
    transport, opener = make_transport([p1, p2])
    orders = list(TransactionOrdersResource(transport).iter("acct_1", page_size=2))
    assert [o.id for o in orders] == ["o1", "o2", "o3"]
    assert "offset=2" in opener.calls[1]["url"]


def test_get_404_raises_order_not_found() -> None:
    transport, _ = make_transport(
        [
            http_error(
                404, {}, {"error": "no", "error_code": "TRANSACTION_ORDER_NOT_FOUND"}
            )
        ]
    )
    with pytest.raises(TransactionOrderNotFoundError):
        TransactionOrdersResource(transport).get("acct_1", "missing")


def test_create_with_payment_method_id_wraps_body_and_forwards_idempotency_key() -> None:
    transport, opener = make_transport([ok_response(_order("o1"), status=201)])
    TransactionOrdersResource(transport).create(
        "acct_1",
        amount="100.00",
        currency="VES",
        description="rent",
        destination_payment_method_id="pm_dst",
        idempotency_key="ik-001",
    )
    call = opener.calls[0]
    body = json.loads(call["body"].decode("utf-8"))
    assert body == {
        "transaction_order": {
            "amount": "100.00",
            "currency": "VES",
            "description": "rent",
            "destination_payment_method_id": "pm_dst",
            "idempotency_key": "ik-001",
        }
    }
    headers = dict(call["headers"])
    assert headers["Idempotency-key"] == "ik-001"


def test_create_with_beneficiary_dict_includes_it_in_body() -> None:
    transport, opener = make_transport([ok_response(_order("o1"), status=201)])
    TransactionOrdersResource(transport).create(
        "acct_1",
        amount="50.00",
        currency="VES",
        description="payout",
        beneficiary={"name": "Vendor", "bank_code": "0102", "account_number": "1234"},
        scheduled_for="2026-05-01T00:00:00Z",
        metadata={"po": "PO-1"},
    )
    body = json.loads(opener.calls[0]["body"].decode("utf-8"))
    inner = body["transaction_order"]
    assert inner["beneficiary"]["name"] == "Vendor"
    assert inner["scheduled_for"] == "2026-05-01T00:00:00Z"
    assert inner["metadata"] == {"po": "PO-1"}


def test_create_validation_error_raises_typed_error() -> None:
    transport, _ = make_transport(
        [http_error(400, {}, {"error": "bad amount", "error_code": "VALIDATION_ERROR"})]
    )
    with pytest.raises(ValidationError):
        TransactionOrdersResource(transport).create(
            "acct_1", amount="-1", currency="VES", description="bad"
        )


def test_submit_passes_token_when_provided() -> None:
    transport, opener = make_transport(
        [ok_response(_order("o1", "processing"), status=202)]
    )
    TransactionOrdersResource(transport).submit("acct_1", "o1", token="otp-123")
    body = json.loads(opener.calls[0]["body"].decode("utf-8"))
    assert body == {"token": "otp-123"}
    assert opener.calls[0]["url"].endswith("/v2/accounts/acct_1/transaction_orders/o1/submit")


def test_submit_invalid_order_state() -> None:
    transport, _ = make_transport(
        [http_error(409, {}, {"error": "wrong state", "error_code": "INVALID_ORDER_STATE"})]
    )
    with pytest.raises(InvalidOrderStateError):
        TransactionOrdersResource(transport).submit("acct_1", "o1")


def test_cancel_posts_empty_body_to_cancel_path() -> None:
    transport, opener = make_transport([ok_response(_order("o1", "cancelled"))])
    res = TransactionOrdersResource(transport).cancel("acct_1", "o1")
    assert res.status == "cancelled"
    assert opener.calls[0]["url"].endswith("/v2/accounts/acct_1/transaction_orders/o1/cancel")
    body = json.loads(opener.calls[0]["body"].decode("utf-8"))
    assert body == {}


def test_415_when_content_type_not_json_maps_to_unprocessable() -> None:
    """Ensures the 415 path is wired -- even though the SDK always sets the header,
    a server-side 415 still propagates through the typed-error map."""
    from tesote_sdk.errors import UnprocessableContentError

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
        TransactionOrdersResource(transport).cancel("acct_1", "o1")
