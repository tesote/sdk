# tesote/sdk (PHP)

Official PHP client SDK for the [equipo.tesote.com](https://equipo.tesote.com) API.

Status: 0.1.0 — unreleased. v3 `accounts.list` / `accounts.get` are wired; everything else is stubbed and throws `LogicException`.

## Install

```bash
composer require tesote/sdk
```

Requires PHP 8.1+, `ext-curl`, `ext-json`. No other runtime dependencies.

## Quick start

```php
<?php
use Tesote\Sdk\V3\Client;

$client = new Client([
    'apiKey' => getenv('TESOTE_API_KEY'),
]);

$accounts = $client->accounts->list(['limit' => 50]);
foreach ($accounts['data'] as $account) {
    echo $account['id'], "\n";
}
```

## Versioned clients

```php
use Tesote\Sdk\V1\Client as V1Client;
use Tesote\Sdk\V2\Client as V2Client;
use Tesote\Sdk\V3\Client as V3Client;
```

Pick a version explicitly. `V1` / `V2` stay shipped indefinitely.

## Configuration

```php
new V3Client([
    'apiKey'           => '...',                           // required
    'baseUrl'          => 'https://equipo.tesote.com/api', // default
    'userAgent'        => 'tesote-sdk-php/0.1.0 (php/8.x)', // override for Odoo/SAP connectors
    'maxAttempts'      => 3,                               // retries on 429/5xx + transient network
    'baseDelayMs'      => 250,
    'maxDelayMs'       => 8000,
    'connectTimeoutMs' => 5000,
    'timeoutMs'        => 30000,
    'cache'            => new Tesote\Sdk\Cache\InMemoryCache(), // optional opt-in TTL cache
]);
```

## Errors

Every error is a typed subclass of `Tesote\Sdk\Errors\TesoteException`. Catch the narrowest type:

```php
try {
    $client->accounts->get('acct_x');
} catch (Tesote\Sdk\Errors\RateLimitExceededException $e) {
    sleep($e->retryAfter ?? 5);
} catch (Tesote\Sdk\Errors\UnauthorizedException $e) {
    // rotate key
}
```

Every exception carries `errorCode`, `httpStatus`, `requestId`, `errorId`, `retryAfter`, `responseBody`, `requestSummary`, `attempts`. The bearer token is always redacted to `Bearer <last4>`.

## Polling model

The platform is poll-based for v1/v2. Don't tight-loop — the SDK will surface `RateLimitExceededException` after retries. v3 adds webhooks (signature-verification helper coming soon).

## Tests

```bash
composer install
composer run-script phpstan
vendor/bin/phpunit
```

## Docs

End-user docs live at https://www.tesote.com/docs/sdk.
