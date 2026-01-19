package main

import (
	"log/slog"
	"os"

	"github.com/ccrsxx/api-go/src/internal/config"
	"github.com/ccrsxx/api-go/src/internal/server"
)

func main() {
	server := server.NewServer()

	slog.Info("server started", "port", server.Addr, "env", config.Env().AppEnv, "best_girl", config.Env().SecretKey)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
