"""v2 batches resource: create, show, approve, submit, cancel."""

from __future__ import annotations

from typing import Any, Dict, Mapping, Optional, Sequence

from ..models import (
    BatchApproveResult,
    BatchCancelResult,
    BatchCreateResult,
    BatchSubmitResult,
    BatchSummary,
)
from ..transport import Transport


class BatchesResource:
    """`/v2/accounts/{id}/batches`."""

    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def create(
        self,
        account_id: str,
        orders: Sequence[Mapping[str, Any]],
        *,
        idempotency_key: Optional[str] = None,
    ) -> BatchCreateResult:
        """POST /v2/accounts/{id}/batches. ``orders`` mirrors the spec's order body shape."""
        body: Dict[str, Any] = {"orders": [dict(o) for o in orders]}
        response = self._transport.request(
            "POST",
            f"/v2/accounts/{account_id}/batches",
            body=body,
            idempotency_key=idempotency_key,
        )
        payload: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return BatchCreateResult.from_dict(payload)

    def get(
        self,
        account_id: str,
        batch_id: str,
        *,
        cache_ttl: Optional[float] = None,
    ) -> BatchSummary:
        response = self._transport.request(
            "GET",
            f"/v2/accounts/{account_id}/batches/{batch_id}",
            cache_ttl=cache_ttl,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return BatchSummary.from_dict(body)

    def approve(
        self,
        account_id: str,
        batch_id: str,
        *,
        idempotency_key: Optional[str] = None,
    ) -> BatchApproveResult:
        response = self._transport.request(
            "POST",
            f"/v2/accounts/{account_id}/batches/{batch_id}/approve",
            body={},
            idempotency_key=idempotency_key,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return BatchApproveResult.from_dict(body)

    def submit(
        self,
        account_id: str,
        batch_id: str,
        *,
        token: Optional[str] = None,
        idempotency_key: Optional[str] = None,
    ) -> BatchSubmitResult:
        body: Dict[str, Any] = {}
        if token is not None:
            body["token"] = token
        response = self._transport.request(
            "POST",
            f"/v2/accounts/{account_id}/batches/{batch_id}/submit",
            body=body,
            idempotency_key=idempotency_key,
        )
        payload: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return BatchSubmitResult.from_dict(payload)

    def cancel(
        self,
        account_id: str,
        batch_id: str,
        *,
        idempotency_key: Optional[str] = None,
    ) -> BatchCancelResult:
        response = self._transport.request(
            "POST",
            f"/v2/accounts/{account_id}/batches/{batch_id}/cancel",
            body={},
            idempotency_key=idempotency_key,
        )
        body: Dict[str, Any] = response.json if isinstance(response.json, dict) else {}
        return BatchCancelResult.from_dict(body)


__all__ = ["BatchesResource"]
