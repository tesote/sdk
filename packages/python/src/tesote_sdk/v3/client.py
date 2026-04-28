"""v3 client -- adds reporting, configuration, webhooks."""

from __future__ import annotations

from typing import Optional

from .._base_client import build_transport
from ..transport import (
    DEFAULT_BASE_URL,
    DEFAULT_CONNECT_TIMEOUT,
    DEFAULT_READ_TIMEOUT,
    CacheBackend,
    LoggerCallback,
    RetryPolicy,
    Transport,
)
from ._stubs import (
    BalanceHistoryResource,
    BatchesResource,
    CategoriesResource,
    ConnectionsResource,
    CounterpartiesResource,
    LegalEntitiesResource,
    McpResource,
    PaymentMethodsResource,
    ReportsResource,
    StatusResource,
    SyncSessionsResource,
    TransactionOrdersResource,
    TransactionsResource,
    WorkspaceResource,
)
from .accounts import AccountsResource
from .webhooks import WebhooksResource


class V3Client:
    """Public v3 client. See ``docs/architecture/resources.md``."""

    def __init__(
        self,
        api_key: str,
        *,
        base_url: str = DEFAULT_BASE_URL,
        connect_timeout: float = DEFAULT_CONNECT_TIMEOUT,
        read_timeout: float = DEFAULT_READ_TIMEOUT,
        retry_policy: Optional[RetryPolicy] = None,
        cache_backend: Optional[CacheBackend] = None,
        user_agent: Optional[str] = None,
        logger: Optional[LoggerCallback] = None,
    ) -> None:
        self._transport: Transport = build_transport(
            api_key,
            base_url=base_url,
            connect_timeout=connect_timeout,
            read_timeout=read_timeout,
            retry_policy=retry_policy,
            cache_backend=cache_backend,
            user_agent=user_agent,
            logger=logger,
        )
        self.accounts = AccountsResource(self._transport)
        self.transactions = TransactionsResource(self._transport)
        self.sync_sessions = SyncSessionsResource(self._transport)
        self.transaction_orders = TransactionOrdersResource(self._transport)
        self.batches = BatchesResource(self._transport)
        self.payment_methods = PaymentMethodsResource(self._transport)
        self.categories = CategoriesResource(self._transport)
        self.counterparties = CounterpartiesResource(self._transport)
        self.legal_entities = LegalEntitiesResource(self._transport)
        self.connections = ConnectionsResource(self._transport)
        self.webhooks = WebhooksResource(self._transport)
        self.reports = ReportsResource(self._transport)
        self.balance_history = BalanceHistoryResource(self._transport)
        self.workspace = WorkspaceResource(self._transport)
        self.mcp = McpResource(self._transport)
        self.status_resource = StatusResource(self._transport)

    @property
    def transport(self) -> Transport:
        return self._transport

    @property
    def last_rate_limit(self) -> object:
        return self._transport.last_rate_limit


__all__ = ["V3Client"]
