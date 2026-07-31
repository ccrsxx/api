package middleware

import (
	"net/http"
	"strings"

	"github.com/ccrsxx/api/internal/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const requestIDHeader = "X-Request-Id"

// RequestID establishes per-request correlation identifiers and stores them in
// the request context so the logger can attach them to every log line.
//
// It runs inside the OpenTelemetry HTTP handler (so an active span exists) but
// outside recovery/logging/rate limiting, so the identifiers are available to
// every log line those middlewares produce, including panics.
//
// Sources, in order of preference:
//  1. The active OpenTelemetry span (authoritative when tracing is enabled).
//  2. W3C "traceparent" header, if a caller propagated one.
//  3. An inbound X-Request-Id header (lets a frontend/proxy set its own id).
//  4. A freshly generated UUID.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID, spanID := traceIDsFromRequest(r)

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

// traceIDsFromRequest resolves the trace and span id for this request.
//
// When tracing is enabled, an OpenTelemetry span is already active (the otel
// HTTP handler wraps this middleware), and that span is the source of truth:
// using it guarantees the trace_id in the logs is the same one Tempo stores,
// which is what makes the log -> trace jump in Grafana work.
//
// When tracing is disabled there is no recording span, so we fall back to
// parsing the inbound W3C traceparent header. This keeps correlation working
// for callers that propagate trace context even while our own tracing is off.
func traceIDsFromRequest(r *http.Request) (traceID, spanID string) {
	spanCtx := trace.SpanContextFromContext(r.Context())

	if spanCtx.IsValid() {
		return spanCtx.TraceID().String(), spanCtx.SpanID().String()
	}

	return parseTraceparent(r.Header.Get("traceparent"))
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
