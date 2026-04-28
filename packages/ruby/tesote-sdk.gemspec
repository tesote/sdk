require_relative 'lib/tesote_sdk/version'

Gem::Specification.new do |spec|
  spec.name = 'tesote-sdk'
  spec.version = TesoteSdk::VERSION
  spec.authors = ['Tesote']
  spec.email = ['support@tesote.com']

  spec.summary = 'Official Ruby SDK for the equipo.tesote.com API.'
  spec.description = 'Versioned (v1/v2) clients for the Tesote API. Zero runtime dependencies; built on Ruby stdlib net/http.'
  spec.homepage = 'https://www.tesote.com/docs/sdk/ruby'
  spec.license = 'MIT'

  spec.required_ruby_version = '>= 3.0'

  spec.metadata['homepage_uri'] = spec.homepage
  spec.metadata['source_code_uri'] = 'https://github.com/tesote/sdk/tree/main/packages/ruby'
  spec.metadata['changelog_uri'] = 'https://github.com/tesote/sdk/blob/main/packages/ruby/CHANGELOG.md'
  spec.metadata['documentation_uri'] = 'https://www.tesote.com/docs/sdk/ruby'
  spec.metadata['rubygems_mfa_required'] = 'true'

  spec.files = Dir.glob('lib/**/*.rb') + ['README.md', 'CHANGELOG.md', 'tesote-sdk.gemspec']
  spec.require_paths = ['lib']

  # why: zero-runtime-dep policy — see CLAUDE.md and docs/architecture/transport.md
  # No spec.add_dependency lines, intentionally.

  spec.add_development_dependency 'rspec'
  spec.add_development_dependency 'rubocop'
  spec.add_development_dependency 'webmock'
end
