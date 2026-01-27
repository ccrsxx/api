package server

import (
	"github.com/ccrsxx/api/src/internal/config"
	"github.com/ccrsxx/api/src/internal/logger"
)

func RegisterLoaders() {
	config.LoadEnv()
	config.LoadConfig()

	logger.Init()
}
