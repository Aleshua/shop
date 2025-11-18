package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Env string `env-required:"true" env:"ENV"`
		Logger
		HTTP
		Database
		EmailVerify
		Email
		JWT
	}

	Logger struct {
		Level            int8          `env:"LOGGER_LEVEL" env-default:"0"`
		LogFilePath      string        `env:"LOGGER_LOG_FILE_PATH"`
		BufferSize       int           `env:"LOGGER_BUFFER_SIZE"`
		FlushInterval    time.Duration `env:"LOGGER_FLUSH_INTERVAL"`
		FileMaxMegabytes int           `env:"LOGGER_FILE_MAX_MEGABYTES"`
		MaxBackups       int           `env:"LOGGER_MAX_BACKUPS"`
		MaxAgeDays       int           `env:"LOGGER_MAX_AGE_DAYS"`
	}

	HTTP struct {
		Host              string        `env-required:"true" env:"HTTP_HOST"`
		Port              string        `env-required:"true" env:"HTTP_PORT"`
		ReadTimeout       time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"10s"`
		ReadHeaderTimeout time.Duration `env:"HTTP_READ_HEADER_TIMEOUT" env-default:"5s"`
		WriteTimeout      time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
		IdleTimeout       time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"30s"`
		MaxHeaderBytes    int           `env:"HTTP_MAX_HEADER_BYTES" env-default:"1048576"`
	}

	Database struct {
		Path              string        `env-required:"true" env:"DATABASE_PATH"`
		MaxConns          int32         `env:"DATABASE_MAX_CONNS" env-default:"10"`
		MinConns          int32         `env:"DATABASE_MIN_CONNS" env-default:"2"`
		MaxConnLifetime   time.Duration `env:"DATABASE_MAX_CONN_LIFETIME" env-default:"180s"`
		MaxConnIdleTime   time.Duration `env:"DATABASE_MAX_CONN_IDLETIME" env-default:"30s"`
		HealthCheckPeriod time.Duration `env:"DATABASE_HEALTH_CHECK_PERIOD" env-default:"60s"`
		ConnectTimeout    time.Duration `env:"DATABASE_CONNECT_TIMEOUT" env-default:"5s"`
	}

	EmailVerify struct {
		CodeLength  int `env-required:"true" env:"EMAIL_VERIFY_CODE_LENGTH"`
		MaxAttempts int `env-required:"true" env:"EMAIL_VERIFY_MAX_ATTEMPTS"`
		TTLSeconds  int `env-required:"true" env:"EMAIL_VERIFY_TTL_SECONDS"`
	}

	Email struct {
		Host     string `env:"EMAIL_HOST"`
		Port     int    `env:"EMAIL_PORT"`
		Username string `env:"EMAIL_USERNAME"`
		Password string `env:"EMAIL_PASSWORD"`
		From     string `env:"EMAIL_FROM"`
	}

	JWT struct {
		RefreshTokenLength     int   `env-required:"true" env:"JWT_REFRESH_TOKEN_LENGTH"`
		AccessTokenTTLSeconds  int64 `env-required:"true" env:"JWT_ACCESS_TOKEN_TTL_SECONDS"`
		RefreshTokenTTLSeconds int64 `env-required:"true" env:"JWT_REFRESH_TOKEN_TTL_SECONDS"`
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
