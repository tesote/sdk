#!/usr/bin/env python3
"""Cross-language parity check.

Reads spec/parity.yaml and asserts that every SDK under packages/<lang>/
exposes (a) every resource method named in the manifest and (b) every typed
error class. Match strategy is intentionally coarse: we grep the source tree
for identifiers and flag any that never appear. False positives are preferred
over false negatives — the manifest is the source of truth.

Exit codes:
  0  every language has everything
  1  one or more languages are missing items
  2  bad invocation / spec parse failure

Run from anywhere; paths resolve relative to the repo root (parent of bin/).
"""

from __future__ import annotations

import os
import re
import sys
from dataclasses import dataclass, field
from pathlib import Path

try:
    import yaml
except ImportError:
    sys.stderr.write(
        "PyYAML is required: pip install pyyaml (or: apt install python3-yaml)\n"
    )
    sys.exit(2)


REPO_ROOT = Path(__file__).resolve().parent.parent
SPEC = REPO_ROOT / "spec" / "parity.yaml"
PACKAGES = REPO_ROOT / "packages"


@dataclass
class LanguageReport:
    name: str
    missing_methods: list[str] = field(default_factory=list)
    missing_errors: list[str] = field(default_factory=list)
    package_present: bool = True

    @property
    def ok(self) -> bool:
        return (
            self.package_present
            and not self.missing_methods
            and not self.missing_errors
        )


def to_camel(snake: str) -> str:
    head, *tail = snake.split("_")
    return head + "".join(part.capitalize() for part in tail)


def to_pascal(snake: str) -> str:
    return "".join(part.capitalize() for part in snake.split("_"))


def method_variants(method_name: str, case: str) -> list[str]:
    """Return identifier spellings to grep for, depending on language case."""
    if case == "snake":
        # The manifest already uses camelCase for method names like listForAccount.
        # Convert to snake_case for snake-case languages.
        snake = re.sub(r"([A-Z])", r"_\1", method_name).lower()
        return [snake]
    if case == "camel":
        return [method_name]
    if case == "pascal":
        return [method_name[:1].upper() + method_name[1:]]
    raise ValueError(f"unknown method_case {case!r}")


def collect_source(lang_dir: Path, extensions: list[str]) -> str:
    """Concatenate every source file's text (extensions filter, no node_modules)."""
    chunks: list[str] = []
    skip = {"node_modules", "vendor", "build", "dist", "target", ".venv", "__pycache__"}
    for root, dirs, files in os.walk(lang_dir):
        dirs[:] = [d for d in dirs if d not in skip and not d.startswith(".")]
        for fname in files:
            ext = fname.rsplit(".", 1)[-1].lower()
            if ext not in extensions:
                continue
            try:
                chunks.append(Path(root, fname).read_text(encoding="utf-8", errors="replace"))
            except OSError:
                continue
    return "\n".join(chunks)


def check_language(lang_cfg: dict, manifest: dict) -> LanguageReport:
    name = lang_cfg["dir"]
    report = LanguageReport(name=name)
    lang_dir = PACKAGES / name

    if not lang_dir.exists():
        report.package_present = False
        return report

    source = collect_source(lang_dir, lang_cfg["extensions"])
    case = lang_cfg["method_case"]
    suffix = lang_cfg["class_suffix"]
    async_suffix = lang_cfg.get("async_suffix", "")

    # --- methods ---
    for version, vcfg in manifest["versions"].items():
        for resource, rcfg in vcfg["resources"].items():
            for method in rcfg["methods"]:
                spellings = method_variants(method, case)
                if async_suffix:
                    spellings = spellings + [s + async_suffix for s in spellings]
                if not any(re.search(rf"\b{re.escape(s)}\b", source) for s in spellings):
                    report.missing_methods.append(f"{version}.{resource}.{method}")

    # --- errors ---
    for err in manifest["errors"]:
        canonical = err["class"]
        # Translate `XxxError` <-> `XxxException` per language convention.
        if suffix == "Exception" and canonical.endswith("Error"):
            target = canonical[:-len("Error")] + "Exception"
        else:
            target = canonical
        # Go has both `*XxxError` types and `ErrXxx` sentinels — accept either.
        if name == "go":
            sentinel = "Err" + canonical[:-len("Error")] if canonical.endswith("Error") else "Err" + canonical
            if not (
                re.search(rf"\b{re.escape(target)}\b", source)
                or re.search(rf"\b{re.escape(sentinel)}\b", source)
            ):
                report.missing_errors.append(err["code"])
        else:
            if not re.search(rf"\b{re.escape(target)}\b", source):
                report.missing_errors.append(err["code"])

    return report


def main() -> int:
    if not SPEC.exists():
        sys.stderr.write(f"missing manifest: {SPEC}\n")
        return 2
    try:
        manifest = yaml.safe_load(SPEC.read_text())
    except yaml.YAMLError as exc:
        sys.stderr.write(f"failed to parse {SPEC}: {exc}\n")
        return 2

    print(f"# parity check (spec: {SPEC.relative_to(REPO_ROOT)})")
    print(f"# repo:   {REPO_ROOT}")
    print()

    any_failures = False
    for lang_cfg in manifest["languages"]:
        report = check_language(lang_cfg, manifest)
        header = f"## {report.name}"
        if not report.package_present:
            print(f"{header}: SKIP (packages/{report.name}/ does not exist yet)")
            print()
            continue
        if report.ok:
            print(f"{header}: OK")
            print()
            continue

        any_failures = True
        print(f"{header}: FAIL")
        if report.missing_methods:
            print(f"  missing methods ({len(report.missing_methods)}):")
            for item in report.missing_methods:
                print(f"    - {item}")
        if report.missing_errors:
            print(f"  missing error classes ({len(report.missing_errors)}):")
            for code in report.missing_errors:
                print(f"    - {code}")
        print()

    if any_failures:
        sys.stderr.write(
            "parity check failed: at least one SDK is missing a method or error.\n"
            "Update the SDK to expose the missing identifiers, OR (only if the\n"
            "manifest is wrong) update spec/parity.yaml plus the matching\n"
            "docs/architecture/{resources,errors}.md page in the same PR.\n"
        )
        return 1
    print("parity check: all languages passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
