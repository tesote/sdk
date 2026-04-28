# tesote-sdk (Ruby)

Official Ruby SDK for the [equipo.tesote.com](https://equipo.tesote.com) API.

- Zero runtime dependencies (stdlib `net/http` + `json` only).
- Ruby >= 3.0.
- Versioned clients side-by-side: `V1::Client`, `V2::Client`.
- Typed errors per `error_code`.
- Transport owns retries, rate-limit awareness, idempotency keys, opt-in TTL caching.

End-user docs: <https://www.tesote.com/docs/ruby> (canonical).

## Install

```ruby
# Gemfile
gem 'tesote-sdk'
```

## Usage

```ruby
require 'tesote_sdk'

client = TesoteSdk::V2::Client.new(api_key: ENV.fetch('TESOTE_API_KEY'))

accounts = client.accounts.list
account  = client.accounts.get('acct_123')

# Rate-limit + request-id of the most recent response
client.last_rate_limit  # => #<struct TesoteSdk::RateLimitInfo limit=200, remaining=197, reset=...>
client.last_request_id  # => "req_..."
```

### Configuration

```ruby
TesoteSdk::V2::Client.new(
  api_key: ENV.fetch('TESOTE_API_KEY'),
  base_url: 'https://equipo-staging.tesote.com/api',
  user_agent: 'MyApp/1.2.3',
  open_timeout: 5,
  read_timeout: 30,
  max_attempts: 3,
  base_delay: 0.250,
  max_delay: 8.0,
  cache_backend: TesoteSdk::CacheBackend.new(max_size: 512),
  logger: ->(event, payload) { puts({ event: event, **payload }.to_json) }
)
```

### Errors

Every SDK error carries: `error_code`, `message`, `http_status`, `request_id`,
`error_id`, `retry_after`, `response_body`, `request_summary`, `attempts`.

```ruby
begin
  client.accounts.get('acct_123')
rescue TesoteSdk::RateLimitExceededError => err
  sleep err.retry_after.to_i
  retry
rescue TesoteSdk::UnauthorizedError => err
  warn "rotate api key — request_id=#{err.request_id}"
end
```

## Polling model

v1 and v2 are poll-based. See the [Implementation Checklist on tesote.com](https://www.tesote.com/docs/ruby).

## Development

```bash
cd packages/ruby
bundle install
bundle exec rubocop
bundle exec rspec
```
