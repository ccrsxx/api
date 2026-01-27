package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ccrsxx/api/src/internal/utils"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// Unwrap provides access to the underlying ResponseWriter
// To reserve compatibility with http.ResponseWriter wrappers
// Like http.CloseNotifier, http.Flusher, etc.
func (w *wrappedWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		ipAddress := utils.GetIpAddressFromRequest(r)

		end := time.Since(start).String()

		slog.Info("http request",
			"path", r.URL.Path,
			"method", r.Method,
			"status_code", wrapped.statusCode,
			"duration", end,
			"ip_address", ipAddress,
		)
	})
}
