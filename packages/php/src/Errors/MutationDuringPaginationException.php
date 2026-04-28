<?php

declare(strict_types=1);

namespace Tesote\Sdk\Errors;

/** 409 MUTATION_CONFLICT — collection changed mid-pagination; restart with new cursor. */
final class MutationDuringPaginationException extends ApiException
{
}
