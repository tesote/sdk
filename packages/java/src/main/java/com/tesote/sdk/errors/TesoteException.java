package com.tesote.sdk.errors;

/**
 * Root of the SDK exception hierarchy.
 *
 * <p>Unchecked because checked exceptions pollute the call sites of every resource
 * method without giving callers any information they couldn't get from typed
 * subclasses dispatched in a single {@code catch}.
 *
 * <p>Underlying causes are preserved through {@link #getCause()} — never lose the
 * chain.
 */
public class TesoteException extends RuntimeException {
    private final String errorCode;
    private final int httpStatus;
    private final String requestId;
    private final String errorId;
    private final Integer retryAfter;
    private final String responseBody;
    private final RequestSummary requestSummary;
    private final int attempts;

    public TesoteException(
            String message,
            String errorCode,
            int httpStatus,
            String requestId,
            String errorId,
            Integer retryAfter,
            String responseBody,
            RequestSummary requestSummary,
            int attempts,
            Throwable cause
    ) {
        super(message, cause);
        this.errorCode = errorCode;
        this.httpStatus = httpStatus;
        this.requestId = requestId;
        this.errorId = errorId;
        this.retryAfter = retryAfter;
        this.responseBody = responseBody;
        this.requestSummary = requestSummary;
        this.attempts = attempts;
    }

    public String errorCode() { return errorCode; }
    public int httpStatus() { return httpStatus; }
    public String requestId() { return requestId; }
    public String errorId() { return errorId; }
    public Integer retryAfter() { return retryAfter; }
    public String responseBody() { return responseBody; }
    public RequestSummary requestSummary() { return requestSummary; }
    public int attempts() { return attempts; }

    @Override
    public String getMessage() {
        StringBuilder sb = new StringBuilder();
        sb.append(getClass().getSimpleName()).append(": ");
        if (httpStatus > 0) sb.append(httpStatus).append(' ');
        sb.append(super.getMessage() == null ? "" : super.getMessage());
        if (errorCode != null) sb.append("\n  error_code: ").append(errorCode);
        if (requestId != null) sb.append("\n  request_id: ").append(requestId);
        if (errorId != null) sb.append("\n  error_id: ").append(errorId);
        if (retryAfter != null) sb.append("\n  retry_after: ").append(retryAfter).append('s');
        if (attempts > 0) sb.append("\n  attempts: ").append(attempts);
        if (requestSummary != null) {
            sb.append("\n  request: ").append(requestSummary.method())
                    .append(' ').append(requestSummary.path());
            if (requestSummary.bodyShape() != null) {
                sb.append(" (body: ").append(requestSummary.bodyShape()).append(')');
            }
        }
        if (responseBody != null && !responseBody.isEmpty()) {
            sb.append("\n  response: ").append(responseBody);
        }
        return sb.toString();
    }
}
