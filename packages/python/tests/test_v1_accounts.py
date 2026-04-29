"""v1 accounts resource."""

from __future__ import annotations

import pytest

from tesote_sdk.errors import AccountNotFoundError, UnauthorizedError
from tesote_sdk.v1.accounts import AccountsResource

from ._helpers import make_transport
from .conftest import http_error, ok_response

_ACCOUNT = {
    "id": "acct_1",
    "name": "Operating",
    "data": {
        "masked_account_number": "****1234",
        "currency": "VES",
        "balance_cents": "1000000",
    },
    "bank": {"name": "Banco X"},
    "legal_entity": {"id": "le_1", "legal_name": "Acme SA"},
    "tesote_created_at": "2026-01-01T00:00:00Z",
    "tesote_updated_at": "2026-01-02T00:00:00Z",
}


def test_list_returns_typed_account_list() -> None:
    payload = {
        "total": 1,
        "accounts": [_ACCOUNT],
        "pagination": {
            "current_page": 1,
            "per_page": 50,
            "total_pages": 1,
            "total_count": 1,
        },
    }
    transport, opener = make_transport([ok_response(payload)])
    result = AccountsResource(transport).list(page=1, per_page=50, sort="name")
    assert result.total == 1
    assert len(result.accounts) == 1
    assert result.accounts[0].id == "acct_1"
    assert result.accounts[0].bank.name == "Banco X"
    assert result.accounts[0].data.currency == "VES"
    assert result.pagination is not None
    assert result.pagination.current_page == 1
    url = opener.calls[0]["url"]
    assert "page=1" in url
    assert "per_page=50" in url
    assert "sort=name" in url


def test_list_drops_none_query_params() -> None:
    transport, opener = make_transport([ok_response({"total": 0, "accounts": []})])
    AccountsResource(transport).list()
    # url has no query string at all
    assert "?" not in opener.calls[0]["url"]


def test_get_returns_account_model() -> None:
    transport, _ = make_transport([ok_response(_ACCOUNT)])
    account = AccountsResource(transport).get("acct_1")
    assert account.id == "acct_1"
    assert account.name == "Operating"
    assert account.legal_entity.legal_name == "Acme SA"


def test_get_404_raises_account_not_found() -> None:
    transport, _ = make_transport(
        [http_error(404, {}, {"error": "missing", "error_code": "ACCOUNT_NOT_FOUND"})]
    )
    with pytest.raises(AccountNotFoundError):
        AccountsResource(transport).get("acct_missing")


def test_list_unauthorized_raises_typed_error() -> None:
    transport, _ = make_transport(
        [http_error(401, {}, {"error": "no key", "error_code": "UNAUTHORIZED"})]
    )
    with pytest.raises(UnauthorizedError):
        AccountsResource(transport).list()
