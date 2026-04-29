"""v2 transactions resource: per-account list, sync, search, bulk, export, get-by-id."""

from __future__ import annotations

from typing import Any, Dict, Iterator, List, Mapping, Optional

from ..models import (
    BulkResult,
    SearchResult,
    SyncDelta,
    Transaction,
    TransactionList,
)
from ..transport import Transport


class TransactionsResource:
    """v2 transactions: filtering on /v2/accounts/{id}/transactions, sync, search, bulk, export."""

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    # -- list / iter for account ----------------------------------------

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
        transaction_date_after: Optional[str] = None,
        transaction_date_before: Optional[str] = None,
        created_after: Optional[str] = None,
        updated_after: Optional[str] = None,
        amount_min: Optional[float] = None,
        amount_max: Optional[float] = None,
        amount: Optional[float] = None,
        status: Optional[str] = None,
        category_id: Optional[str] = None,
        counterparty_id: Optional[str] = None,
        q: Optional[str] = None,
        type: Optional[str] = None,  # noqa: A002 -- matches API param name
        reference_code: Optional[str] = None,
        cache_ttl: Optional[float] = None,
    ) -> TransactionList:
        """GET /v2/accounts/{id}/transactions."""
        query: Dict[str, Any] = {
            "start_date": start_date,
            "end_date": end_date,
            "scope": scope,
            "page": page,
            "per_page": per_page,
            "transactions_after_id": transactions_after_id,
            "transactions_before_id": transactions_before_id,
            "transaction_date_after": transaction_date_after,
            "transaction_date_before": transaction_date_before,
            "created_after": created_after,
            "updated_after": updated_after,
            "amount_min": amount_min,
            "amount_max": amount_max,
            "amount": amount,
            "status": status,
            "category_id": category_id,
            "counterparty_id": counterparty_id,
            "q": q,
            "type": type,
            "reference_code": reference_code,
        }
        response = self._transport.request(
            "GET",
            f"/v2/accounts/{account_id}/transactions",
            query=query,
            cache_ttl=cache_ttl,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return TransactionList.from_dict(body)

    def iter_for_account(
        self,
        account_id: str,
        *,
        per_page: Optional[int] = None,
        **filters: Any,
    ) -> Iterator[Transaction]:
        """Iterate every transaction for the account, following the cursor."""
        after: Optional[str] = None
        while True:
            page = self.list_for_account(
                account_id,
                per_page=per_page,
                transactions_after_id=after,
                **filters,
            )
            if not page.transactions:
                return
            yield from page.transactions
            cursor = page.pagination
            if cursor is None or not cursor.has_more or not cursor.after_id:
                return
            after = cursor.after_id

    # -- export ---------------------------------------------------------

    def export(
        self,
        account_id: str,
        *,
        format: str = "csv",  # noqa: A002 -- matches API param
        start_date: Optional[str] = None,
        end_date: Optional[str] = None,
        **filters: Any,
    ) -> str:
        """GET /v2/accounts/{id}/transactions/export. Returns the file body as text."""
        query: Dict[str, Any] = {
            "format": format,
            "start_date": start_date,
            "end_date": end_date,
        }
        for k, v in filters.items():
            query[k] = v
        response = self._transport.request(
            "GET",
            f"/v2/accounts/{account_id}/transactions/export",
            query=query,
        )
        return response.body

    # -- get by id ------------------------------------------------------

    def get(self, transaction_id: str, *, cache_ttl: Optional[float] = None) -> Transaction:
        """GET /v2/transactions/{id} -- v1 schema (not SyncTransaction)."""
        response = self._transport.request(
            "GET", f"/v2/transactions/{transaction_id}", cache_ttl=cache_ttl
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return Transaction.from_dict(body)

    # -- sync (per-account + legacy) ------------------------------------

    def sync(
        self,
        account_id: str,
        *,
        count: Optional[int] = None,
        cursor: Optional[str] = None,
        include_running_balance: Optional[bool] = None,
        idempotency_key: Optional[str] = None,
    ) -> SyncDelta:
        """POST /v2/accounts/{id}/transactions/sync."""
        body = self._build_sync_body(count, cursor, include_running_balance)
        response = self._transport.request(
            "POST",
            f"/v2/accounts/{account_id}/transactions/sync",
            body=body,
            idempotency_key=idempotency_key,
        )
        payload: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return SyncDelta.from_dict(payload)

    def sync_legacy(
        self,
        *,
        account_id: Optional[str] = None,
        count: Optional[int] = None,
        cursor: Optional[str] = None,
        include_running_balance: Optional[bool] = None,
        idempotency_key: Optional[str] = None,
    ) -> SyncDelta:
        """POST /v2/transactions/sync -- legacy non-nested route."""
        body = self._build_sync_body(count, cursor, include_running_balance)
        if account_id is not None:
            body["account_id"] = account_id
        response = self._transport.request(
            "POST",
            "/v2/transactions/sync",
            body=body,
            idempotency_key=idempotency_key,
        )
        payload: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return SyncDelta.from_dict(payload)

    @staticmethod
    def _build_sync_body(
        count: Optional[int],
        cursor: Optional[str],
        include_running_balance: Optional[bool],
    ) -> Dict[str, Any]:
        body: Dict[str, Any] = {}
        if count is not None:
            body["count"] = count
        if cursor is not None:
            body["cursor"] = cursor
        if include_running_balance is not None:
            body["options"] = {"include_running_balance": bool(include_running_balance)}
        return body

    # -- bulk -----------------------------------------------------------

    def bulk(
        self,
        account_ids: List[str],
        *,
        page: Optional[int] = None,
        per_page: Optional[int] = None,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
        idempotency_key: Optional[str] = None,
    ) -> BulkResult:
        """POST /v2/transactions/bulk -- max 100 accounts."""
        body: Dict[str, Any] = {"account_ids": list(account_ids)}
        if page is not None:
            body["page"] = page
        if per_page is not None:
            body["per_page"] = per_page
        if limit is not None:
            body["limit"] = limit
        if offset is not None:
            body["offset"] = offset
        response = self._transport.request(
            "POST", "/v2/transactions/bulk", body=body, idempotency_key=idempotency_key
        )
        payload: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return BulkResult.from_dict(payload)

    # -- search ---------------------------------------------------------

    def search(
        self,
        q: str,
        *,
        account_id: Optional[str] = None,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
        filters: Optional[Mapping[str, Any]] = None,
        cache_ttl: Optional[float] = None,
    ) -> SearchResult:
        """GET /v2/transactions/search. ``q`` is required."""
        query: Dict[str, Any] = {
            "q": q,
            "account_id": account_id,
            "limit": limit,
            "offset": offset,
        }
        if filters:
            for k, v in filters.items():
                query[k] = v
        response = self._transport.request(
            "GET", "/v2/transactions/search", query=query, cache_ttl=cache_ttl
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return SearchResult.from_dict(body)


__all__ = ["TransactionsResource"]
