<?php

declare(strict_types=1);

namespace Tesote\Sdk;

use LogicException;

/**
 * Stub used by version Clients to expose a resource property whose methods
 * are not yet wired. Calling any method throws so missing wiring fails
 * loudly at the call site rather than silently no-oping.
 *
 * Concrete resource classes (Accounts, Transactions, ...) replace these
 * stubs as wiring lands.
 */
final class NotImplemented
{
    public function __construct(private readonly string $resource)
    {
    }

    /**
     * @param  array<int, mixed> $arguments
     * @throws LogicException
     */
    public function __call(string $name, array $arguments): never
    {
        throw new LogicException(sprintf(
            'not implemented: %s.%s() — wiring pending in this SDK version',
            $this->resource,
            $name,
        ));
    }
}
