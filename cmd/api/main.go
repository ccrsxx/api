package main

import (
	_ "time/tzdata"

	_ "golang.org/x/crypto/x509roots/fallback"

	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/db/sqlc"
	"github.com/ccrsxx/api/internal/server"
)

func main() {
	cfg := config.Load()

	shutdownCtx, cancelShutdown := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	defer cancelShutdown()

	pool, db := sqlc.NewQueries(shutdownCtx, cfg.DatabaseURL)

	defer pool.Close()

	mainServer := server.New(shutdownCtx, cfg, pool, db)
	monitoringServer := server.NewMonitoringServer(cfg)

	go func() {
		slog.Info("main server start listening", "port", mainServer.Addr, "env", cfg.AppEnv)

		if err := mainServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("main server stop listening", "error", err)
		}
	}()

	go func() {
		slog.Info("monitoring server start listening", "port", monitoringServer.Addr, "env", cfg.AppEnv)

		if err := monitoringServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("monitoring server stop listening", "error", err)
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

	if err := mainServer.Shutdown(shutdownTimeoutCtx); err != nil {
		slog.Error("main server shutdown failed", "error", err)
	}

	if err := monitoringServer.Shutdown(shutdownTimeoutCtx); err != nil {
		slog.Error("monitoring server shutdown failed", "error", err)
	}

	slog.Info("server stopped gracefully")
}
