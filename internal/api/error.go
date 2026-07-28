package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ccrsxx/api/internal/utils"
	"github.com/google/uuid"
)

var isDevelopmentMode bool

func Load(isDevelopment bool) {
	isDevelopmentMode = isDevelopment
}

type PanicError struct {
	Value   any
	Stack   string
	Message string
}

func (e *PanicError) Error() string {
	return e.Message
}

type HTTPError struct {
	Message    string
	Details    []string
	StatusCode int
	// Code is an optional stable, machine-readable identifier for this error
	// (e.g. "turnstile_invalid"). It is used as the error_kind fingerprint for
	// log-based error grouping. When empty, error_kind falls back to the
	// status code. Adding this field is backwards compatible because every
	// construction site uses named struct literals.
	Code string
}

func (e *HTTPError) Error() string {
	return e.Message
}

func HandleHTTPError(w http.ResponseWriter, r *http.Request, err error) {
	errorID := uuid.New().String()

	ctx := r.Context()

	ipAddress := utils.GetIPAddressFromRequest(r)

	// route is low cardinality (matched mux pattern), safe for grouping.
	route := r.Pattern

	if route == "" {
		route = "unmatched"
	}

	if panicErr, ok := errors.AsType[*PanicError](err); ok {
		parsedStack := panicErr.Stack

		if isDevelopmentMode {
			fmt.Printf("panic stack trace:\n%s\n", panicErr.Stack)

			parsedStack = "printed to stdout in development mode"
		}

		// Panics are grouped by crash location (first application stack frame),
		// mirroring how hosted error reporting fingerprints crashes.
		errorKind := extractPanicOrigin(panicErr.Stack)

		slog.ErrorContext(ctx, "http panic error",
			"message", panicErr.Message,
			"value", panicErr.Value,
			"stack", parsedStack,
			"error", err,
			"error_id", errorID,
			"error_kind", errorKind,
			"route", route,
			"path", r.URL.Path,
			"method", r.Method,
			"ip_address", ipAddress,
		)

		if err := NewErrorResponse(w, http.StatusInternalServerError, "An internal server error occurred", nil, errorID); err != nil {
			logErrorResponse(err, errorID)
		}

		return
	}

	if httpErr, ok := errors.AsType[*HTTPError](err); ok {
		errorKind := httpErr.Code

		if errorKind == "" {
			errorKind = fmt.Sprintf("http_%d", httpErr.StatusCode)
		}

		// Severity by status class so level="error" means "my fault":
		//   5xx -> error, 429 -> info (expected, high volume), other 4xx -> warn.
		logErrorByStatus(ctx, httpErr.StatusCode, "http handled error",
			"message", httpErr.Message,
			"status_code", httpErr.StatusCode,
			"details", httpErr.Details,
			"error", err,
			"error_id", errorID,
			"error_kind", errorKind,
			"route", route,
			"path", r.URL.Path,
			"method", r.Method,
			"ip_address", ipAddress,
		)

		if err := NewErrorResponse(w, httpErr.StatusCode, httpErr.Message, httpErr.Details, errorID); err != nil {
			logErrorResponse(err, errorID)
		}

		return
	}

	// Any unhandled errors. Fingerprint by Go type name, which is naturally
	// stable and bounded.
	slog.ErrorContext(ctx, "http unhandled error",
		"error", err,
		"error_id", errorID,
		"error_kind", fmt.Sprintf("%T", err),
		"route", route,
		"path", r.URL.Path,
		"method", r.Method,
		"ip_address", ipAddress,
	)

	if err := NewErrorResponse(w, http.StatusInternalServerError, "An internal server error occurred", nil, errorID); err != nil {
		logErrorResponse(err, errorID)
	}
}

// logErrorByStatus emits at a severity appropriate for the HTTP status class.
func logErrorByStatus(ctx context.Context, statusCode int, msg string, args ...any) {
	switch {
	case statusCode >= 500:
		slog.ErrorContext(ctx, msg, args...)
	case statusCode == http.StatusTooManyRequests:
		// Rate limiting is expected and high volume; do not pollute errors.
		slog.InfoContext(ctx, msg, args...)
	default:
		// Other 4xx are client mistakes, not server faults.
		slog.WarnContext(ctx, msg, args...)
	}
}

// extractPanicOrigin returns the first application stack frame from a panic
// stack, used as a stable fingerprint for grouping panics. Runtime, debug and
// recovery frames are skipped. Returns "unknown" when nothing suitable is
// found.
func extractPanicOrigin(stack string) string {
	lines := strings.Split(stack, "\n")

	for _, line := range lines {
		// Function frames are not indented; file frames start with a tab.
		if line == "" || strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ") {
			continue
		}

		// Skip the goroutine header line.
		if strings.HasPrefix(line, "goroutine ") {
			continue
		}

		// Skip runtime/debug/recovery internals.
		if strings.HasPrefix(line, "runtime.") ||
			strings.HasPrefix(line, "runtime/debug.") ||
			strings.Contains(line, "internal/api.Recovery") ||
			strings.Contains(line, "internal/api.HandleHTTPError") ||
			strings.Contains(line, "internal/middleware.Recovery") {
			continue
		}

		// Trim the argument list, e.g. "pkg.Func(0x1, 0x2)" -> "pkg.Func".
		if idx := strings.Index(line, "("); idx != -1 {
			line = line[:idx]
		}

		return strings.TrimSpace(line)
	}

	return "unknown"
}

func logErrorResponse(err error, errorID string) {
	slog.Error("send error response failed", "error", err, "error_id", errorID)
}
