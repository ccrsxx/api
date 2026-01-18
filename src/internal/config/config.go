package config

import (
	"os"
)

type Environment string

const (
	EnvProduction  Environment = "production"
	EnvDevelopment Environment = "development"
)

type appConfig struct {
	AppEnv        Environment
	IsProduction  bool
	IsDevelopment bool
}

var configInstance appConfig

func LoadConfig() {
	val := os.Getenv("APP_ENV")

	env := EnvDevelopment

	if val == string(EnvProduction) {
		env = EnvProduction
	}

	configInstance = appConfig{
		AppEnv:        env,
		IsProduction:  env == EnvProduction,
		IsDevelopment: env == EnvDevelopment,
	}
}

func Config() *appConfig {
	return &configInstance
}
