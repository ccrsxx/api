package observability

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"
	"go.opentelemetry.io/otel/trace"
)

// TraceHandler wraps an HTTP handler with OpenTelemetry instrumentation.
//
// The span is initially named after the HTTP method only. The default
// formatter would use the raw URL path, producing one distinct span name per
// URL (e.g. every blog slug), which explodes Tempo's span name index.
// TraceRoute renames the span to the matched route once the mux has resolved
// it.
//
// /sse is excluded: long-lived streaming connections would produce multi-hour
// spans that distort trace views and waste storage.
//
// When tracing is disabled the global tracer provider is a no-op, so this adds
// negligible overhead and needs no conditional wiring.
func TraceHandler(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "",
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return r.Method
		}),
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/sse"
		}),
	)
}

// TraceRoute renames the active span to the matched mux pattern and records it
// as the http.route attribute.
//
// It must be registered as the INNERMOST middleware (directly wrapping the
// router), because http.ServeMux only populates r.Pattern while dispatching.
// Renaming after next.ServeHTTP returns is safe: otelhttp ends the span after
// the whole chain unwinds, and SetName before End is valid.
//
// Unmatched requests keep the bare method as their span name, which stays low
// cardinality.
func TraceRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if r.Pattern == "" {
			return
		}

		span := trace.SpanFromContext(r.Context())

		if !span.IsRecording() {
			return
		}

		span.SetName(r.Method + " " + r.Pattern)
		span.SetAttributes(semconv.HTTPRoute(r.Pattern))
	})
}
