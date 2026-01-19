package server

import (
	"github.com/ccrsxx/api-go/src/internal/config"
	"github.com/ccrsxx/api-go/src/internal/logger"
)

func RegisterLoaders() {
	config.LoadEnv()
	config.LoadConfig()

	logger.Init()
}
