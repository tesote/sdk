"""v2 status / whoami endpoints (`/v2/status`, `/v2/whoami`)."""

from __future__ import annotations

from typing import Any, Dict

from ..models import StatusResponse, WhoAmI
from ..transport import Transport


class StatusResource:
    """`/v2/status` (no auth) and `/v2/whoami` (auth)."""

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def status(self) -> StatusResponse:
        response = self._transport.request("GET", "/v2/status")
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return StatusResponse.from_dict(body)

    def whoami(self) -> WhoAmI:
        response = self._transport.request("GET", "/v2/whoami")
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return WhoAmI.from_dict(body)


__all__ = ["StatusResource"]
