package main

import (
	"context"
	"log/slog"
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

	server := server.New(shutdownCtx, cfg, db)

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

	slog.Info("server stopped gracefully")
}
