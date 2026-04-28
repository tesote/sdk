"""v3 accounts resource. ``list`` and ``get`` wired end-to-end."""

from __future__ import annotations

from typing import Any, Dict, List, Optional

from ..transport import Transport


class AccountsResource:
    """`/v3/accounts`. Mirrors v2 with the v3 prefix."""

    _PREFIX = "/v3/accounts"

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def list(
        self,
        *,
        cursor: Optional[str] = None,
        limit: Optional[int] = None,
        cache_ttl: Optional[float] = None,
    ) -> List[Dict[str, Any]]:
        query: Dict[str, Any] = {}
        if cursor is not None:
            query["cursor"] = cursor
        if limit is not None:
            query["limit"] = limit
        response = self._transport.request(
            "GET", self._PREFIX, query=query or None, cache_ttl=cache_ttl
        )
        body = response.json
        if isinstance(body, list):
            return [item for item in body if isinstance(item, dict)]
        if isinstance(body, dict) and isinstance(body.get("data"), list):
            return [item for item in body["data"] if isinstance(item, dict)]
        return []

    def get(self, account_id: str, *, cache_ttl: Optional[float] = None) -> Dict[str, Any]:
        response = self._transport.request(
            "GET", f"{self._PREFIX}/{account_id}", cache_ttl=cache_ttl
        )
        body = response.json
        if isinstance(body, dict):
            return body
        return {}

    def sync(self, account_id: str, *, idempotency_key: Optional[str] = None) -> Dict[str, Any]:
        raise NotImplementedError("v3 accounts.sync not yet implemented")


__all__ = ["AccountsResource"]
