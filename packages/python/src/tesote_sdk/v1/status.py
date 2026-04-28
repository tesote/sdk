"""v1 status / whoami stub."""

from __future__ import annotations

from typing import Any, Dict

from ..transport import Transport


class StatusResource:
    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def status(self) -> Dict[str, Any]:
        raise NotImplementedError("v1 status not yet implemented")

    def whoami(self) -> Dict[str, Any]:
        raise NotImplementedError("v1 whoami not yet implemented")


__all__ = ["StatusResource"]
