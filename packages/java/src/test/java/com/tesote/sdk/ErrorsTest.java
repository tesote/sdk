package com.tesote.sdk;

import com.tesote.sdk.errors.AccountDisabledException;
import com.tesote.sdk.errors.AccountNotFoundException;
import com.tesote.sdk.errors.ApiException;
import com.tesote.sdk.errors.ApiKeyRevokedException;
import com.tesote.sdk.errors.BankConnectionNotFoundException;
import com.tesote.sdk.errors.BankSubmissionException;
import com.tesote.sdk.errors.BankUnderMaintenanceException;
import com.tesote.sdk.errors.BatchNotFoundException;
import com.tesote.sdk.errors.BatchValidationException;
import com.tesote.sdk.errors.ErrorDispatcher;
import com.tesote.sdk.errors.HistorySyncForbiddenException;
import com.tesote.sdk.errors.InternalErrorException;
import com.tesote.sdk.errors.InvalidCountException;
import com.tesote.sdk.errors.InvalidCursorException;
import com.tesote.sdk.errors.InvalidDateRangeException;
import com.tesote.sdk.errors.InvalidLimitException;
import com.tesote.sdk.errors.InvalidOrderStateException;
import com.tesote.sdk.errors.InvalidQueryException;
import com.tesote.sdk.errors.MissingDateRangeException;
import com.tesote.sdk.errors.MutationDuringPaginationException;
import com.tesote.sdk.errors.NotFoundException;
import com.tesote.sdk.errors.PaymentMethodNotFoundException;
import com.tesote.sdk.errors.RateLimitExceededException;
import com.tesote.sdk.errors.RequestSummary;
import com.tesote.sdk.errors.ServiceUnavailableException;
import com.tesote.sdk.errors.SyncInProgressException;
import com.tesote.sdk.errors.SyncRateLimitExceededException;
import com.tesote.sdk.errors.SyncSessionNotFoundException;
import com.tesote.sdk.errors.TransactionNotFoundException;
import com.tesote.sdk.errors.TransactionOrderNotFoundException;
import com.tesote.sdk.errors.UnauthorizedException;
import com.tesote.sdk.errors.UnprocessableContentException;
import com.tesote.sdk.errors.ValidationException;
import com.tesote.sdk.errors.WorkspaceSuspendedException;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ErrorsTest {
    private static RequestSummary summary() {
        return new RequestSummary("GET", "/v3/accounts",
                Map.of(), null, "Bearer ****1234");
    }

    private static ApiException dispatch(String code, int status) {
        return ErrorDispatcher.dispatch("msg", code, status, "req_1", "err_1",
                null, "{}", summary(), 1, null);
    }

    @Test
    void unauthorizedMaps() {
        assertInstanceOf(UnauthorizedException.class, dispatch("UNAUTHORIZED", 401));
    }

    @Test
    void apiKeyRevokedMaps() {
        assertInstanceOf(ApiKeyRevokedException.class, dispatch("API_KEY_REVOKED", 401));
    }

    @Test
    void workspaceSuspendedMaps() {
        assertInstanceOf(WorkspaceSuspendedException.class, dispatch("WORKSPACE_SUSPENDED", 403));
    }

    @Test
    void accountDisabledMaps() {
        assertInstanceOf(AccountDisabledException.class, dispatch("ACCOUNT_DISABLED", 403));
    }

    @Test
    void historySyncForbiddenMaps() {
        assertInstanceOf(HistorySyncForbiddenException.class, dispatch("HISTORY_SYNC_FORBIDDEN", 403));
    }

    @Test
    void mutationConflictMaps() {
        assertInstanceOf(MutationDuringPaginationException.class, dispatch("MUTATION_CONFLICT", 409));
    }

    @Test
    void unprocessableContentMaps() {
        assertInstanceOf(UnprocessableContentException.class, dispatch("UNPROCESSABLE_CONTENT", 422));
    }

    @Test
    void invalidDateRangeMaps() {
        ApiException ex = dispatch("INVALID_DATE_RANGE", 422);
        assertInstanceOf(InvalidDateRangeException.class, ex);
        // why: subclass of UnprocessableContentException so callers can catch the parent.
        assertInstanceOf(UnprocessableContentException.class, ex);
    }

    @Test
    void rateLimitExceededMaps() {
        assertInstanceOf(RateLimitExceededException.class, dispatch("RATE_LIMIT_EXCEEDED", 429));
    }

    @Test
    void serviceUnavailableMapsByStatusWithEmptyCode() {
        ApiException ex = dispatch("", 503);
        assertInstanceOf(ServiceUnavailableException.class, ex);
    }

    @Test
    void unknownCodeFallsBackToApiException() {
        ApiException ex = dispatch("MYSTERY_CODE", 418);
        assertEquals(ApiException.class, ex.getClass());
        assertEquals("MYSTERY_CODE", ex.errorCode());
    }

    @Test
    void requiredFieldsPopulated() {
        ApiException ex = dispatch("UNAUTHORIZED", 401);
        assertEquals(401, ex.httpStatus());
        assertEquals("UNAUTHORIZED", ex.errorCode());
        assertEquals("req_1", ex.requestId());
        assertEquals("err_1", ex.errorId());
        assertEquals(1, ex.attempts());
        assertNotNull(ex.requestSummary());
    }

    @Test
    void bearerRedactionUtility() {
        String redacted = Transport.redactBearer("sk_test_abcd1234");
        assertTrue(redacted.startsWith("Bearer ****"));
        assertTrue(redacted.endsWith("1234"));
    }

    @Test
    void shortKeyStillRedacted() {
        String redacted = Transport.redactBearer("ab");
        assertEquals("Bearer ****", redacted);
    }

    @Test
    void notFoundFamilyAllSubclassNotFound() {
        assertInstanceOf(AccountNotFoundException.class, dispatch("ACCOUNT_NOT_FOUND", 404));
        assertInstanceOf(NotFoundException.class, dispatch("ACCOUNT_NOT_FOUND", 404));
        assertInstanceOf(TransactionNotFoundException.class, dispatch("TRANSACTION_NOT_FOUND", 404));
        assertInstanceOf(NotFoundException.class, dispatch("TRANSACTION_NOT_FOUND", 404));
        assertInstanceOf(SyncSessionNotFoundException.class, dispatch("SYNC_SESSION_NOT_FOUND", 404));
        assertInstanceOf(PaymentMethodNotFoundException.class, dispatch("PAYMENT_METHOD_NOT_FOUND", 404));
        assertInstanceOf(TransactionOrderNotFoundException.class, dispatch("TRANSACTION_ORDER_NOT_FOUND", 404));
        assertInstanceOf(BatchNotFoundException.class, dispatch("BATCH_NOT_FOUND", 404));
        assertInstanceOf(BankConnectionNotFoundException.class, dispatch("BANK_CONNECTION_NOT_FOUND", 404));
    }

    @Test
    void unprocessableFamilyAllSubclassUnprocessable() {
        assertInstanceOf(InvalidCursorException.class, dispatch("INVALID_CURSOR", 422));
        assertInstanceOf(UnprocessableContentException.class, dispatch("INVALID_CURSOR", 422));
        assertInstanceOf(InvalidCountException.class, dispatch("INVALID_COUNT", 422));
        assertInstanceOf(InvalidLimitException.class, dispatch("INVALID_LIMIT", 422));
        assertInstanceOf(InvalidQueryException.class, dispatch("INVALID_QUERY", 422));
        assertInstanceOf(MissingDateRangeException.class, dispatch("MISSING_DATE_RANGE", 422));
        assertInstanceOf(BankSubmissionException.class, dispatch("BANK_SUBMISSION_ERROR", 422));
    }

    @Test
    void validationFamilyAllSubclassValidation() {
        assertInstanceOf(ValidationException.class, dispatch("VALIDATION_ERROR", 400));
        assertInstanceOf(BatchValidationException.class, dispatch("BATCH_VALIDATION_ERROR", 400));
        assertInstanceOf(ValidationException.class, dispatch("BATCH_VALIDATION_ERROR", 400));
    }

    @Test
    void conflictAndSyncErrorsMap() {
        assertInstanceOf(InvalidOrderStateException.class, dispatch("INVALID_ORDER_STATE", 409));
        assertInstanceOf(SyncInProgressException.class, dispatch("SYNC_IN_PROGRESS", 409));
        assertInstanceOf(SyncRateLimitExceededException.class, dispatch("SYNC_RATE_LIMIT_EXCEEDED", 429));
        assertInstanceOf(BankUnderMaintenanceException.class, dispatch("BANK_UNDER_MAINTENANCE", 503));
        assertInstanceOf(InternalErrorException.class, dispatch("INTERNAL_ERROR", 500));
    }
}
