<?php

declare(strict_types=1);

namespace Tesote\Sdk\V1;

use Tesote\Sdk\Models\StatusInfo;
use Tesote\Sdk\Models\WhoAmI;
use Tesote\Sdk\Transport;

/** GET /status (no auth) and GET /whoami. */
final class Status
{
    public function __construct(private readonly Transport $transport)
    {
    }

    public function check(): StatusInfo
    {
        $body = $this->transport->request('GET', '/status') ?? [];
        return StatusInfo::fromArray($body);
    }

    public function whoami(): WhoAmI
    {
        $body = $this->transport->request('GET', '/whoami') ?? [];
        return WhoAmI::fromArray($body);
    }
}
