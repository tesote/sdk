"""v1 accounts resource. Read-only on `/v1/accounts`."""

from __future__ import annotations

from typing import Any, Dict, Optional

from ..models import Account, AccountList
from ..transport import Transport


class AccountsResource:
    """Read-only accounts on `/v1/accounts`."""

    _PREFIX = "/v1/accounts"

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
        """GET /v1/accounts -- page-based pagination, ETag-cached for 60s upstream."""
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
        """GET /v1/accounts/{id}."""
        response = self._transport.request(
            "GET", f"{self._PREFIX}/{account_id}", cache_ttl=cache_ttl
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return Account.from_dict(body)


__all__ = ["AccountsResource"]
