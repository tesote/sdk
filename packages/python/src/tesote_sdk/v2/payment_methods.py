"""v2 payment_methods resource: list, get, create, update, delete."""

from __future__ import annotations

from typing import Any, Dict, Iterator, Mapping, Optional

from ..models import PaymentMethod, PaymentMethodList
from ..transport import Transport


class PaymentMethodsResource:
    """`/v2/payment_methods`."""

    _PREFIX = "/v2/payment_methods"

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def list(
        self,
        *,
        limit: Optional[int] = None,
        offset: Optional[int] = None,
        method_type: Optional[str] = None,
        currency: Optional[str] = None,
        counterparty_id: Optional[str] = None,
        verified: Optional[bool] = None,
        cache_ttl: Optional[float] = None,
    ) -> PaymentMethodList:
        verified_str: Optional[str] = None
        if verified is not None:
            verified_str = "true" if verified else "false"
        query: Dict[str, Any] = {
            "limit": limit,
            "offset": offset,
            "method_type": method_type,
            "currency": currency,
            "counterparty_id": counterparty_id,
            "verified": verified_str,
        }
        response = self._transport.request(
            "GET", self._PREFIX, query=query, cache_ttl=cache_ttl
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return PaymentMethodList.from_dict(body)

    def iter(
        self,
        *,
        method_type: Optional[str] = None,
        currency: Optional[str] = None,
        counterparty_id: Optional[str] = None,
        verified: Optional[bool] = None,
        page_size: int = 50,
    ) -> Iterator[PaymentMethod]:
        offset = 0
        while True:
            page = self.list(
                limit=page_size,
                offset=offset,
                method_type=method_type,
                currency=currency,
                counterparty_id=counterparty_id,
                verified=verified,
            )
            if not page.items:
                return
            yield from page.items
            if not page.has_more:
                return
            offset += len(page.items)

    def get(
        self,
        payment_method_id: str,
        *,
        cache_ttl: Optional[float] = None,
    ) -> PaymentMethod:
        response = self._transport.request(
            "GET", f"{self._PREFIX}/{payment_method_id}", cache_ttl=cache_ttl
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return PaymentMethod.from_dict(body)

    def create(
        self,
        *,
        method_type: str,
        currency: str,
        details: Mapping[str, Any],
        label: Optional[str] = None,
        counterparty_id: Optional[str] = None,
        counterparty: Optional[Mapping[str, Any]] = None,
        idempotency_key: Optional[str] = None,
    ) -> PaymentMethod:
        """POST /v2/payment_methods. Pass ``counterparty_id`` OR ``counterparty`` (auto-create)."""
        pm_body: Dict[str, Any] = {
            "method_type": method_type,
            "currency": currency,
            "details": dict(details),
        }
        if label is not None:
            pm_body["label"] = label
        if counterparty_id is not None:
            pm_body["counterparty_id"] = counterparty_id
        if counterparty is not None:
            pm_body["counterparty"] = dict(counterparty)
        response = self._transport.request(
            "POST",
            self._PREFIX,
            body={"payment_method": pm_body},
            idempotency_key=idempotency_key,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return PaymentMethod.from_dict(body)

    def update(
        self,
        payment_method_id: str,
        *,
        method_type: Optional[str] = None,
        currency: Optional[str] = None,
        label: Optional[str] = None,
        details: Optional[Mapping[str, Any]] = None,
        counterparty_id: Optional[str] = None,
        counterparty: Optional[Mapping[str, Any]] = None,
        idempotency_key: Optional[str] = None,
    ) -> PaymentMethod:
        """PATCH /v2/payment_methods/{id}. Only sends fields the caller passed."""
        pm_body: Dict[str, Any] = {}
        if method_type is not None:
            pm_body["method_type"] = method_type
        if currency is not None:
            pm_body["currency"] = currency
        if label is not None:
            pm_body["label"] = label
        if details is not None:
            pm_body["details"] = dict(details)
        if counterparty_id is not None:
            pm_body["counterparty_id"] = counterparty_id
        if counterparty is not None:
            pm_body["counterparty"] = dict(counterparty)
        response = self._transport.request(
            "PATCH",
            f"{self._PREFIX}/{payment_method_id}",
            body={"payment_method": pm_body},
            idempotency_key=idempotency_key,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return PaymentMethod.from_dict(body)

    def delete(
        self,
        payment_method_id: str,
        *,
        idempotency_key: Optional[str] = None,
    ) -> None:
        """DELETE /v2/payment_methods/{id}. Returns ``None`` on 204."""
        self._transport.request(
            "DELETE",
            f"{self._PREFIX}/{payment_method_id}",
            idempotency_key=idempotency_key,
        )


__all__ = ["PaymentMethodsResource"]
