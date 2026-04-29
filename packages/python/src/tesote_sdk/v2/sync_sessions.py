"""v2 sync_sessions resource: list + get per account."""

from __future__ import annotations

from typing import Any, Dict, Iterator, Optional

from ..models import SyncSession, SyncSessionList
from ..transport import Transport


class SyncSessionsResource:
    """`/v2/accounts/{id}/sync_sessions` -- list + show."""

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def list(
        self,
        account_id: str,
        *,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
        status: Optional[str] = None,
        cache_ttl: Optional[float] = None,
    ) -> SyncSessionList:
        query: Dict[str, Any] = {
            "limit": limit,
            "offset": offset,
            "status": status,
        }
        response = self._transport.request(
            "GET",
            f"/v2/accounts/{account_id}/sync_sessions",
            query=query,
            cache_ttl=cache_ttl,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return SyncSessionList.from_dict(body)

    def iter(
        self,
        account_id: str,
        *,
        status: Optional[str] = None,
        page_size: int = 50,
    ) -> Iterator[SyncSession]:
        """Iterate every sync session for the account using offset pagination."""
        offset = 0
        while True:
            page = self.list(
                account_id, limit=page_size, offset=offset, status=status
            )
            if not page.sync_sessions:
                return
            yield from page.sync_sessions
            if not page.has_more:
                return
            offset += len(page.sync_sessions)

    def get(
        self,
        account_id: str,
        session_id: str,
        *,
        cache_ttl: Optional[float] = None,
    ) -> SyncSession:
        response = self._transport.request(
            "GET",
            f"/v2/accounts/{account_id}/sync_sessions/{session_id}",
            cache_ttl=cache_ttl,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return SyncSession.from_dict(body)


__all__ = ["SyncSessionsResource"]
