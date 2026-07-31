package observability

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ccrsxx/api/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Pinger is satisfied by *pgxpool.Pool. Kept as an interface so this package
// does not depend on the database layer.
type Pinger interface {
	Ping(ctx context.Context) error
}

// NewServer builds the internal observability server exposing /metrics and
// /healthz.
//
// This is deliberately a SEPARATE server on its own port, not a route on the
// public router: the public router is exposed at BACKEND_PUBLIC_URL, and
// publishing /metrics there would leak the application's internal state and
// infrastructure inventory. The port is not published in compose; the
// collector scrapes it over the Docker network.
func NewServer(cfg config.AppConfig, pinger Pinger) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("GET /metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{
		// Surface scrape problems in our own logs rather than silently 500ing.
		ErrorLog:      slog.NewLogLogger(slog.Default().Handler(), slog.LevelWarn),
		ErrorHandling: promhttp.HTTPErrorOnError,
	}))

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		healthCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)

		defer cancel()

		if pinger != nil {
			if err := pinger.Ping(healthCtx); err != nil {
				slog.WarnContext(r.Context(), "health check failed", "error", err)

				http.Error(w, "database unavailable", http.StatusServiceUnavailable)

				return
			}
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("ok")); err != nil {
			slog.WarnContext(r.Context(), "health check write failed", "error", err)
		}
	})

	return &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.MetricsPort),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
