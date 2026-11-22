# tesote/sdk — packages/ruby

> Workspace brain — established by the developerz.ai maintainer agent for the maintainers of this repository from this workspace, then owned by
> its maintainers. Correct it in place; setup never overwrites a brain that exists.

Ruby gem SDK, independently tested and published to RubyGems.

The gate is the repo root’s: `bin/check` (see the root `CLAUDE.md`).

## Structure

- `lib` is the gem code and `spec` is the RSpec suite. `packages/ruby/tesote-sdk.gemspec:22-23`, `packages/ruby/Rakefile:5-8`
- The Gemfile resolves the gemspec, keeping development dependencies in the gemspec. `packages/ruby/Gemfile:1-3`, `packages/ruby/tesote-sdk.gemspec:25-30`

## Conventions

- Support Ruby 3.0 and newer. `packages/ruby/tesote-sdk.gemspec:14-14`
- Use explicit rescues. Do not use safe navigation. `packages/ruby/.rubocop.yml:45-54`
- Run `bundle install`, `bundle exec rubocop`, and `bundle exec rspec`. `.github/workflows/ruby.yml:35-44`
- Use WebMock-backed transport specs. `packages/ruby/spec/transport_spec.rb:21-33`
