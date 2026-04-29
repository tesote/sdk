"""v1 status / whoami resource."""

from __future__ import annotations

import pytest

from tesote_sdk.errors import UnauthorizedError
from tesote_sdk.v1.status import StatusResource

from ._helpers import make_transport
from .conftest import http_error, ok_response


def test_status_returns_typed_response() -> None:
    transport, opener = make_transport([ok_response({"status": "ok", "authenticated": False})])
    res = StatusResource(transport).status()
    assert res.status == "ok"
    assert res.authenticated is False
    assert opener.calls[0]["url"].endswith("/status")
    assert opener.calls[0]["method"] == "GET"


def test_whoami_returns_client_block() -> None:
    transport, _ = make_transport(
        [ok_response({"client": {"id": "cli_1", "name": "Acme", "type": "workspace"}})]
    )
    res = StatusResource(transport).whoami()
    assert res.client is not None
    assert res.client.id == "cli_1"
    assert res.client.name == "Acme"
    assert res.client.type == "workspace"


def test_whoami_unauthorized_raises_typed_error() -> None:
    transport, _ = make_transport(
        [http_error(401, {}, {"error": "bad", "error_code": "UNAUTHORIZED"})]
    )
    with pytest.raises(UnauthorizedError):
        StatusResource(transport).whoami()
