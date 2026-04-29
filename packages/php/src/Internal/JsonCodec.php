<?php

declare(strict_types=1);

namespace Tesote\Sdk\Internal;

use Tesote\Sdk\Errors\ConfigException;

/**
 * JSON encode/decode with the SDK's chosen flag set and error contract.
 */
final class JsonCodec
{
    private const ENCODE_FLAGS = JSON_THROW_ON_ERROR | JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE;

    /**
     * @param array<mixed> $body
     */
    public static function encode(array $body): string
    {
        try {
            return json_encode($body, self::ENCODE_FLAGS);
        } catch (\JsonException $e) {
            throw new ConfigException('Failed to JSON-encode request body: ' . $e->getMessage());
        }
    }

    /**
     * @return array<mixed>|null
     */
    public static function decode(string $body): ?array
    {
        if ($body === '') {
            return null;
        }
        try {
            $decoded = json_decode($body, true, 512, JSON_THROW_ON_ERROR);
        } catch (\JsonException) {
            return null;
        }
        return is_array($decoded) ? $decoded : null;
    }
}
