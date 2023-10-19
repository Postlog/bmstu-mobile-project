package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Env string

var (
	EnvLocal = Env("local")
	EnvProd  = Env("prod")
)

type Config struct {
	Env Env `env:"APP_ENV,required"`

	DependenciesConfig DependenciesConfig
	RabbitConfig       RabbitConfig
	PostgresConfig     PostgresConfig
}

type ServerConfig struct {
	Port int `env:"SERVER_PORT,required"`
}

type RabbitConfig struct {
	DSN string `env:"RABBITMQ_DSN,required"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     string `env:"POSTGRES_PORT,required"`
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	DBName   string `env:"POSTGRES_DB_NAME,required"`
}

type DependenciesConfig struct {
	ServiceImageScalerURL    string        `env:"SERVICE_IMAGE_SCALER_URL,required"`
	ServiceImageScalerTimout time.Duration `env:"SERVICE_IMAGE_SCALER_TIMEOUT,required"`
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("error loading .env file")
	}

	c := Config{}
	err := env.Parse(&c)
	if err != nil {
		return Config{}, fmt.Errorf("parse env variables to config: %w", err)
	}

	return c, nil
}
