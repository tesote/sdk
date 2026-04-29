"""v2 transaction_orders resource: list, get, create, submit, cancel."""

from __future__ import annotations

from typing import Any, Dict, Iterator, Mapping, Optional

from ..models import TransactionOrder, TransactionOrderList
from ..transport import Transport


class TransactionOrdersResource:
    """`/v2/accounts/{id}/transaction_orders`."""

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def list(
        self,
        account_id: str,
        *,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
        status: Optional[str] = None,
        created_after: Optional[str] = None,
        created_before: Optional[str] = None,
        batch_id: Optional[str] = None,
        cache_ttl: Optional[float] = None,
    ) -> TransactionOrderList:
        query: Dict[str, Any] = {
            "limit": limit,
            "offset": offset,
            "status": status,
            "created_after": created_after,
            "created_before": created_before,
            "batch_id": batch_id,
        }
        response = self._transport.request(
            "GET",
            f"/v2/accounts/{account_id}/transaction_orders",
            query=query,
            cache_ttl=cache_ttl,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return TransactionOrderList.from_dict(body)

    def iter(
        self,
        account_id: str,
        *,
        status: Optional[str] = None,
        batch_id: Optional[str] = None,
        page_size: int = 50,
    ) -> Iterator[TransactionOrder]:
        """Iterate every transaction order for the account via offset pagination."""
        offset = 0
        while True:
            page = self.list(
                account_id,
                limit=page_size,
                offset=offset,
                status=status,
                batch_id=batch_id,
            )
            if not page.items:
                return
            yield from page.items
            if not page.has_more:
                return
            offset += len(page.items)

    def get(
        self,
        account_id: str,
        order_id: str,
        *,
        cache_ttl: Optional[float] = None,
    ) -> TransactionOrder:
        response = self._transport.request(
            "GET",
            f"/v2/accounts/{account_id}/transaction_orders/{order_id}",
            cache_ttl=cache_ttl,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return TransactionOrder.from_dict(body)

    def create(
        self,
        account_id: str,
        *,
        amount: str,
        currency: str,
        description: str,
        destination_payment_method_id: Optional[str] = None,
        beneficiary: Optional[Mapping[str, Any]] = None,
        scheduled_for: Optional[str] = None,
        metadata: Optional[Mapping[str, Any]] = None,
        idempotency_key: Optional[str] = None,
    ) -> TransactionOrder:
        """POST /v2/accounts/{id}/transaction_orders.

        Pass ``destination_payment_method_id`` OR a ``beneficiary`` dict (server
        creates the on-the-fly PaymentMethod). ``idempotency_key`` is forwarded
        both to the request body (for server-side dedupe) and as the transport
        ``Idempotency-Key`` header.
        """
        order_body: Dict[str, Any] = {
            "amount": amount,
            "currency": currency,
            "description": description,
        }
        if destination_payment_method_id is not None:
            order_body["destination_payment_method_id"] = destination_payment_method_id
        if beneficiary is not None:
            order_body["beneficiary"] = dict(beneficiary)
        if scheduled_for is not None:
            order_body["scheduled_for"] = scheduled_for
        if metadata is not None:
            order_body["metadata"] = dict(metadata)
        if idempotency_key is not None:
            order_body["idempotency_key"] = idempotency_key
        response = self._transport.request(
            "POST",
            f"/v2/accounts/{account_id}/transaction_orders",
            body={"transaction_order": order_body},
            idempotency_key=idempotency_key,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return TransactionOrder.from_dict(body)

    def submit(
        self,
        account_id: str,
        order_id: str,
        *,
        token: Optional[str] = None,
        idempotency_key: Optional[str] = None,
    ) -> TransactionOrder:
        """POST /v2/accounts/{id}/transaction_orders/{order_id}/submit."""
        body: Dict[str, Any] = {}
        if token is not None:
            body["token"] = token
        response = self._transport.request(
            "POST",
            f"/v2/accounts/{account_id}/transaction_orders/{order_id}/submit",
            body=body,
            idempotency_key=idempotency_key,
        )
        payload: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return TransactionOrder.from_dict(payload)

    def cancel(
        self,
        account_id: str,
        order_id: str,
        *,
        idempotency_key: Optional[str] = None,
    ) -> TransactionOrder:
        """POST /v2/accounts/{id}/transaction_orders/{order_id}/cancel."""
        response = self._transport.request(
            "POST",
            f"/v2/accounts/{account_id}/transaction_orders/{order_id}/cancel",
            body={},
            idempotency_key=idempotency_key,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return TransactionOrder.from_dict(body)


__all__ = ["TransactionOrdersResource"]
