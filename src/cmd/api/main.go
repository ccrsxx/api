package main

import (
	"log/slog"
	"os"

	"github.com/ccrsxx/api-go/src/internal/config"
	"github.com/ccrsxx/api-go/src/internal/logger"
	"github.com/ccrsxx/api-go/src/internal/server"
)

func main() {
	config.LoadConfig()

	logger.Init()

	config.LoadEnv()

	server := server.NewServer()

	slog.Info("server started", "port", server.Addr, "env", config.Config().AppEnv)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
