package server

import (
	"github.com/ccrsxx/api/internal/api"
	"github.com/ccrsxx/api/internal/config"
	"github.com/ccrsxx/api/internal/logger"
)

func LoadLoaders(cfg config.AppConfig) {
	api.Load(cfg.IsDevelopment)
	logger.Load(cfg.IsDevelopment)
}
