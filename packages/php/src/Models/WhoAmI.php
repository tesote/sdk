<?php

declare(strict_types=1);

namespace Tesote\Sdk\Models;

/** Response from GET /whoami and GET /v2/whoami. */
final readonly class WhoAmI
{
    public function __construct(
        public string $id,
        public string $name,
        public string $type,
    ) {
    }

    /**
     * @param array<string, mixed> $data
     */
    public static function fromArray(array $data): self
    {
        $client = is_array($data['client'] ?? null) ? $data['client'] : [];
        return new self(
            id: (string) ($client['id'] ?? ''),
            name: (string) ($client['name'] ?? ''),
            type: (string) ($client['type'] ?? ''),
        );
    }
}
