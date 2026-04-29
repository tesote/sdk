"""v2 transactions resource: list, sync, sync_legacy, bulk, search, export, get."""

from __future__ import annotations

import pytest

from tesote_sdk.errors import (
    HistorySyncForbiddenError,
    InvalidCountError,
    InvalidCursorError,
    TransactionNotFoundError,
    UnprocessableContentError,
)
from tesote_sdk.v2.transactions import TransactionsResource

from ._helpers import make_transport
from .conftest import http_error, ok_response


def _v1_txn(tid: str) -> dict:
    return {
        "id": tid,
        "status": "posted",
        "data": {"amount_cents": 100, "currency": "VES", "description": tid},
        "tesote_imported_at": "2026-04-01T00:00:00Z",
        "tesote_updated_at": "2026-04-01T00:00:00Z",
    }


def _sync_txn(tid: str) -> dict:
    return {
        "transaction_id": tid,
        "account_id": "acct_1",
        "amount": 12.34,
        "iso_currency_code": "VES",
        "date": "2026-04-01",
        "name": "Coffee",
        "merchant_name": "Cafe",
        "pending": False,
        "category": ["food"],
    }


def test_list_for_account_includes_v2_filters() -> None:
    transport, opener = make_transport(
        [ok_response({"total": 0, "transactions": [], "pagination": {"has_more": False}})]
    )
    TransactionsResource(transport).list_for_account(
        "acct_1",
        amount_min=1.5,
        amount_max=100.0,
        status="posted",
        category_id="cat_1",
        counterparty_id="cp_1",
        q="coffee",
    )
    url = opener.calls[0]["url"]
    assert "amount_min=1.5" in url
    assert "amount_max=100" in url
    assert "status=posted" in url
    assert "category_id=cat_1" in url
    assert "counterparty_id=cp_1" in url
    assert "q=coffee" in url


def test_iter_for_account_follows_cursor() -> None:
    p1 = ok_response(
        {
            "transactions": [_v1_txn("t1"), _v1_txn("t2")],
            "pagination": {"has_more": True, "after_id": "t2", "before_id": "t1"},
        }
    )
    p2 = ok_response(
        {
            "transactions": [_v1_txn("t3")],
            "pagination": {"has_more": False, "after_id": "t3", "before_id": "t3"},
        }
    )
    transport, opener = make_transport([p1, p2])
    txns = list(TransactionsResource(transport).iter_for_account("acct_1", per_page=2))
    assert [t.id for t in txns] == ["t1", "t2", "t3"]
    assert "transactions_after_id=t2" in opener.calls[1]["url"]


def test_get_uses_v2_path() -> None:
    transport, opener = make_transport([ok_response(_v1_txn("t_1"))])
    txn = TransactionsResource(transport).get("t_1")
    assert txn.id == "t_1"
    assert opener.calls[0]["url"].endswith("/v2/transactions/t_1")


def test_get_404() -> None:
    transport, _ = make_transport(
        [
            http_error(
                404, {}, {"error": "missing", "error_code": "TRANSACTION_NOT_FOUND"}
            )
        ]
    )
    with pytest.raises(TransactionNotFoundError):
        TransactionsResource(transport).get("t_missing")


def test_sync_round_trip_with_options_and_idempotency() -> None:
    transport, opener = make_transport(
        [
            ok_response(
                {
                    "added": [_sync_txn("a1")],
                    "modified": [_sync_txn("m1")],
                    "removed": [{"transaction_id": "r1", "account_id": "acct_1"}],
                    "next_cursor": "cur-next",
                    "has_more": False,
                }
            )
        ]
    )
    delta = TransactionsResource(transport).sync(
        "acct_1",
        count=100,
        cursor="now",
        include_running_balance=True,
        idempotency_key="sync-key",
    )
    assert len(delta.added) == 1
    assert delta.added[0].transaction_id == "a1"
    assert delta.removed[0].transaction_id == "r1"
    assert delta.next_cursor == "cur-next"
    assert delta.has_more is False
    headers = dict(opener.calls[0]["headers"])
    assert headers["Idempotency-key"] == "sync-key"
    # body shape: count + cursor + options.include_running_balance
    import json as _json

    body = _json.loads(opener.calls[0]["body"].decode("utf-8"))
    assert body == {
        "count": 100,
        "cursor": "now",
        "options": {"include_running_balance": True},
    }


def test_sync_invalid_count_raises_typed_error() -> None:
    transport, _ = make_transport(
        [http_error(422, {}, {"error": "bad count", "error_code": "INVALID_COUNT"})]
    )
    with pytest.raises(InvalidCountError):
        TransactionsResource(transport).sync("acct_1", count=99999)


def test_sync_invalid_cursor_raises_typed_error() -> None:
    transport, _ = make_transport(
        [http_error(422, {}, {"error": "bad cur", "error_code": "INVALID_CURSOR"})]
    )
    with pytest.raises(InvalidCursorError):
        TransactionsResource(transport).sync("acct_1", cursor="garbage")


def test_sync_history_forbidden_raises_typed_error() -> None:
    transport, _ = make_transport(
        [
            http_error(
                403, {}, {"error": "too old", "error_code": "HISTORY_SYNC_FORBIDDEN"}
            )
        ]
    )
    with pytest.raises(HistorySyncForbiddenError):
        TransactionsResource(transport).sync("acct_1", cursor="ancient")


def test_sync_legacy_uses_non_nested_path_and_passes_account_id_in_body() -> None:
    transport, opener = make_transport(
        [
            ok_response(
                {"added": [], "modified": [], "removed": [], "next_cursor": None, "has_more": False}
            )
        ]
    )
    TransactionsResource(transport).sync_legacy(account_id="acct_1", count=10)
    assert opener.calls[0]["url"].endswith("/v2/transactions/sync")
    import json as _json

    body = _json.loads(opener.calls[0]["body"].decode("utf-8"))
    assert body["account_id"] == "acct_1"
    assert body["count"] == 10


def test_bulk_sends_account_ids_and_returns_typed_results() -> None:
    transport, opener = make_transport(
        [
            ok_response(
                {
                    "bulk_results": [
                        {
                            "account_id": "acct_1",
                            "transactions": [_v1_txn("t1")],
                            "pagination": {"has_more": False, "after_id": "t1"},
                        }
                    ]
                }
            )
        ]
    )
    result = TransactionsResource(transport).bulk(["acct_1"], per_page=10)
    assert len(result.bulk_results) == 1
    assert result.bulk_results[0].account_id == "acct_1"
    assert result.bulk_results[0].transactions[0].id == "t1"
    import json as _json

    body = _json.loads(opener.calls[0]["body"].decode("utf-8"))
    assert body["account_ids"] == ["acct_1"]
    assert body["per_page"] == 10


def test_bulk_empty_account_ids_raises_unprocessable_content() -> None:
    transport, _ = make_transport(
        [
            http_error(
                422, {}, {"error": "no accounts", "error_code": "UNPROCESSABLE_CONTENT"}
            )
        ]
    )
    with pytest.raises(UnprocessableContentError):
        TransactionsResource(transport).bulk([])


def test_search_required_q_in_query_string() -> None:
    transport, opener = make_transport(
        [ok_response({"transactions": [_v1_txn("t1")], "total": 1})]
    )
    result = TransactionsResource(transport).search(
        "coffee", account_id="acct_1", limit=10
    )
    assert result.total == 1
    assert "q=coffee" in opener.calls[0]["url"]
    assert "account_id=acct_1" in opener.calls[0]["url"]


def test_search_missing_q_returns_unprocessable_content() -> None:
    transport, _ = make_transport(
        [http_error(422, {}, {"error": "missing q", "error_code": "UNPROCESSABLE_CONTENT"})]
    )
    with pytest.raises(UnprocessableContentError):
        TransactionsResource(transport).search("")


def test_export_returns_raw_body_and_uses_format_param() -> None:
    csv_body = b"Transaction ID,Date\nt1,2026-04-01\n"
    transport, opener = make_transport(
        [ok_response(csv_body, headers={"Content-Type": "text/csv"})]
    )
    body = TransactionsResource(transport).export("acct_1", format="csv")
    assert "Transaction ID" in body
    url = opener.calls[0]["url"]
    assert "format=csv" in url
    assert "/v2/accounts/acct_1/transactions/export" in url
