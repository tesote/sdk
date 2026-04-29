<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 422 BANK_SUBMISSION_ERROR — upstream bank rejected the order at submission time. */
final class BankSubmissionException extends ApiException
{
}
