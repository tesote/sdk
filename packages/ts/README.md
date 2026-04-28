# @tesote.com/sdk

Official TypeScript SDK for the [Tesote API](https://www.tesote.com).

Zero runtime dependencies. Native `fetch`. Node 18+.

## Install

```sh
npm install @tesote.com/sdk
```

## Quick start

```ts
import { V2Client } from '@tesote.com/sdk';

const client = new V2Client({ apiKey: process.env.TESOTE_API_KEY! });

const accounts = await client.accounts.list();
const acct = await client.accounts.get('acct_123');

// Rate-limit headers from the most recent response.
console.log(client.lastRateLimit);
```

V1 and V2 clients ship side-by-side and are picked explicitly:

```ts
import { V1Client, V2Client } from '@tesote.com/sdk';
```

## Errors

Every error the SDK throws is a typed subclass of `TesoteError` and carries
`errorCode`, `httpStatus`, `requestId`, `errorId`, `retryAfter`,
`responseBody`, `requestSummary`, and `attempts`.

```ts
import { RateLimitExceededError } from '@tesote.com/sdk';

try {
  await client.accounts.list();
} catch (err) {
  if (err instanceof RateLimitExceededError) {
    console.log('retry after', err.retryAfter, 'seconds');
  }
  throw err;
}
```

## Docs

Full reference: <https://www.tesote.com/docs/ts>

## License

MIT
