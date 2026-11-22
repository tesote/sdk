# tesote/sdk — packages/python

> Workspace brain — established by the developerz.ai maintainer agent for the maintainers of this repository from this workspace, then owned by
> its maintainers. Correct it in place; setup never overwrites a brain that exists.

Python SDK package, independently tested and published to PyPI.

The gate is the repo root’s: `bin/check` (see the root `CLAUDE.md`).

## Structure

- `src/tesote_sdk` is the wheel package; `tests` is the pytest test path. `packages/python/pyproject.toml:38-42`, `packages/python/pyproject.toml:74-76`
- Test tooling is the `test` optional dependency group. `packages/python/pyproject.toml:29-30`

## Conventions

- Support Python 3.9 and newer. `packages/python/pyproject.toml:5-10`
- Keep mypy strict for `src`; use Ruff at 100 columns with the configured rule set. `packages/python/pyproject.toml:44-70`
- Run `python -m pip install --upgrade pip`, `pip install -e .[test]`, `ruff check .`, `mypy src`, and `pytest`. `.github/workflows/python.yml:40-49`
- Mock the injected opener in transport tests. `packages/python/tests/test_transport.py:1-5`
