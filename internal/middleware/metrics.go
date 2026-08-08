package middleware

import (
	"cmp"
	"net/http"
	"strconv"
	"time"

	"github.com/ccrsxx/api/internal/observability"
)

// Metrics records request counts and latency for Prometheus.
//
// Only the matched mux pattern (r.Pattern) is used as the route label, never
// the raw path, to keep cardinality bounded.
//
// SSE is skipped for the same reason it is skipped in Logging: a long-lived
// stream would be recorded as one multi-hour request and ruin the latency
// histogram.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/sse" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		route := cmp.Or(r.Pattern, "unmatched")

		observability.HTTPRequestsTotal.WithLabelValues(
			route,
			r.Method,
			strconv.Itoa(wrapped.statusCode),
		).Inc()

		observability.HTTPRequestDuration.WithLabelValues(
			route,
			r.Method,
		).Observe(time.Since(start).Seconds())
	})
}
