package config

import (
	"fmt"
	"log"
	"time"

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
		JWT
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

	// JWT -.
	JWT struct {
		AccessTokenSecretKey  string        `env-required:"true" env:"ACCESS_TOKEN_SECRET_KEY"`
		RefreshTokenSecretKey string        `env-required:"true" env:"REFRESH_TOKEN_SECRET_KEY"`
		AccessTokenTTL        time.Duration `env-required:"true" env:"ACCESS_TOKEN_TTL"`
		RefreshTokenTTL       time.Duration `env-required:"true" env:"REFRESH_TOKEN_TTL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := godotenv.Load()
	if err != nil {
		log.Printf("error load env file: %s", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	return cfg, nil
}
