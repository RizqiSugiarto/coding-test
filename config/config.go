package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App  App
		HTTP HTTP
		Log  Log
		PG   PG
	}

	// App -.
	App struct {
		Name    string `env-required:"true" env:"APP_NAME"`
		Version string `env-required:"true" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		Host     string `env-required:"true" env:"POSTGRES_HOST"`
		Port     int    `env-required:"true" env:"POSTGRES_PORT"`
		User     string `env-required:"true" env:"POSTGRES_USER"`
		Password string `env-required:"true" env:"POSTGRES_PASSWORD"`
		DBName   string `env-required:"true" env:"POSTGRES_DB"`
		PoolMax  int    `env-required:"true" env:"POSTGRES_POOL_MAX"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Load Config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("Read Config error: %w", err)
	}

	return cfg, nil
}
