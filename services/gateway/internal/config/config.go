package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Env string `env-required:"true" env:"ENV"`
		Http
		Services
		Logger
	}

	Http struct {
		Port string `env-required:"true" env:"HTTP_PORT"`
	}

	Services struct {
		AuthServiceUrl string `env-required:"true" env:"AUTH_SERVICE_URL"`
	}

	Logger struct {
		Level            int           `env:"LOGGER_LEVEL" env-default:"0"`
		LogFilePath      string        `env:"LOGGER_LOG_FILE_PATH"`
		BufferSize       int           `env:"LOGGER_BUFFER_SIZE"`
		FlushInterval    time.Duration `env:"LOGGER_FLUSH_INTERVAL"`
		FileMaxMegabytes int           `env:"LOGGER_FILE_MAX_MEGABYTES"`
		MaxBackups       int           `env:"LOGGER_MAX_BACKUPS"`
		MaxAgeDays       int           `env:"LOGGER_MAX_AGE_DAYS"`
	}
)

func New() (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	return &cfg, nil
}
