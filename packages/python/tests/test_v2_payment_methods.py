"""v2 payment_methods resource: list, get, create, update, delete."""

from __future__ import annotations

import json

import pytest

from tesote_sdk.errors import (
    PaymentMethodNotFoundError,
    UnprocessableContentError,
    ValidationError,
)
from tesote_sdk.v2.payment_methods import PaymentMethodsResource

from ._helpers import make_transport
from .conftest import http_error, ok_response


def _pm(pid: str = "pm_1") -> dict:
    return {
        "id": pid,
        "method_type": "bank_account",
        "currency": "VES",
        "label": "main",
        "details": {
            "bank_code": "0102",
            "account_number": "1234",
            "holder_name": "Acme",
        },
        "verified": True,
        "verified_at": "2026-01-01T00:00:00Z",
        "counterparty": {"id": "cp_1", "name": "Vendor"},
        "tesote_account": None,
        "created_at": "2026-01-01T00:00:00Z",
        "updated_at": "2026-01-01T00:00:00Z",
    }


def test_list_serializes_filters_including_verified_bool_to_string() -> None:
    transport, opener = make_transport(
        [ok_response({"items": [_pm()], "has_more": False, "limit": 50, "offset": 0})]
    )
    PaymentMethodsResource(transport).list(
        method_type="bank_account",
        currency="VES",
        counterparty_id="cp_1",
        verified=True,
    )
    url = opener.calls[0]["url"]
    assert "method_type=bank_account" in url
    assert "currency=VES" in url
    assert "counterparty_id=cp_1" in url
    assert "verified=true" in url


def test_iter_follows_offset() -> None:
    p1 = ok_response(
        {"items": [_pm("pm_1"), _pm("pm_2")], "has_more": True, "limit": 2, "offset": 0}
    )
    p2 = ok_response(
        {"items": [_pm("pm_3")], "has_more": False, "limit": 2, "offset": 2}
    )
    transport, opener = make_transport([p1, p2])
    items = list(PaymentMethodsResource(transport).iter(page_size=2))
    assert [p.id for p in items] == ["pm_1", "pm_2", "pm_3"]
    assert "offset=2" in opener.calls[1]["url"]


def test_get_404_payment_method_not_found() -> None:
    transport, _ = make_transport(
        [
            http_error(
                404, {}, {"error": "missing", "error_code": "PAYMENT_METHOD_NOT_FOUND"}
            )
        ]
    )
    with pytest.raises(PaymentMethodNotFoundError):
        PaymentMethodsResource(transport).get("pm_missing")


def test_create_wraps_body_with_payment_method_key_and_sets_idempotency() -> None:
    transport, opener = make_transport([ok_response(_pm(), status=201)])
    PaymentMethodsResource(transport).create(
        method_type="bank_account",
        currency="VES",
        details={"bank_code": "0102", "account_number": "1234", "holder_name": "Acme"},
        label="main",
        counterparty={"name": "Vendor"},
        idempotency_key="pm-key",
    )
    body = json.loads(opener.calls[0]["body"].decode("utf-8"))
    assert body["payment_method"]["method_type"] == "bank_account"
    assert body["payment_method"]["counterparty"] == {"name": "Vendor"}
    headers = dict(opener.calls[0]["headers"])
    assert headers["Idempotency-key"] == "pm-key"
    assert headers["Content-type"] == "application/json"


def test_create_validation_error() -> None:
    transport, _ = make_transport(
        [http_error(400, {}, {"error": "bad", "error_code": "VALIDATION_ERROR"})]
    )
    with pytest.raises(ValidationError):
        PaymentMethodsResource(transport).create(
            method_type="bank_account", currency="VES", details={}
        )


def test_update_only_sends_provided_fields() -> None:
    transport, opener = make_transport([ok_response(_pm())])
    PaymentMethodsResource(transport).update("pm_1", label="renamed")
    assert opener.calls[0]["method"] == "PATCH"
    body = json.loads(opener.calls[0]["body"].decode("utf-8"))
    assert body == {"payment_method": {"label": "renamed"}}


def test_delete_204_returns_none_and_idempotency_propagates() -> None:
    transport, opener = make_transport([ok_response(b"", status=204)])
    PaymentMethodsResource(transport).delete("pm_1", idempotency_key="del-key")
    assert opener.calls[0]["method"] == "DELETE"
    headers = dict(opener.calls[0]["headers"])
    assert headers["Idempotency-key"] == "del-key"


def test_delete_409_in_use_maps_to_validation_error() -> None:
    """Spec says 409 VALIDATION_ERROR when payment method has active orders."""
    transport, _ = make_transport(
        [http_error(409, {}, {"error": "in use", "error_code": "VALIDATION_ERROR"})]
    )
    with pytest.raises(ValidationError):
        PaymentMethodsResource(transport).delete("pm_1")


def test_415_when_content_type_missing_maps_to_unprocessable() -> None:
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
        PaymentMethodsResource(transport).create(
            method_type="bank_account", currency="VES", details={}
        )
