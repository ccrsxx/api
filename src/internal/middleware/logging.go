package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ccrsxx/api-go/src/internal/utils"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
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

		end := time.Since(start)

		parsedEnd := fmt.Sprintf("%v", end)

		slog.Info("http request",
			"path", r.URL.Path,
			"method", r.Method,
			"status_code", wrapped.statusCode,
			"duration", parsedEnd,
			"ip_address", ipAddress,
		)
	})
}
