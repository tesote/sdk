"""v2 resources not yet wired beyond their public shape."""

from __future__ import annotations

from typing import Any, Dict

from ..transport import Transport


class _StubResource:
    def __init__(self, transport: Transport) -> None:
        self._transport = transport


class TransactionsResource(_StubResource):
    def list_for_account(self, account_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 transactions.list_for_account not yet implemented")

    def get(self, transaction_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 transactions.get not yet implemented")

    def export(self, **kwargs: Any) -> Dict[str, Any]:
        raise NotImplementedError("v2 transactions.export not yet implemented")

    def sync(self, **kwargs: Any) -> Dict[str, Any]:
        raise NotImplementedError("v2 transactions.sync not yet implemented")

    def bulk(self, **kwargs: Any) -> Dict[str, Any]:
        raise NotImplementedError("v2 transactions.bulk not yet implemented")

    def search(self, **kwargs: Any) -> Dict[str, Any]:
        raise NotImplementedError("v2 transactions.search not yet implemented")


class SyncSessionsResource(_StubResource):
    def list(self, account_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 sync_sessions.list not yet implemented")

    def get(self, account_id: str, session_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 sync_sessions.get not yet implemented")


class TransactionOrdersResource(_StubResource):
    def list(self, account_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 transaction_orders.list not yet implemented")

    def get(self, account_id: str, order_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 transaction_orders.get not yet implemented")

    def create(self, account_id: str, **kwargs: Any) -> Dict[str, Any]:
        raise NotImplementedError("v2 transaction_orders.create not yet implemented")

    def submit(self, account_id: str, order_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 transaction_orders.submit not yet implemented")

    def cancel(self, account_id: str, order_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 transaction_orders.cancel not yet implemented")


class BatchesResource(_StubResource):
    def create(self, **kwargs: Any) -> Dict[str, Any]:
        raise NotImplementedError("v2 batches.create not yet implemented")

    def get(self, batch_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 batches.get not yet implemented")

    def approve(self, batch_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 batches.approve not yet implemented")

    def submit(self, batch_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 batches.submit not yet implemented")

    def cancel(self, batch_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 batches.cancel not yet implemented")


class PaymentMethodsResource(_StubResource):
    def list(self) -> Dict[str, Any]:
        raise NotImplementedError("v2 payment_methods.list not yet implemented")

    def get(self, payment_method_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 payment_methods.get not yet implemented")

    def create(self, **kwargs: Any) -> Dict[str, Any]:
        raise NotImplementedError("v2 payment_methods.create not yet implemented")

    def update(self, payment_method_id: str, **kwargs: Any) -> Dict[str, Any]:
        raise NotImplementedError("v2 payment_methods.update not yet implemented")

    def delete(self, payment_method_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v2 payment_methods.delete not yet implemented")


class StatusResource(_StubResource):
    def status(self) -> Dict[str, Any]:
        raise NotImplementedError("v2 status not yet implemented")

    def whoami(self) -> Dict[str, Any]:
        raise NotImplementedError("v2 whoami not yet implemented")


__all__ = [
    "TransactionsResource",
    "SyncSessionsResource",
    "TransactionOrdersResource",
    "BatchesResource",
    "PaymentMethodsResource",
    "StatusResource",
]
