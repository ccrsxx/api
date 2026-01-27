package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ccrsxx/api/src/internal/server"
)

func main() {
	server := server.NewServer()

	go func() {
		slog.Info("server start listening", "port", server.Addr, "env", config.Env().AppEnv)

		if err := server.ListenAndServe(); err != nil {
			slog.Error("server stop listening", "error", err)
		}
	}()

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	defer stop()

	<-shutdown.Done()

	slog.Info("server stopping gracefully")

	// Allow forced stop signal to exit immediately
	// Use case: if graceful shutdown is waiting too long, user can send
	// a second signal (CTRL+C) to force stop the application immediately
	stop()

	// Give the server 60 seconds to shutdown gracefully
	// Basically a hard timeout to avoid hanging forever
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}

	slog.Info("server stopped gracefully")
}
