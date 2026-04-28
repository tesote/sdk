"""v3 webhooks resource + signature verification helper.

The signature scheme is a **stub** awaiting platform confirmation -- the helper
shape exists so callers can wire it now and the implementation lands when the
spec is final. ``verify_webhook_signature`` raises ``NotImplementedError`` for
non-empty inputs and treats the empty-string sentinel as a no-op for tests.
"""

from __future__ import annotations

from typing import Any, Dict, Union

from ..transport import Transport


class WebhooksResource:
    def __init__(self, transport: Transport) -> None:
        self._transport = transport

    def list(self) -> Dict[str, Any]:
        raise NotImplementedError("v3 webhooks.list not yet implemented")

    def get(self, webhook_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v3 webhooks.get not yet implemented")

    def create(self, **kwargs: Any) -> Dict[str, Any]:
        raise NotImplementedError("v3 webhooks.create not yet implemented")

    def update(self, webhook_id: str, **kwargs: Any) -> Dict[str, Any]:
        raise NotImplementedError("v3 webhooks.update not yet implemented")

    def delete(self, webhook_id: str) -> Dict[str, Any]:
        raise NotImplementedError("v3 webhooks.delete not yet implemented")


def verify_webhook_signature(
    *, body: Union[str, bytes], signature_header: str, secret: str
) -> None:
    """Verify a Tesote webhook signature.

    Stub: signature scheme awaits platform confirmation
    (see ``docs/architecture/resources.md``). Raises ``NotImplementedError``.
    """
    if not isinstance(body, (str, bytes)):
        raise TypeError("body must be str or bytes")
    if not isinstance(signature_header, str):
        raise TypeError("signature_header must be str")
    if not isinstance(secret, str):
        raise TypeError("secret must be str")
    raise NotImplementedError(
        "verify_webhook_signature stub -- signature scheme TBD; see resources.md"
    )


__all__ = ["WebhooksResource", "verify_webhook_signature"]
