package cmdutil

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/algolia/cli/api/dashboard"
)

// error_source buckets for failure telemetry: where an error came from, so
// failures can be triaged without parsing the (high-cardinality) message.
const (
	ErrorSourceNetwork    = "network"
	ErrorSourceAPI        = "api"
	ErrorSourceValidation = "validation"
	ErrorSourceLocal      = "local"
)

// httpStatusError is satisfied by errors carrying an HTTP status (e.g.
// *dashboard.APIError), declared here to keep the classifier decoupled from it.
type httpStatusError interface{ HTTPStatusCode() int }

// ClassifyError returns a stable, low-cardinality error class, a source bucket
// (ErrorSource*), and the HTTP status when present (0 otherwise), so every
// "*Failed" event classifies errors the same way. Most specific checks first.
func ClassifyError(err error) (class, source string, httpStatus int) {
	if err == nil {
		return "", "", 0
	}

	var flagErr *FlagError
	if errors.As(err, &flagErr) {
		return "validation_error", ErrorSourceValidation, 0
	}

	if errors.Is(err, dashboard.ErrSessionExpired) {
		return "session_expired", ErrorSourceAPI, http.StatusUnauthorized
	}
	var clusterErr *dashboard.ErrClusterUnavailable
	if errors.As(err, &clusterErr) {
		return "cluster_unavailable", ErrorSourceAPI, 0
	}
	var statusErr httpStatusError
	if errors.As(err, &statusErr) {
		status := statusErr.HTTPStatusCode()
		return fmt.Sprintf("http_%d", status), ErrorSourceAPI, status
	}

	if errors.Is(err, context.Canceled) {
		return "canceled", ErrorSourceLocal, 0
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "timeout", ErrorSourceNetwork, 0
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return "dns_error", ErrorSourceNetwork, 0
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return "timeout", ErrorSourceNetwork, 0
		}
		return "network_error", ErrorSourceNetwork, 0
	}

	// Fall back to the root cause's type name: stable and low-cardinality for
	// typed errors, "unknown" for generic ones.
	return rootCauseType(err), ErrorSourceLocal, 0
}

// ErrorClass returns only the class component of ClassifyError.
func ErrorClass(err error) string {
	class, _, _ := ClassifyError(err)
	return class
}

// rootCauseType returns the type name of the deepest wrapped error, or
// "unknown" for the generic types every errors.New/fmt.Errorf produces —
// otherwise that one type name would dominate the error_class dimension
// while carrying no signal.
func rootCauseType(err error) string {
	for {
		next := errors.Unwrap(err)
		if next == nil {
			break
		}
		err = next
	}
	switch name := fmt.Sprintf("%T", err); name {
	case "*errors.errorString", "*errors.joinError", "*fmt.wrapError", "*fmt.wrapErrors":
		return "unknown"
	default:
		return name
	}
}
