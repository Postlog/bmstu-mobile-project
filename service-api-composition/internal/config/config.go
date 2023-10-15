package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v9"
)

type Env string

var (
	EnvLocal = Env("local")
	EnvProd  = Env("prod")
)

type Config struct {
	Env Env `env:"APP_ENV,required"`

	ServerConfig       ServerConfig
	DependenciesConfig DependenciesConfig
}

type ServerConfig struct {
	Port int `env:"SERVER_PORT,required"`
}

type DependenciesConfig struct {
	ServiceImageStorageURL     string        `env:"SERVICE_IMAGE_STORAGE_URL,required"`
	ServiceImageStorageTimeout time.Duration `env:"SERVICE_IMAGE_STORAGE_TIMEOUT,required"`

	RabbitMQDSN string `env:"RABBITMQ_DSN,required"`
}

func Load() (Config, error) {
	c := Config{}
	err := env.Parse(&c)
	if err != nil {
		return Config{}, fmt.Errorf("parse env variables to config: %w", err)
	}

	return c, nil
}
