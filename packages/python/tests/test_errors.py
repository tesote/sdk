"""Errors module: dispatch table, required fields, bearer-redaction helper."""

from __future__ import annotations

import pytest

from tesote_sdk.errors import (
    AccountDisabledError,
    ApiError,
    ApiKeyRevokedError,
    HistorySyncForbiddenError,
    InvalidDateRangeError,
    MutationDuringPaginationError,
    RateLimitExceededError,
    ServiceUnavailableError,
    TesoteError,
    UnauthorizedError,
    UnprocessableContentError,
    WorkspaceSuspendedError,
    classify_api_error,
)


@pytest.mark.parametrize(
    "code,cls",
    [
        ("UNAUTHORIZED", UnauthorizedError),
        ("API_KEY_REVOKED", ApiKeyRevokedError),
        ("WORKSPACE_SUSPENDED", WorkspaceSuspendedError),
        ("ACCOUNT_DISABLED", AccountDisabledError),
        ("HISTORY_SYNC_FORBIDDEN", HistorySyncForbiddenError),
        ("MUTATION_CONFLICT", MutationDuringPaginationError),
        ("UNPROCESSABLE_CONTENT", UnprocessableContentError),
        ("INVALID_DATE_RANGE", InvalidDateRangeError),
        ("RATE_LIMIT_EXCEEDED", RateLimitExceededError),
    ],
)
def test_each_error_code_maps_to_typed_class(code: str, cls: type) -> None:
    assert classify_api_error(code, http_status=400) is cls


def test_unknown_code_falls_back_to_http_status() -> None:
    assert classify_api_error("WHO_KNOWS", 503) is ServiceUnavailableError
    assert classify_api_error(None, 422) is UnprocessableContentError


def test_unknown_everything_falls_back_to_apierror() -> None:
    assert classify_api_error(None, 418) is ApiError


def test_required_fields_present_on_every_instance() -> None:
    err = UnauthorizedError(
        "bad key",
        error_code="UNAUTHORIZED",
        http_status=401,
        request_id="rid-1",
        error_id="eid-1",
        retry_after=10,
        response_body='{"x":1}',
        request_summary={"method": "GET", "path": "/v3/accounts"},
        attempts=2,
    )
    for attr in (
        "error_code",
        "message",
        "http_status",
        "request_id",
        "error_id",
        "retry_after",
        "response_body",
        "request_summary",
        "attempts",
    ):
        assert hasattr(err, attr)
    assert err.attempts == 2
    assert isinstance(err, ApiError)
    assert isinstance(err, TesoteError)


def test_bearer_token_redaction_helper() -> None:
    from tesote_sdk.transport import _redact_bearer

    assert _redact_bearer("sk_test_abcdef1234").endswith("1234")
    assert "abcdef" not in _redact_bearer("sk_test_abcdef1234")
    assert _redact_bearer("abc") == "Bearer ****"


def test_cause_is_preserved_when_chained() -> None:
    underlying = RuntimeError("dns fail")
    try:
        try:
            raise underlying
        except RuntimeError as exc:
            raise UnauthorizedError("wrapped", error_code="UNAUTHORIZED") from exc
    except UnauthorizedError as e:
        assert e.__cause__ is underlying
