package config

import (
	"fmt"
	"os"
)

type config struct {
	Env   string
	Token string
}

func InitConfig() (*config, error) {
	var config config
	config.Env = envOrDefault("ENV", "prod")
	if config.Env != "prod" && config.Env != "dev" && config.Env != "local" {
		return nil, fmt.Errorf("Недопустимое значение ENV")
	}

	config.Token = os.Getenv("TOKEN")
	if config.Token == "" {
		return nil, fmt.Errorf("Токен не указан")
	}

	return &config, nil
}

func envOrDefault(varName string, defaultValue string) string {
	value := os.Getenv(varName)
	if value == "" {
		value = defaultValue
	}

	return value
}
