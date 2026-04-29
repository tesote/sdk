"""v2 status / whoami resource."""

from __future__ import annotations

from tesote_sdk.v2.status import StatusResource

from ._helpers import make_transport
from .conftest import ok_response


def test_v2_status_uses_v2_path() -> None:
    transport, opener = make_transport(
        [ok_response({"status": "ok", "authenticated": False})]
    )
    res = StatusResource(transport).status()
    assert res.status == "ok"
    assert opener.calls[0]["url"].endswith("/v2/status")


def test_v2_whoami_uses_v2_path() -> None:
    transport, opener = make_transport(
        [ok_response({"client": {"id": "cli", "name": "x", "type": "user"}})]
    )
    StatusResource(transport).whoami()
    assert opener.calls[0]["url"].endswith("/v2/whoami")
