package tesote

import "fmt"

// TransportError is the base for network-level failures (no usable HTTP
// response). Subtypes embed it and override Unwrap.
type TransportError struct {
	Op             string
	Message        string
	RequestSummary RequestSummary
	Attempts       int
	Cause          error
}

// Error implements error.
func (e *TransportError) Error() string {
	return fmt.Sprintf("tesote: %s: %s (attempts=%d)", e.Op, e.Message, e.Attempts)
}

// NetworkError is raised on DNS, connection refused, reset.
type NetworkError struct{ *TransportError }

// Unwrap returns the wrapped cause and matches ErrNetwork via errors.Is.
func (e *NetworkError) Unwrap() error { return e.Cause }

// Is satisfies errors.Is for the ErrNetwork sentinel.
func (e *NetworkError) Is(target error) bool { return target == ErrNetwork }

// TimeoutError is raised on connect/read timeout.
type TimeoutError struct{ *TransportError }

// Unwrap returns the wrapped cause.
func (e *TimeoutError) Unwrap() error { return e.Cause }

// Is satisfies errors.Is for the ErrTimeout sentinel.
func (e *TimeoutError) Is(target error) bool { return target == ErrTimeout }

// TLSError is raised on certificate / handshake failures.
type TLSError struct{ *TransportError }

// Unwrap returns the wrapped cause.
func (e *TLSError) Unwrap() error { return e.Cause }

// Is satisfies errors.Is for the ErrTLS sentinel.
func (e *TLSError) Is(target error) bool { return target == ErrTLS }

// ConfigError is raised at client construction for bad SDK config.
type ConfigError struct {
	Field   string
	Message string
}

// Error implements error.
func (e *ConfigError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("tesote: config error: %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("tesote: config error: %s", e.Message)
}

// Is satisfies errors.Is for the ErrConfig sentinel.
func (e *ConfigError) Is(target error) bool { return target == ErrConfig }

// EndpointRemovedError is returned when the SDK exposes a method whose upstream
// endpoint has been removed from this API version.
type EndpointRemovedError struct {
	Method      string
	Replacement string
}

// Error implements error.
func (e *EndpointRemovedError) Error() string {
	if e.Replacement != "" {
		return fmt.Sprintf("tesote: %s removed; use %s", e.Method, e.Replacement)
	}
	return fmt.Sprintf("tesote: %s removed", e.Method)
}

// Is satisfies errors.Is for the ErrEndpointRemoved sentinel.
func (e *EndpointRemovedError) Is(target error) bool { return target == ErrEndpointRemoved }
