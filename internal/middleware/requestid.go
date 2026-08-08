package middleware

import (
	"net/http"
	"strings"

	"github.com/ccrsxx/api/internal/logger"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-Id"

// RequestID establishes per-request correlation identifiers and stores them in
// the request context so the logger can attach them to every log line.
//
// It must be registered as the OUTERMOST middleware so the identifiers exist
// before any other middleware (recovery, logging, rate limiting) runs.
//
// Sources, in order of preference:
//  1. W3C "traceparent" header (trace_id/span_id) if a caller propagated one.
//  2. An inbound X-Request-Id header (lets a frontend/proxy set its own id).
//  3. A freshly generated UUID.
//
// When OpenTelemetry is added, the active span becomes the source of truth for
// trace_id/span_id; this middleware then moves inside the otel handler.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID, spanID := parseTraceparent(r.Header.Get("traceparent"))

		requestID := r.Header.Get(requestIDHeader)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		info := logger.RequestInfo{
			RequestID: requestID,
			TraceID:   traceID,
			SpanID:    spanID,
		}

		// Echo the request id back so clients can quote it when reporting issues.
		w.Header().Set(requestIDHeader, requestID)

		ctx := logger.WithRequestInfo(r.Context(), info)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// parseTraceparent extracts the trace-id and span-id from a W3C traceparent
// header. Returns empty strings when the header is absent or malformed.
//
// Format: "<version>-<trace-id>-<parent-id>-<flags>"
// Example: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
func parseTraceparent(header string) (traceID, spanID string) {
	if header == "" {
		return "", ""
	}

	parts := strings.Split(header, "-")

	if len(parts) != 4 {
		return "", ""
	}

	// trace-id is 32 hex chars, parent/span-id is 16 hex chars.
	if len(parts[1]) != 32 || len(parts[2]) != 16 {
		return "", ""
	}

	// An all-zero trace-id is invalid per spec.
	if strings.Trim(parts[1], "0") == "" {
		return "", ""
	}

	return parts[1], parts[2]
}
