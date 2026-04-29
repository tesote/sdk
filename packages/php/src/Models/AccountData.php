<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/**
 * Account.data subobject — wire-side snake_case preserved on raw map; idiomatic
 * camelCase exposed on properties.
 *
 * balanceCents / availableBalanceCents only present when the workspace allows
 * balance display. Returned as strings on the wire (decimal-safe).
 */
final class AccountData
{
    public function __construct(
        public readonly ?string $maskedAccountNumber,
        public readonly ?string $currency,
        public readonly ?string $transactionsDataCurrentAsOf,
        public readonly ?string $balanceDataCurrentAsOf,
        public readonly ?string $customUserProvidedIdentifier,
        public readonly ?string $balanceCents,
        public readonly ?string $availableBalanceCents,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        return new self(
            maskedAccountNumber: isset($data['masked_account_number']) ? (string) $data['masked_account_number'] : null,
            currency: isset($data['currency']) ? (string) $data['currency'] : null,
            transactionsDataCurrentAsOf: isset($data['transactions_data_current_as_of']) ? (string) $data['transactions_data_current_as_of'] : null,
            balanceDataCurrentAsOf: isset($data['balance_data_current_as_of']) ? (string) $data['balance_data_current_as_of'] : null,
            customUserProvidedIdentifier: isset($data['custom_user_provided_identifier']) ? (string) $data['custom_user_provided_identifier'] : null,
            balanceCents: isset($data['balance_cents']) ? (string) $data['balance_cents'] : null,
            availableBalanceCents: isset($data['available_balance_cents']) ? (string) $data['available_balance_cents'] : null,
        );
    }
}
