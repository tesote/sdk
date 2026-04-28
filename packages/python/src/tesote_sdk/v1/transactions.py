"""v1 transactions resource. Read-only stub."""

from __future__ import annotations

from typing import Any, Dict, List

from ..transport import Transport


class TransactionsResource:
    """Stub: methods raise ``NotImplementedError`` until wired."""

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def list_for_account(self, account_id: str) -> List[Dict[str, Any]]:
        raise NotImplementedError("v1 transactions.list_for_account not yet implemented")

    def get(self, transaction_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v1 transactions.get not yet implemented")


__all__ = ["TransactionsResource"]
