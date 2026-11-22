# API

> established by the developerz.ai maintainer agent for the maintainers of this repository from what onboarding read and ran in this repo, then owned by
> its maintainers. Correct it in place; setup never overwrites a page that exists.

The SDK sends bearer API-key authentication, with explicit client configuration taking precedence over `TESOTE_SDK_API_KEY`; see `docs/architecture/auth.md:3-24`. v1 exposes read-only accounts and transactions, while v2 adds sync sessions, transaction orders, batches, payment methods, bulk, and search. `docs/architecture/resources.md:5-27`.
