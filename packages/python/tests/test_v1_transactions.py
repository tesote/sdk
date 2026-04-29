"""v1 transactions resource."""

from __future__ import annotations

import pytest

from tesote_sdk.errors import (
    InvalidDateRangeError,
    TransactionNotFoundError,
)
from tesote_sdk.v1.transactions import TransactionsResource

from ._helpers import make_transport
from .conftest import http_error, ok_response


def _txn(txn_id: str) -> dict:
    return {
        "id": txn_id,
        "status": "posted",
        "data": {
            "amount_cents": 1000,
            "currency": "VES",
            "description": f"txn {txn_id}",
            "transaction_date": "2026-04-01",
        },
        "tesote_imported_at": "2026-04-01T00:00:00Z",
        "tesote_updated_at": "2026-04-01T00:00:00Z",
        "transaction_categories": [{"name": "groceries"}],
        "counterparty": {"name": "Acme"},
    }


def test_list_for_account_returns_cursor_paginated_list() -> None:
    payload = {
        "total": 2,
        "transactions": [_txn("t_1"), _txn("t_2")],
        "pagination": {
            "has_more": True,
            "per_page": 50,
            "after_id": "t_2",
            "before_id": "t_1",
        },
    }
    transport, opener = make_transport([ok_response(payload)])
    result = TransactionsResource(transport).list_for_account(
        "acct_1", per_page=50, transactions_after_id="prev"
    )
    assert len(result.transactions) == 2
    assert result.transactions[0].counterparty is not None
    assert result.transactions[0].counterparty.name == "Acme"
    assert result.transactions[0].transaction_categories[0].name == "groceries"
    assert result.pagination is not None
    assert result.pagination.has_more is True
    assert result.pagination.after_id == "t_2"
    url = opener.calls[0]["url"]
    assert "transactions_after_id=prev" in url
    assert "per_page=50" in url


def test_iter_for_account_follows_cursor_until_has_more_false() -> None:
    page1 = ok_response(
        {
            "transactions": [_txn("t_1"), _txn("t_2")],
            "pagination": {"has_more": True, "after_id": "t_2", "before_id": "t_1"},
        }
    )
    page2 = ok_response(
        {
            "transactions": [_txn("t_3")],
            "pagination": {"has_more": False, "after_id": "t_3", "before_id": "t_3"},
        }
    )
    transport, opener = make_transport([page1, page2])
    txns = list(TransactionsResource(transport).iter_for_account("acct_1"))
    assert [t.id for t in txns] == ["t_1", "t_2", "t_3"]
    assert len(opener.calls) == 2
    # second call uses after_id from first
    assert "transactions_after_id=t_2" in opener.calls[1]["url"]


def test_iter_stops_when_first_page_empty() -> None:
    transport, opener = make_transport(
        [ok_response({"transactions": [], "pagination": {"has_more": False}})]
    )
    txns = list(TransactionsResource(transport).iter_for_account("acct_1"))
    assert txns == []
    assert len(opener.calls) == 1


def test_get_returns_transaction_model() -> None:
    transport, _ = make_transport([ok_response(_txn("t_1"))])
    txn = TransactionsResource(transport).get("t_1")
    assert txn.id == "t_1"
    assert txn.data.amount_cents == 1000


def test_get_404_raises_transaction_not_found() -> None:
    transport, _ = make_transport(
        [
            http_error(
                404,
                {},
                {"error": "missing", "error_code": "TRANSACTION_NOT_FOUND"},
            )
        ]
    )
    with pytest.raises(TransactionNotFoundError):
        TransactionsResource(transport).get("t_missing")


def test_invalid_date_range_raises_typed_error() -> None:
    transport, _ = make_transport(
        [
            http_error(
                422,
                {},
                {"error": "bad dates", "error_code": "INVALID_DATE_RANGE"},
            )
        ]
    )
    with pytest.raises(InvalidDateRangeError):
        TransactionsResource(transport).list_for_account(
            "acct_1", start_date="bogus", end_date="bogus"
        )
