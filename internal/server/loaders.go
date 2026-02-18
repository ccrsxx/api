package server

import (
	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/logger"
)

func RegisterLoaders() {
	config.LoadEnv()
	config.LoadConfig()

	logger.Init()
}
