"""v1 accounts resource. Read-only."""

from __future__ import annotations

from typing import Any, Dict, List, Optional

from ..transport import Transport


class AccountsResource:
    """Read-only accounts on `/v1/accounts`."""

    _PREFIX = "/v1/accounts"

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def list(self, *, cache_ttl: Optional[float] = None) -> List[Dict[str, Any]]:
        response = self._transport.request(
            "GET", self._PREFIX, cache_ttl=cache_ttl
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


__all__ = ["AccountsResource"]
