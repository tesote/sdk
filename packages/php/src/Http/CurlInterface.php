<?php

declare(strict_types=1);

namespace Tesote\Sdk\Http;

/**
 * Minimal cURL surface, behind an interface so tests can swap a fake.
 *
 * The shape mirrors the ext-curl procedural API one-for-one rather than
 * attempting to hide it; we want the production path to stay obviously
 * a thin wrapper, not a leaky abstraction.
 */
interface CurlInterface
{
    /**
     * @param array<int, mixed> $options curl_setopt_array-compatible options.
     *
     * @return CurlResult
     */
    public function execute(array $options): CurlResult;
}
