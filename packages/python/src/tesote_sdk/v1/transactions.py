"""v1 transactions resource. Read-only.

Endpoints: `/v1/accounts/{id}/transactions`, `/v1/transactions/{id}`.
"""

from __future__ import annotations

from typing import Any, Dict, Iterator, Optional

from ..models import Transaction, TransactionList
from ..transport import Transport


class TransactionsResource:
    """v1 transactions: list-for-account (cursor) + show-by-id."""

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def list_for_account(
        self,
        account_id: str,
        *,
        start_date: Optional[str] = None,
        end_date: Optional[str] = None,
        scope: Optional[str] = None,
        page: Optional[int] = None,
        per_page: Optional[int] = None,
        transactions_after_id: Optional[str] = None,
        transactions_before_id: Optional[str] = None,
        cache_ttl: Optional[float] = None,
    ) -> TransactionList:
        """GET /v1/accounts/{id}/transactions."""
        query: Dict[str, Any] = {
            "start_date": start_date,
            "end_date": end_date,
            "scope": scope,
            "page": page,
            "per_page": per_page,
            "transactions_after_id": transactions_after_id,
            "transactions_before_id": transactions_before_id,
        }
        response = self._transport.request(
            "GET",
            f"/v1/accounts/{account_id}/transactions",
            query=query,
            cache_ttl=cache_ttl,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return TransactionList.from_dict(body)

    def iter_for_account(
        self,
        account_id: str,
        *,
        start_date: Optional[str] = None,
        end_date: Optional[str] = None,
        scope: Optional[str] = None,
        per_page: Optional[int] = None,
    ) -> Iterator[Transaction]:
        """Iterate every transaction for the account, following the cursor.

        Stops when the server signals ``has_more=False`` or returns an empty page.
        """
        after: Optional[str] = None
        while True:
            page = self.list_for_account(
                account_id,
                start_date=start_date,
                end_date=end_date,
                scope=scope,
                per_page=per_page,
                transactions_after_id=after,
            )
            if not page.transactions:
                return
            yield from page.transactions
            cursor = page.pagination
            if cursor is None or not cursor.has_more or not cursor.after_id:
                return
            after = cursor.after_id

    def get(self, transaction_id: str, *, cache_ttl: Optional[float] = None) -> Transaction:
        """GET /v1/transactions/{id}."""
        response = self._transport.request(
            "GET", f"/v1/transactions/{transaction_id}", cache_ttl=cache_ttl
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return Transaction.from_dict(body)


__all__ = ["TransactionsResource"]
