"""Shared helpers for resource tests."""

from __future__ import annotations

from typing import Any, Sequence

from tesote_sdk.transport import RetryPolicy, Transport

from .conftest import ScriptedOpener


def make_transport(responses: Sequence[Any]) -> tuple[Transport, ScriptedOpener]:
    opener = ScriptedOpener(list(responses))
    transport = Transport(
        api_key="sk_test_abcdef1234",
        retry_policy=RetryPolicy(max_attempts=3, base_delay=0.0, max_delay=0.0),
        opener=opener,
    )
    return transport, opener
