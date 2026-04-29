"""Typed dataclass models for v1+v2 response bodies.

All models are immutable (`@dataclass(frozen=True)`). Field names mirror the
wire snake_case exactly. ``from_dict`` factories accept partial server payloads:
unknown keys are dropped, missing keys default to ``None`` / empty containers.

Stdlib only -- no runtime deps. Built for Python 3.9 (no PEP 604 unions, no
PEP 585 builtin generics in annotations).
"""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any, Dict, List, Mapping, Optional


def _as_str(value: Any) -> Optional[str]:
    if value is None:
        return None
    if isinstance(value, str):
        return value
    return str(value)


def _as_int(value: Any) -> Optional[int]:
    if value is None:
        return None
    if isinstance(value, bool):
        # why: bool is an int in Python; we don't want True -> 1 silently
        return None
    if isinstance(value, int):
        return value
    try:
        return int(value)
    except (TypeError, ValueError):
        return None


def _as_float(value: Any) -> Optional[float]:
    if value is None:
        return None
    if isinstance(value, bool):
        return None
    if isinstance(value, (int, float)):
        return float(value)
    try:
        return float(value)
    except (TypeError, ValueError):
        return None


def _as_bool(value: Any) -> Optional[bool]:
    if value is None:
        return None
    if isinstance(value, bool):
        return value
    if isinstance(value, str):
        if value.lower() in {"true", "1", "yes"}:
            return True
        if value.lower() in {"false", "0", "no"}:
            return False
    return None


def _as_dict(value: Any) -> Dict[str, Any]:
    if isinstance(value, dict):
        return dict(value)
    return {}


def _as_list(value: Any) -> List[Any]:
    if isinstance(value, list):
        return list(value)
    return []


# Account ---------------------------------------------------------------------


@dataclass(frozen=True)
class AccountBank:
    name: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> AccountBank:
        return cls(name=_as_str(data.get("name")))


@dataclass(frozen=True)
class AccountLegalEntity:
    id: Optional[str]
    legal_name: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> AccountLegalEntity:
        return cls(
            id=_as_str(data.get("id")),
            legal_name=_as_str(data.get("legal_name")),
        )


@dataclass(frozen=True)
class AccountData:
    masked_account_number: Optional[str]
    currency: Optional[str]
    transactions_data_current_as_of: Optional[str]
    balance_data_current_as_of: Optional[str]
    custom_user_provided_identifier: Optional[str]
    balance_cents: Optional[str]
    available_balance_cents: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> AccountData:
        return cls(
            masked_account_number=_as_str(data.get("masked_account_number")),
            currency=_as_str(data.get("currency")),
            transactions_data_current_as_of=_as_str(data.get("transactions_data_current_as_of")),
            balance_data_current_as_of=_as_str(data.get("balance_data_current_as_of")),
            custom_user_provided_identifier=_as_str(data.get("custom_user_provided_identifier")),
            balance_cents=_as_str(data.get("balance_cents")),
            available_balance_cents=_as_str(data.get("available_balance_cents")),
        )


@dataclass(frozen=True)
class Account:
    id: Optional[str]
    name: Optional[str]
    data: AccountData
    bank: AccountBank
    legal_entity: AccountLegalEntity
    tesote_created_at: Optional[str]
    tesote_updated_at: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> Account:
        return cls(
            id=_as_str(data.get("id")),
            name=_as_str(data.get("name")),
            data=AccountData.from_dict(_as_dict(data.get("data"))),
            bank=AccountBank.from_dict(_as_dict(data.get("bank"))),
            legal_entity=AccountLegalEntity.from_dict(_as_dict(data.get("legal_entity"))),
            tesote_created_at=_as_str(data.get("tesote_created_at")),
            tesote_updated_at=_as_str(data.get("tesote_updated_at")),
        )


# Transaction (v1 schema) -----------------------------------------------------


@dataclass(frozen=True)
class TransactionCategory:
    name: Optional[str]
    external_category_code: Optional[str]
    created_at: Optional[str]
    updated_at: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> TransactionCategory:
        return cls(
            name=_as_str(data.get("name")),
            external_category_code=_as_str(data.get("external_category_code")),
            created_at=_as_str(data.get("created_at")),
            updated_at=_as_str(data.get("updated_at")),
        )


@dataclass(frozen=True)
class Counterparty:
    name: Optional[str]
    id: Optional[str] = None

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> Counterparty:
        return cls(
            name=_as_str(data.get("name")),
            id=_as_str(data.get("id")),
        )


@dataclass(frozen=True)
class TransactionData:
    amount_cents: Optional[int]
    currency: Optional[str]
    description: Optional[str]
    transaction_date: Optional[str]
    created_at: Optional[str]
    created_at_date: Optional[str]
    note: Optional[str]
    external_service_id: Optional[str]
    running_balance_cents: Optional[int]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> TransactionData:
        return cls(
            amount_cents=_as_int(data.get("amount_cents")),
            currency=_as_str(data.get("currency")),
            description=_as_str(data.get("description")),
            transaction_date=_as_str(data.get("transaction_date")),
            created_at=_as_str(data.get("created_at")),
            created_at_date=_as_str(data.get("created_at_date")),
            note=_as_str(data.get("note")),
            external_service_id=_as_str(data.get("external_service_id")),
            running_balance_cents=_as_int(data.get("running_balance_cents")),
        )


@dataclass(frozen=True)
class Transaction:
    id: Optional[str]
    status: Optional[str]
    data: TransactionData
    tesote_imported_at: Optional[str]
    tesote_updated_at: Optional[str]
    transaction_categories: List[TransactionCategory] = field(default_factory=list)
    counterparty: Optional[Counterparty] = None

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> Transaction:
        cats_raw = _as_list(data.get("transaction_categories"))
        categories = [TransactionCategory.from_dict(c) for c in cats_raw if isinstance(c, dict)]
        cp_raw = data.get("counterparty")
        counterparty = (
            Counterparty.from_dict(cp_raw) if isinstance(cp_raw, dict) else None
        )
        return cls(
            id=_as_str(data.get("id")),
            status=_as_str(data.get("status")),
            data=TransactionData.from_dict(_as_dict(data.get("data"))),
            tesote_imported_at=_as_str(data.get("tesote_imported_at")),
            tesote_updated_at=_as_str(data.get("tesote_updated_at")),
            transaction_categories=categories,
            counterparty=counterparty,
        )


# SyncTransaction (v2 sync) ---------------------------------------------------


@dataclass(frozen=True)
class SyncTransaction:
    transaction_id: Optional[str]
    account_id: Optional[str]
    amount: Optional[float]
    iso_currency_code: Optional[str]
    unofficial_currency_code: Optional[str]
    date: Optional[str]
    datetime: Optional[str]
    name: Optional[str]
    merchant_name: Optional[str]
    pending: Optional[bool]
    category: List[str] = field(default_factory=list)
    running_balance_cents: Optional[int] = None

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> SyncTransaction:
        cats = [str(c) for c in _as_list(data.get("category")) if isinstance(c, str)]
        return cls(
            transaction_id=_as_str(data.get("transaction_id")),
            account_id=_as_str(data.get("account_id")),
            amount=_as_float(data.get("amount")),
            iso_currency_code=_as_str(data.get("iso_currency_code")),
            unofficial_currency_code=_as_str(data.get("unofficial_currency_code")),
            date=_as_str(data.get("date")),
            datetime=_as_str(data.get("datetime")),
            name=_as_str(data.get("name")),
            merchant_name=_as_str(data.get("merchant_name")),
            pending=_as_bool(data.get("pending")),
            category=cats,
            running_balance_cents=_as_int(data.get("running_balance_cents")),
        )


@dataclass(frozen=True)
class RemovedSyncTransaction:
    transaction_id: Optional[str]
    account_id: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> RemovedSyncTransaction:
        return cls(
            transaction_id=_as_str(data.get("transaction_id")),
            account_id=_as_str(data.get("account_id")),
        )


@dataclass(frozen=True)
class SyncDelta:
    added: List[SyncTransaction]
    modified: List[SyncTransaction]
    removed: List[RemovedSyncTransaction]
    next_cursor: Optional[str]
    has_more: bool

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> SyncDelta:
        added = [
            SyncTransaction.from_dict(item)
            for item in _as_list(data.get("added"))
            if isinstance(item, dict)
        ]
        modified = [
            SyncTransaction.from_dict(item)
            for item in _as_list(data.get("modified"))
            if isinstance(item, dict)
        ]
        removed = [
            RemovedSyncTransaction.from_dict(item)
            for item in _as_list(data.get("removed"))
            if isinstance(item, dict)
        ]
        has_more_value = _as_bool(data.get("has_more"))
        return cls(
            added=added,
            modified=modified,
            removed=removed,
            next_cursor=_as_str(data.get("next_cursor")),
            has_more=bool(has_more_value) if has_more_value is not None else False,
        )


# TransactionOrder ------------------------------------------------------------


@dataclass(frozen=True)
class OrderSourceAccount:
    id: Optional[str]
    name: Optional[str]
    payment_method_id: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> OrderSourceAccount:
        return cls(
            id=_as_str(data.get("id")),
            name=_as_str(data.get("name")),
            payment_method_id=_as_str(data.get("payment_method_id")),
        )


@dataclass(frozen=True)
class OrderDestination:
    payment_method_id: Optional[str]
    counterparty_id: Optional[str]
    counterparty_name: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> OrderDestination:
        return cls(
            payment_method_id=_as_str(data.get("payment_method_id")),
            counterparty_id=_as_str(data.get("counterparty_id")),
            counterparty_name=_as_str(data.get("counterparty_name")),
        )


@dataclass(frozen=True)
class Fee:
    amount: Optional[float]
    currency: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> Fee:
        return cls(
            amount=_as_float(data.get("amount")),
            currency=_as_str(data.get("currency")),
        )


@dataclass(frozen=True)
class TesoteTransactionRef:
    id: Optional[str]
    status: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> TesoteTransactionRef:
        return cls(
            id=_as_str(data.get("id")),
            status=_as_str(data.get("status")),
        )


@dataclass(frozen=True)
class TransactionOrderAttempt:
    id: Optional[str]
    status: Optional[str]
    attempt_number: Optional[int]
    external_reference: Optional[str]
    submitted_at: Optional[str]
    completed_at: Optional[str]
    error_code: Optional[str]
    error_message: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> TransactionOrderAttempt:
        return cls(
            id=_as_str(data.get("id")),
            status=_as_str(data.get("status")),
            attempt_number=_as_int(data.get("attempt_number")),
            external_reference=_as_str(data.get("external_reference")),
            submitted_at=_as_str(data.get("submitted_at")),
            completed_at=_as_str(data.get("completed_at")),
            error_code=_as_str(data.get("error_code")),
            error_message=_as_str(data.get("error_message")),
        )


@dataclass(frozen=True)
class TransactionOrder:
    id: Optional[str]
    status: Optional[str]
    amount: Optional[float]
    currency: Optional[str]
    description: Optional[str]
    reference: Optional[str]
    external_reference: Optional[str]
    idempotency_key: Optional[str]
    batch_id: Optional[str]
    scheduled_for: Optional[str]
    approved_at: Optional[str]
    submitted_at: Optional[str]
    completed_at: Optional[str]
    failed_at: Optional[str]
    cancelled_at: Optional[str]
    source_account: Optional[OrderSourceAccount]
    destination: Optional[OrderDestination]
    fee: Optional[Fee]
    execution_strategy: Optional[str]
    tesote_transaction: Optional[TesoteTransactionRef]
    latest_attempt: Optional[TransactionOrderAttempt]
    created_at: Optional[str]
    updated_at: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> TransactionOrder:
        src = data.get("source_account")
        dst = data.get("destination")
        fee = data.get("fee")
        ttx = data.get("tesote_transaction")
        latest = data.get("latest_attempt")
        return cls(
            id=_as_str(data.get("id")),
            status=_as_str(data.get("status")),
            amount=_as_float(data.get("amount")),
            currency=_as_str(data.get("currency")),
            description=_as_str(data.get("description")),
            reference=_as_str(data.get("reference")),
            external_reference=_as_str(data.get("external_reference")),
            idempotency_key=_as_str(data.get("idempotency_key")),
            batch_id=_as_str(data.get("batch_id")),
            scheduled_for=_as_str(data.get("scheduled_for")),
            approved_at=_as_str(data.get("approved_at")),
            submitted_at=_as_str(data.get("submitted_at")),
            completed_at=_as_str(data.get("completed_at")),
            failed_at=_as_str(data.get("failed_at")),
            cancelled_at=_as_str(data.get("cancelled_at")),
            source_account=OrderSourceAccount.from_dict(src) if isinstance(src, dict) else None,
            destination=OrderDestination.from_dict(dst) if isinstance(dst, dict) else None,
            fee=Fee.from_dict(fee) if isinstance(fee, dict) else None,
            execution_strategy=_as_str(data.get("execution_strategy")),
            tesote_transaction=(
                TesoteTransactionRef.from_dict(ttx) if isinstance(ttx, dict) else None
            ),
            latest_attempt=(
                TransactionOrderAttempt.from_dict(latest) if isinstance(latest, dict) else None
            ),
            created_at=_as_str(data.get("created_at")),
            updated_at=_as_str(data.get("updated_at")),
        )


# PaymentMethod ---------------------------------------------------------------


@dataclass(frozen=True)
class PaymentMethodAccountRef:
    id: Optional[str]
    name: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> PaymentMethodAccountRef:
        return cls(
            id=_as_str(data.get("id")),
            name=_as_str(data.get("name")),
        )


@dataclass(frozen=True)
class PaymentMethod:
    id: Optional[str]
    method_type: Optional[str]
    currency: Optional[str]
    label: Optional[str]
    details: Dict[str, Any]
    verified: Optional[bool]
    verified_at: Optional[str]
    last_used_at: Optional[str]
    counterparty: Optional[Counterparty]
    tesote_account: Optional[PaymentMethodAccountRef]
    created_at: Optional[str]
    updated_at: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> PaymentMethod:
        cp = data.get("counterparty")
        acc = data.get("tesote_account")
        return cls(
            id=_as_str(data.get("id")),
            method_type=_as_str(data.get("method_type")),
            currency=_as_str(data.get("currency")),
            label=_as_str(data.get("label")),
            details=_as_dict(data.get("details")),
            verified=_as_bool(data.get("verified")),
            verified_at=_as_str(data.get("verified_at")),
            last_used_at=_as_str(data.get("last_used_at")),
            counterparty=Counterparty.from_dict(cp) if isinstance(cp, dict) else None,
            tesote_account=(
                PaymentMethodAccountRef.from_dict(acc) if isinstance(acc, dict) else None
            ),
            created_at=_as_str(data.get("created_at")),
            updated_at=_as_str(data.get("updated_at")),
        )


# SyncSession -----------------------------------------------------------------


@dataclass(frozen=True)
class SyncSessionError:
    type: Optional[str]
    message: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> SyncSessionError:
        return cls(
            type=_as_str(data.get("type")),
            message=_as_str(data.get("message")),
        )


@dataclass(frozen=True)
class SyncSessionPerformance:
    total_duration: Optional[float]
    complexity_score: Optional[float]
    sync_speed_score: Optional[float]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> SyncSessionPerformance:
        return cls(
            total_duration=_as_float(data.get("total_duration")),
            complexity_score=_as_float(data.get("complexity_score")),
            sync_speed_score=_as_float(data.get("sync_speed_score")),
        )


@dataclass(frozen=True)
class SyncSession:
    id: Optional[str]
    status: Optional[str]
    started_at: Optional[str]
    completed_at: Optional[str]
    transactions_synced: Optional[int]
    accounts_count: Optional[int]
    error: Optional[SyncSessionError]
    performance: Optional[SyncSessionPerformance]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> SyncSession:
        err = data.get("error")
        perf = data.get("performance")
        return cls(
            id=_as_str(data.get("id")),
            status=_as_str(data.get("status")),
            started_at=_as_str(data.get("started_at")),
            completed_at=_as_str(data.get("completed_at")),
            transactions_synced=_as_int(data.get("transactions_synced")),
            accounts_count=_as_int(data.get("accounts_count")),
            error=SyncSessionError.from_dict(err) if isinstance(err, dict) else None,
            performance=(
                SyncSessionPerformance.from_dict(perf) if isinstance(perf, dict) else None
            ),
        )


# Sync trigger response (POST /v2/accounts/{id}/sync) ------------------------


@dataclass(frozen=True)
class AccountSyncStarted:
    message: Optional[str]
    sync_session_id: Optional[str]
    status: Optional[str]
    started_at: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> AccountSyncStarted:
        return cls(
            message=_as_str(data.get("message")),
            sync_session_id=_as_str(data.get("sync_session_id")),
            status=_as_str(data.get("status")),
            started_at=_as_str(data.get("started_at")),
        )


# Batch summary ---------------------------------------------------------------


@dataclass(frozen=True)
class BatchSummary:
    batch_id: Optional[str]
    total_orders: Optional[int]
    total_amount_cents: Optional[int]
    amount_currency: Optional[str]
    statuses: Dict[str, int]
    batch_status: Optional[str]
    created_at: Optional[str]
    orders: List[TransactionOrder]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> BatchSummary:
        statuses_raw = _as_dict(data.get("statuses"))
        statuses: Dict[str, int] = {}
        for k, v in statuses_raw.items():
            iv = _as_int(v)
            if iv is not None:
                statuses[str(k)] = iv
        orders = [
            TransactionOrder.from_dict(o)
            for o in _as_list(data.get("orders"))
            if isinstance(o, dict)
        ]
        return cls(
            batch_id=_as_str(data.get("batch_id")),
            total_orders=_as_int(data.get("total_orders")),
            total_amount_cents=_as_int(data.get("total_amount_cents")),
            amount_currency=_as_str(data.get("amount_currency")),
            statuses=statuses,
            batch_status=_as_str(data.get("batch_status")),
            created_at=_as_str(data.get("created_at")),
            orders=orders,
        )


@dataclass(frozen=True)
class BatchCreateResult:
    batch_id: Optional[str]
    orders: List[TransactionOrder]
    errors: List[Dict[str, Any]]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> BatchCreateResult:
        orders = [
            TransactionOrder.from_dict(o)
            for o in _as_list(data.get("orders"))
            if isinstance(o, dict)
        ]
        errs = [e for e in _as_list(data.get("errors")) if isinstance(e, dict)]
        return cls(
            batch_id=_as_str(data.get("batch_id")),
            orders=orders,
            errors=errs,
        )


@dataclass(frozen=True)
class BatchApproveResult:
    approved: int
    failed: int

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> BatchApproveResult:
        return cls(
            approved=_as_int(data.get("approved")) or 0,
            failed=_as_int(data.get("failed")) or 0,
        )


@dataclass(frozen=True)
class BatchSubmitResult:
    enqueued: int
    failed: int

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> BatchSubmitResult:
        return cls(
            enqueued=_as_int(data.get("enqueued")) or 0,
            failed=_as_int(data.get("failed")) or 0,
        )


@dataclass(frozen=True)
class BatchCancelResult:
    cancelled: int
    skipped: int
    errors: List[Dict[str, Any]]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> BatchCancelResult:
        errs = [e for e in _as_list(data.get("errors")) if isinstance(e, dict)]
        return cls(
            cancelled=_as_int(data.get("cancelled")) or 0,
            skipped=_as_int(data.get("skipped")) or 0,
            errors=errs,
        )


# Status / WhoAmI -------------------------------------------------------------


@dataclass(frozen=True)
class StatusResponse:
    status: Optional[str]
    authenticated: Optional[bool]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> StatusResponse:
        return cls(
            status=_as_str(data.get("status")),
            authenticated=_as_bool(data.get("authenticated")),
        )


@dataclass(frozen=True)
class WhoAmIClient:
    id: Optional[str]
    name: Optional[str]
    type: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> WhoAmIClient:
        return cls(
            id=_as_str(data.get("id")),
            name=_as_str(data.get("name")),
            type=_as_str(data.get("type")),
        )


@dataclass(frozen=True)
class WhoAmI:
    client: Optional[WhoAmIClient]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> WhoAmI:
        cl = data.get("client")
        return cls(client=WhoAmIClient.from_dict(cl) if isinstance(cl, dict) else None)


# Pagination shapes -----------------------------------------------------------


@dataclass(frozen=True)
class PageInfo:
    current_page: Optional[int]
    per_page: Optional[int]
    total_pages: Optional[int]
    total_count: Optional[int]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> PageInfo:
        return cls(
            current_page=_as_int(data.get("current_page")),
            per_page=_as_int(data.get("per_page")),
            total_pages=_as_int(data.get("total_pages")),
            total_count=_as_int(data.get("total_count")),
        )


@dataclass(frozen=True)
class CursorInfo:
    has_more: bool
    per_page: Optional[int]
    after_id: Optional[str]
    before_id: Optional[str]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> CursorInfo:
        has_more_value = _as_bool(data.get("has_more"))
        return cls(
            has_more=bool(has_more_value) if has_more_value is not None else False,
            per_page=_as_int(data.get("per_page")),
            after_id=_as_str(data.get("after_id")),
            before_id=_as_str(data.get("before_id")),
        )


@dataclass(frozen=True)
class AccountList:
    total: Optional[int]
    accounts: List[Account]
    pagination: Optional[PageInfo]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> AccountList:
        accounts = [
            Account.from_dict(a)
            for a in _as_list(data.get("accounts"))
            if isinstance(a, dict)
        ]
        page = data.get("pagination")
        return cls(
            total=_as_int(data.get("total")),
            accounts=accounts,
            pagination=PageInfo.from_dict(page) if isinstance(page, dict) else None,
        )


@dataclass(frozen=True)
class TransactionList:
    total: Optional[int]
    transactions: List[Transaction]
    pagination: Optional[CursorInfo]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> TransactionList:
        txns = [
            Transaction.from_dict(t)
            for t in _as_list(data.get("transactions"))
            if isinstance(t, dict)
        ]
        page = data.get("pagination")
        return cls(
            total=_as_int(data.get("total")),
            transactions=txns,
            pagination=CursorInfo.from_dict(page) if isinstance(page, dict) else None,
        )


@dataclass(frozen=True)
class BulkAccountResult:
    account_id: Optional[str]
    transactions: List[Transaction]
    pagination: Optional[CursorInfo]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> BulkAccountResult:
        txns = [
            Transaction.from_dict(t)
            for t in _as_list(data.get("transactions"))
            if isinstance(t, dict)
        ]
        page = data.get("pagination")
        return cls(
            account_id=_as_str(data.get("account_id")),
            transactions=txns,
            pagination=CursorInfo.from_dict(page) if isinstance(page, dict) else None,
        )


@dataclass(frozen=True)
class BulkResult:
    bulk_results: List[BulkAccountResult]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> BulkResult:
        items = [
            BulkAccountResult.from_dict(b)
            for b in _as_list(data.get("bulk_results"))
            if isinstance(b, dict)
        ]
        return cls(bulk_results=items)


@dataclass(frozen=True)
class SearchResult:
    transactions: List[Transaction]
    total: Optional[int]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> SearchResult:
        txns = [
            Transaction.from_dict(t)
            for t in _as_list(data.get("transactions"))
            if isinstance(t, dict)
        ]
        return cls(
            transactions=txns,
            total=_as_int(data.get("total")),
        )


@dataclass(frozen=True)
class SyncSessionList:
    sync_sessions: List[SyncSession]
    limit: Optional[int]
    offset: Optional[int]
    has_more: bool

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> SyncSessionList:
        sessions = [
            SyncSession.from_dict(s)
            for s in _as_list(data.get("sync_sessions"))
            if isinstance(s, dict)
        ]
        has_more_value = _as_bool(data.get("has_more"))
        return cls(
            sync_sessions=sessions,
            limit=_as_int(data.get("limit")),
            offset=_as_int(data.get("offset")),
            has_more=bool(has_more_value) if has_more_value is not None else False,
        )


@dataclass(frozen=True)
class TransactionOrderList:
    items: List[TransactionOrder]
    has_more: bool
    limit: Optional[int]
    offset: Optional[int]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> TransactionOrderList:
        items = [
            TransactionOrder.from_dict(o)
            for o in _as_list(data.get("items"))
            if isinstance(o, dict)
        ]
        has_more_value = _as_bool(data.get("has_more"))
        return cls(
            items=items,
            has_more=bool(has_more_value) if has_more_value is not None else False,
            limit=_as_int(data.get("limit")),
            offset=_as_int(data.get("offset")),
        )


@dataclass(frozen=True)
class PaymentMethodList:
    items: List[PaymentMethod]
    has_more: bool
    limit: Optional[int]
    offset: Optional[int]

    @classmethod
    def from_dict(cls, data: Mapping[str, Any]) -> PaymentMethodList:
        items = [
            PaymentMethod.from_dict(p)
            for p in _as_list(data.get("items"))
            if isinstance(p, dict)
        ]
        has_more_value = _as_bool(data.get("has_more"))
        return cls(
            items=items,
            has_more=bool(has_more_value) if has_more_value is not None else False,
            limit=_as_int(data.get("limit")),
            offset=_as_int(data.get("offset")),
        )


__all__ = [
    "Account",
    "AccountBank",
    "AccountData",
    "AccountLegalEntity",
    "AccountList",
    "AccountSyncStarted",
    "BatchApproveResult",
    "BatchCancelResult",
    "BatchCreateResult",
    "BatchSubmitResult",
    "BatchSummary",
    "BulkAccountResult",
    "BulkResult",
    "Counterparty",
    "CursorInfo",
    "Fee",
    "OrderDestination",
    "OrderSourceAccount",
    "PageInfo",
    "PaymentMethod",
    "PaymentMethodAccountRef",
    "PaymentMethodList",
    "RemovedSyncTransaction",
    "SearchResult",
    "StatusResponse",
    "SyncDelta",
    "SyncSession",
    "SyncSessionError",
    "SyncSessionList",
    "SyncSessionPerformance",
    "SyncTransaction",
    "TesoteTransactionRef",
    "Transaction",
    "TransactionCategory",
    "TransactionData",
    "TransactionList",
    "TransactionOrder",
    "TransactionOrderAttempt",
    "TransactionOrderList",
    "WhoAmI",
    "WhoAmIClient",
]
