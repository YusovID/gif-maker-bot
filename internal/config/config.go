package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yml:"env" default:"local"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, errors.New("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		return nil, fmt.Errorf("config file does not exist: %w", err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	return &cfg, nil
}
