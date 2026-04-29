"""v1 status / whoami endpoints (`/status`, `/whoami`)."""

from __future__ import annotations

from typing import Any, Dict

from ..models import StatusResponse, WhoAmI
from ..transport import Transport


class StatusResource:
    """`/status` (no auth) and `/whoami` (auth)."""

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def status(self) -> StatusResponse:
        """GET /status. Always succeeds."""
        response = self._transport.request("GET", "/status")
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return StatusResponse.from_dict(body)

    def whoami(self) -> WhoAmI:
        """GET /whoami. Requires auth."""
        response = self._transport.request("GET", "/whoami")
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return WhoAmI.from_dict(body)


__all__ = ["StatusResource"]
