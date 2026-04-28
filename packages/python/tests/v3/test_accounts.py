"""End-to-end wiring for v3 accounts.list / accounts.get."""

from __future__ import annotations

import pytest

from tesote_sdk import V3Client
from tesote_sdk.errors import ConfigError
from tesote_sdk.transport import RetryPolicy

from ..conftest import ScriptedOpener, ok_response


def _client(opener: ScriptedOpener) -> V3Client:
    c = V3Client(api_key="sk_test_abcdef1234")
    # why: swap out opener post-construction so the v3 client API stays clean
    c._transport._opener = opener  # type: ignore[attr-defined]
    c._transport._retry_policy = RetryPolicy(max_attempts=2, base_delay=0.0, max_delay=0.0)  # type: ignore[attr-defined]
    return c


def test_missing_api_key_raises_config_error() -> None:
    with pytest.raises(ConfigError):
        V3Client(api_key="")


def test_accounts_list_unwraps_data_envelope() -> None:
    opener = ScriptedOpener(
        [ok_response({"data": [{"id": "acct_1"}, {"id": "acct_2"}]})]
    )
    c = _client(opener)
    accounts = c.accounts.list()
    assert accounts == [{"id": "acct_1"}, {"id": "acct_2"}]
    assert opener.calls[0]["method"] == "GET"
    assert "/v3/accounts" in opener.calls[0]["url"]


def test_accounts_list_handles_bare_array() -> None:
    opener = ScriptedOpener([ok_response([{"id": "acct_1"}])])
    c = _client(opener)
    assert c.accounts.list() == [{"id": "acct_1"}]


def test_accounts_list_passes_cursor_and_limit() -> None:
    opener = ScriptedOpener([ok_response({"data": []})])
    c = _client(opener)
    c.accounts.list(cursor="abc", limit=50)
    url = opener.calls[0]["url"]
    assert "cursor=abc" in url
    assert "limit=50" in url


def test_accounts_get_returns_dict() -> None:
    opener = ScriptedOpener([ok_response({"id": "acct_1", "currency": "USD"})])
    c = _client(opener)
    acct = c.accounts.get("acct_1")
    assert acct == {"id": "acct_1", "currency": "USD"}
    assert opener.calls[0]["url"].endswith("/v3/accounts/acct_1")


def test_unimplemented_methods_raise_not_implemented() -> None:
    opener = ScriptedOpener([])
    c = _client(opener)
    with pytest.raises(NotImplementedError):
        c.accounts.sync("acct_1")
