"""v2 accounts resource: list, get, sync."""

from __future__ import annotations

from typing import Any, Dict, Optional

from ..models import Account, AccountList, AccountSyncStarted
from ..transport import Transport


class AccountsResource:
    """`/v2/accounts` -- list, get, sync."""

    _PREFIX = "/v2/accounts"

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def list(
        self,
        *,
        page: Optional[int] = None,
        per_page: Optional[int] = None,
        include: Optional[str] = None,
        sort: Optional[str] = None,
        cache_ttl: Optional[float] = None,
    ) -> AccountList:
        query: Dict[str, Any] = {
            "page": page,
            "per_page": per_page,
            "include": include,
            "sort": sort,
        }
        response = self._transport.request("GET", self._PREFIX, query=query, cache_ttl=cache_ttl)
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return AccountList.from_dict(body)

    def get(self, account_id: str, *, cache_ttl: Optional[float] = None) -> Account:
        response = self._transport.request(
            "GET", f"{self._PREFIX}/{account_id}", cache_ttl=cache_ttl
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return Account.from_dict(body)

    def sync(
        self,
        account_id: str,
        *,
        idempotency_key: Optional[str] = None,
    ) -> AccountSyncStarted:
        """POST /v2/accounts/{id}/sync. 202 Accepted; idempotency-key header optional."""
        response = self._transport.request(
            "POST",
            f"{self._PREFIX}/{account_id}/sync",
            body={},
            idempotency_key=idempotency_key,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return AccountSyncStarted.from_dict(body)


__all__ = ["AccountsResource"]
