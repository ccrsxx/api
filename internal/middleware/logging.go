package middleware

import (
	"cmp"
	"log/slog"
	"net/http"
	"time"

	"github.com/ccrsxx/api/internal/utils"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += n
	return n, err
}

// Unwrap provides access to the underlying ResponseWriter
// To reserve compatibility with http.ResponseWriter wrappers
// Like http.CloseNotifier, http.Flusher, etc.
func (w *wrappedWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip SSE: long-lived connections would be logged as multi-hour
		// "requests" and wreck any latency aggregation built on duration_ms.
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

		duration := time.Since(start)

		ipAddress := utils.GetIPAddressFromRequest(r)

		// route is the matched mux pattern (low cardinality, safe for metrics).
		// path is the raw URL (high cardinality, debugging only). r.Pattern is
		// populated by the mux during ServeHTTP, so it is available here.
		route := cmp.Or(r.Pattern, "unmatched")

		slog.InfoContext(r.Context(), "http request",
			"route", route,
			"path", r.URL.Path,
			"method", r.Method,
			"status_code", wrapped.statusCode,
			"duration_ms", float64(duration.Microseconds())/1000,
			"bytes", wrapped.bytesWritten,
			"ip_address", ipAddress,
		)
	})
}
