<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 503 BANK_UNDER_MAINTENANCE — upstream bank reporting maintenance window; retry later. */
final class BankUnderMaintenanceException extends ApiException
{
}
