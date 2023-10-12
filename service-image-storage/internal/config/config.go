package config

import (
	"fmt"
	"path/filepath"

	"github.com/caarlos0/env/v9"
)

type Env string

var (
	EnvLocal = Env("local")
	EnvProd  = Env("prod")
)

type Config struct {
	Env Env `env:"APP_ENV,required"`

	ServerConfig  ServerConfig
	StorageConfig StorageConfig
}

type ServerConfig struct {
	Port int `env:"SERVER_PORT,required"`
}

type StorageConfig struct {
	FolderPath string `env:"IMAGES_FOLDER_PATH,required"`
}

func Load() (Config, error) {
	c := Config{}
	err := env.Parse(&c)
	if err != nil {
		return Config{}, fmt.Errorf("parse env variables to config: %w", err)
	}

	c.StorageConfig.FolderPath, err = filepath.Abs(c.StorageConfig.FolderPath)
	if err != nil {
		return Config{}, fmt.Errorf("convert folder path to absoulte: %w", err)
	}

	return c, nil
}
