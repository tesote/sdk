<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 409 SYNC_IN_PROGRESS — a previous sync for this account is still running. */
final class SyncInProgressException extends ApiException
{
}
