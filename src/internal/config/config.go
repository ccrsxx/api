package config

import "os"

type appConfig struct {
	IsProduction  bool
	IsDevelopment bool
}

var configInstance = appConfig{
	IsProduction:  os.Getenv("APP_ENV") == "production",
	IsDevelopment: os.Getenv("APP_ENV") == "development",
}

func Config() *appConfig {
	return &configInstance
}
