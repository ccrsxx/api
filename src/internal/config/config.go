package config

type appConfig struct {
	AppEnv        EnvironmentApp
	IsProduction  bool
	IsDevelopment bool
}

var configInstance appConfig

func LoadConfig() {
	env := Env().AppEnv

	configInstance = appConfig{
		IsProduction:  env == EnvironmentProduction,
		IsDevelopment: env == EnvironmentDevelopment,
	}
}

func Config() *appConfig {
	return &configInstance
}
