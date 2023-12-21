package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		GRPC `yaml:"grpc"`
		SQL  `yaml:"sql"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	GRPC struct {
		Port int `env-required:"true" yaml:"port" env:"GRPC_PORT"`
	}

	SQL struct {
		Timeout string `env-required:"true" yaml:"timeout" env:"SQL_TIMEOUT"`
		URL     string `env:"SQL_URL"`
	}
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func NewConfig(path string) (*Config, error) {
	config := &Config{}

	err := cleanenv.ReadConfig(path, config)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
