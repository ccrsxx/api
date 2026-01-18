package main

import (
	"log/slog"
	"os"

	"github.com/ccrsxx/api-go/src/internal/config"
	"github.com/ccrsxx/api-go/src/internal/logger"
	"github.com/ccrsxx/api-go/src/internal/server"
)

func main() {
	logger.Init()
	config.LoadEnv()

	server := server.NewServer()

	slog.Info("Server starting", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("Cannot start server", "error", err)
		os.Exit(1)
	}
}
