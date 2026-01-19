package config

type appConfig struct {
	AppEnv        Environment
	IsProduction  bool
	IsDevelopment bool
}

var configInstance appConfig

func LoadConfig() {
	env := Env().AppEnv

	configInstance = appConfig{
		IsProduction:  env == EnvProduction,
		IsDevelopment: env == EnvDevelopment,
	}
}

func Config() *appConfig {
	return &configInstance
}
