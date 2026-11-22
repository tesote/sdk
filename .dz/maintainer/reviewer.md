# Reviewer instructions — tesote/sdk

> Read by the developerz.ai reviewer at EVERY review of this repository, from the base
> branch's tip (an edit steers reviews once it merges), as this repo's own layer of the reviewer's prompt — established by the developerz.ai maintainer agent for the maintainers of this repository.
> It refines the platform rules (a finding is cited or withheld)
> and never overrides them. Edit freely; keep it under 12,000 characters — past that it is cut.

## Stack

Go (go) + Python (packages/python) + Ruby + Bundler (packages/ruby) + Node + npm (packages/ts). The gate is `bin/check` (`(cd 'go' && go vet ./... && go test ./...) && (cd 'packages/python' && ./.venv/bin/python -m ruff check . && { ./.venv/bin/python -m pytest || test $? -eq 5; }) && (cd 'packages/ruby' && bundle exec rubocop && bundle exec rake test) && (cd 'packages/ts' && npm run typecheck && npm run lint && npm run test)`): a change that does not pass it is not proven, whatever its description says.

## Review conventions

## Gates that matter

## Never flag
