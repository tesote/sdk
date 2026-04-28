"""Shared test helpers: a fake `urlopen` that records calls and returns scripted responses."""

from __future__ import annotations

import io
import json
import urllib.error
from email.message import Message
from typing import Any, Dict, List, Mapping, Optional, Sequence, Union


class FakeHTTPResponse:
    """Mimics the subset of `http.client.HTTPResponse` that Transport reads."""

    def __init__(self, status: int, headers: Mapping[str, str], body: bytes) -> None:
        self.status = status
        self._body = body
        self.headers = self._build_headers(headers)

    @staticmethod
    def _build_headers(headers: Mapping[str, str]) -> Message:
        msg = Message()
        for k, v in headers.items():
            msg[k] = v
        return msg

    def read(self) -> bytes:
        return self._body

    def getcode(self) -> int:
        return self.status

    def close(self) -> None:
        return None


def http_error(
    status: int, headers: Mapping[str, str], body: Union[bytes, dict]
) -> urllib.error.HTTPError:
    if isinstance(body, dict):
        body = json.dumps(body).encode("utf-8")
    msg = Message()
    for k, v in headers.items():
        msg[k] = v
    return urllib.error.HTTPError(
        url="https://example.test", code=status, msg=str(status), hdrs=msg, fp=io.BytesIO(body)
    )


def ok_response(
    body: Any, *, status: int = 200, headers: Optional[Mapping[str, str]] = None
) -> FakeHTTPResponse:
    payload = json.dumps(body).encode("utf-8") if not isinstance(body, bytes) else body
    return FakeHTTPResponse(status, headers or {}, payload)


class ScriptedOpener:
    """Calls(): list of (Request, timeout). Each call pops a scripted response."""

    def __init__(self, responses: Sequence[Any]) -> None:
        self._responses = list(responses)
        self.calls: List[Dict[str, Any]] = []

    def __call__(self, request: Any, timeout: Optional[float] = None) -> Any:
        self.calls.append(
            {
                "method": request.get_method(),
                "url": request.full_url,
                "headers": {k: v for k, v in request.header_items()},
                "body": request.data,
                "timeout": timeout,
            }
        )
        if not self._responses:
            raise AssertionError("ScriptedOpener exhausted")
        item = self._responses.pop(0)
        if isinstance(item, Exception):
            raise item
        return item
