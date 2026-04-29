<?php

declare(strict_types=1);

namespace Tesote\Sdk\V2;

use Tesote\Sdk\Models\StatusInfo;
use Tesote\Sdk\Models\WhoAmI;
use Tesote\Sdk\Transport;

/** GET /v2/status (no auth) and GET /v2/whoami. */
final class Status
{
    public function __construct(private readonly Transport $transport)
    {
    }

    public function check(): StatusInfo
    {
        $body = $this->transport->request('GET', '/v2/status') ?? [];
        return StatusInfo::fromArray($body);
    }

    public function whoami(): WhoAmI
    {
        $body = $this->transport->request('GET', '/v2/whoami') ?? [];
        return WhoAmI::fromArray($body);
    }
}
