package server

import (
	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/logger"
)

func InitLoaders(cfg config.AppConfig) {
	api.Init(cfg.IsDevelopment)
	logger.Init(cfg.IsDevelopment)
}
