<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /v1/accounts and GET /v2/accounts. */
final readonly class AccountList
{
    /**
     * @param list<Account> $accounts
     */
    public function __construct(
        public int $total,
        public array $accounts,
        public PagePagination $pagination,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $accounts = [];
        foreach ((is_array($data['accounts'] ?? null) ? $data['accounts'] : []) as $entry) {
            if (is_array($entry)) {
                $accounts[] = Account::fromArray($entry);
            }
        }
        $pagination = is_array($data['pagination'] ?? null) ? $data['pagination'] : [];

        return new self(
            total: (int) ($data['total'] ?? 0),
            accounts: $accounts,
            pagination: PagePagination::fromArray($pagination),
        );
    }
}
