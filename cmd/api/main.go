package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/observability"
	"github.com/ccrsxx/api/internal/server"
)

func main() {
	cfg := config.Load()

	// Configure logging (and error rendering) first so every subsequent line,
	// including telemetry setup, is properly structured.
	server.LoadLoaders(cfg)

	shutdownCtx, cancelShutdown := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	defer cancelShutdown()

	// Tracing must be initialized before the pgx pool and HTTP handlers are
	// built: otelpgx and otelhttp capture the global tracer provider at
	// construction time, so initializing later would silently lose all spans.
	// A telemetry failure must never prevent the API from serving, so we log
	// and continue with the returned no-op shutdown.
	shutdownTracing, err := observability.InitTracing(shutdownCtx, cfg)

	if err != nil {
		slog.Error("tracing init failed", "error", err)
	}

	pool, db := sqlc.NewQueries(shutdownCtx, cfg.DatabaseURL)

	defer pool.Close()

	server := server.New(shutdownCtx, cfg, pool, db)

	// Internal-only server for /metrics and /healthz. Its port is deliberately
	// not published; the collector scrapes it over the Docker network.
	metricsServer := observability.NewServer(cfg, pool)

	go func() {
		slog.Info("metrics server start listening", "port", metricsServer.Addr)

		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("metrics server stop listening", "error", err)
		}
	}()

	go func() {
		slog.Info("server start listening", "port", server.Addr, "env", cfg.AppEnv)

		if err := server.ListenAndServe(); err != nil {
			slog.Error("server stop listening", "error", err)
		}
	}()

	<-shutdownCtx.Done()

	slog.Info("server stopping gracefully")

	// Allow forced stop signal to exit immediately
	// Use case: if graceful shutdown is waiting too long, user can send
	// a second signal (CTRL+C) to force stop the application immediately
	cancelShutdown()

	// Give the server 60 seconds to shutdown gracefully
	// Basically a hard timeout to avoid hanging forever
	// Any open handler will not receive further requests
	// Ongoing handlers will have 60 seconds to finish before the application is forcefully terminated
	shutdownTimeoutCtx, cancelShutdownTimeout := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancelShutdownTimeout()

	if err := server.Shutdown(shutdownTimeoutCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}

	if err := metricsServer.Shutdown(shutdownTimeoutCtx); err != nil {
		slog.Error("metrics server shutdown failed", "error", err)
	}

	// Flush buffered spans last, so spans produced by in-flight requests during
	// the graceful shutdown window are still exported.
	if err := shutdownTracing(shutdownTimeoutCtx); err != nil {
		slog.Error("tracing shutdown failed", "error", err)
	}

	slog.Info("server stopped gracefully")
}
