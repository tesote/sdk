# tesote/sdk — packages/php

> Workspace brain — established by the developerz.ai maintainer agent for the maintainers of this repository from this workspace, then owned by
> its maintainers. Correct it in place; setup never overwrites a brain that exists.

PHP SDK library, independently tested and published through Packagist.

The gate is the repo root’s: `bin/check` (see the root `CLAUDE.md`).

## Structure

- `src` is autoloaded as `Tesote\Sdk\`; `tests` is autoloaded as `Tesote\Sdk\Tests\`. `composer.json:18-26`
- PHPUnit discovers `tests` and covers `src`; PHPStan checks both. `packages/php/phpunit.xml.dist:10-19`, `packages/php/phpstan.neon:1-6`

## Conventions

- Require PHP 8.1 with curl and JSON extensions. `composer.json:8-16`
- Use strict types, PSR-12, sorted imports, single quotes, and multiline trailing commas. `packages/php/.php-cs-fixer.dist.php:1-19`
- Run `composer install --no-interaction --no-progress --prefer-dist`, `composer cs-check`, `composer phpstan -- --no-progress`, and `composer test -- --colors=never`. `.github/workflows/php.yml:46-53`
- Use fake curl transport tests, not live HTTP. `packages/php/tests/TransportTest.php:17-45`
