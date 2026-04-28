package com.tesote.sdk;

import com.tesote.sdk.errors.AccountDisabledException;
import com.tesote.sdk.errors.ApiException;
import com.tesote.sdk.errors.ApiKeyRevokedException;
import com.tesote.sdk.errors.ErrorDispatcher;
import com.tesote.sdk.errors.HistorySyncForbiddenException;
import com.tesote.sdk.errors.InvalidDateRangeException;
import com.tesote.sdk.errors.MutationDuringPaginationException;
import com.tesote.sdk.errors.RateLimitExceededException;
import com.tesote.sdk.errors.RequestSummary;
import com.tesote.sdk.errors.ServiceUnavailableException;
import com.tesote.sdk.errors.UnauthorizedException;
import com.tesote.sdk.errors.UnprocessableContentException;
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
}
