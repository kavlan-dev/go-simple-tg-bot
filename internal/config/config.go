package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Env   string `json:"env"`
	Token string `json:"token"`
}

func New() (*Config, error) {
	pathCmd := flag.String(
		"p",
		"config/config.json",
		"Введите относительный путь до файла конфигурации",
	)

	flag.Parse()

	return configFile(*pathCmd)
}

func configFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("Токен не указан")
	}

	if cfg.Env == "" {
		return nil, fmt.Errorf("Env не указан")
	}

	return &cfg, nil
}
